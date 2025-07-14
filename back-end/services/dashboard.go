package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"furniture-store-backend/db"
	"furniture-store-backend/models"
	"furniture-store-backend/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func uploadImageToBucket(header *multipart.FileHeader, ch chan<- string, file multipart.File) {
	s3Client, err := utils.LoadS3Client()

	if err != nil {
		log.Printf("failed to load S3 client: %v", err)
		ch <- ""
		return
	}

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("BUCKET_NAME")),
		Key:         aws.String("public-images/" + header.Filename),
		Body:        file,
		ContentType: aws.String(header.Header.Get("Content-Type")),
	})

	if err != nil {
		log.Printf("failed to load S3 client: %v", err)
		ch <- ""
		return
	}

	bucket := os.Getenv("BUCKET_NAME")
	region := os.Getenv("AWS_REGION")
	key := "public-images/" + header.Filename
	ch <- fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key)
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	s3ch := make(chan string)

	err := r.ParseMultipartForm(10 << 20)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Unable to parse form",
		})
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	price, err := strconv.ParseFloat(r.FormValue("price"), 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to parse price",
		})
		return
	}

	amount, err := strconv.Atoi(r.FormValue("amount"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to parse amount",
		})
		return
	}

	files := r.MultipartForm.File["pictures"]

	if files == nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "No files provided",
		})
		return
	}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		header := fileHeader

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Failed to open file",
			})
			return
		}
		file.Close()

		go uploadImageToBucket(header, s3ch, file)
	}

	var pictureURLs []string

	for i := 0; i < len(files); i++ {
		pictureURL := <-s3ch

		if pictureURL == "" {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Failed to upload image to  the bucket",
			})
			return
		}

		pictureURLs = append(pictureURLs, pictureURL)
	}

	if price <= 0 || amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Price AND/OR amount must be greater than zero",
		})
		return
	}

	if name == "" || description == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Name and Description cannot be empty",
		})
		return
	}

	err = db.DB.QueryRow(`
				INSERT INTO products (name, amount, price, picture_urls, description)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id`,
		name,
		amount,
		price,
		pq.StringArray(pictureURLs),
		description).Scan(&product.ID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to create a product",
		})
		return
	}

	product.Name = name
	product.Amount = amount
	product.Price = float64(price)
	product.PictureUrls = pictureURLs
	product.Description = description

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]models.Product{
		"created": product,
	})
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)

	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)

	if err != nil || limit < 12 {
		limit = 12
	}

	offset := (page - 1) * limit

	var products []models.Product

	rows, err := db.DB.Query(`
				SELECT * FROM products 
				LIMIT $1 OFFSET $2`, limit, offset)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to get products",
		})
		return
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Internal Server Error",
			})
			return
		}
	}(rows)

	for rows.Next() {
		var product models.Product
		var pictureUrls pq.StringArray
		err = rows.Scan(&product.ID, &product.Name, &product.Amount, &product.Price, &product.Description, &pictureUrls)
		product.PictureUrls = pictureUrls

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Internal Server Error",
			})
			return
		}
		products = append(products, product)
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string][]models.Product{
		"products": products,
	})
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productIdStr := r.URL.Query().Get("id")
	productId, err := uuid.Parse(productIdStr)
	var product models.Product
	s3ch := make(chan string)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Product ID should be a valid UUID",
		})
		return
	}

	err = r.ParseMultipartForm(10 << 20)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Unable to parse form",
		})
		return
	}

	err = db.DB.Get(&product, "SELECT * FROM products WHERE id = $1", productId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Product not found",
		})
		return
	}

	var updates []string
	params := map[string]interface{}{
		"id": productId,
	}

	if name := r.FormValue("name"); name != "" {
		updates = append(updates, "name = :name")
		params["name"] = name
		product.Name = name
	}
	if description := r.FormValue("description"); description != "" {
		updates = append(updates, "description = :description")
		params["description"] = description
		product.Description = description
	}
	if priceStr := r.FormValue("price"); priceStr != "" {
		price, err := strconv.ParseFloat(priceStr, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Failed to parse price",
			})
			return
		}

		updates = append(updates, "price = :price")
		params["price"] = price
		product.Price = price
	}
	if amountStr := r.FormValue("amount"); amountStr != "" {
		amount, err := strconv.Atoi(amountStr)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Failed to parse amount",
			})
			return
		}

		updates = append(updates, "amount = :amount")
		params["amount"] = amount
		product.Amount = amount
	}
	if files := r.MultipartForm.File["pictures"]; files != nil {
		for _, fileHeader := range files {
			file, err := fileHeader.Open()

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"error": "Failed to open file",
				})
				return
			}

			go uploadImageToBucket(fileHeader, s3ch, file)
		}

		var fileURLs []string
		for i := 0; i < len(files); i++ {
			fileURL := <-s3ch

			if fileURL == "" {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"error": "Failed to upload image to the bucket",
				})
				return
			}
			fileURLs = append(fileURLs, fileURL)
		}
		updates = append(updates, "picture_urls = :fileURLs")
		params["fileURLs"] = pq.StringArray(fileURLs)
		product.PictureUrls = fileURLs
	}

	if len(updates) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "No fields to update",
		})
		return
	}

	query := "UPDATE products SET " + strings.Join(updates, ", ") + " WHERE id = :id"
	_, err = db.DB.NamedExec(query, params)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to update products",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]models.Product{
		"updated": product,
	})
}
