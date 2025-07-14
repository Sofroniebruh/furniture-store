package services

import (
	"context"
	"encoding/json"
	"fmt"
	"furniture-store-backend/db"
	"furniture-store-backend/models"
	"furniture-store-backend/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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
	event := r.FormValue("event")
	model := r.FormValue("model")
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

	if model == "" || event == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Model and Event cannot be empty",
		})
		return
	}

	err = db.DB.QueryRow(`
				INSERT INTO products (name, amount, price, picture_urls, description, event, model)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				RETURNING id`,
		name,
		amount,
		price,
		pq.StringArray(pictureURLs),
		description,
		event,
		model).Scan(&product.ID)

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
	product.Event = event
	product.Model = model

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]models.Product{
		"created": product,
	})
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	priceFromStr := r.URL.Query().Get("price_from")
	priceToStr := r.URL.Query().Get("price_to")
	priceFromInt, err := strconv.Atoi(priceFromStr)

	if err != nil && priceFromStr != "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to parse price_from",
		})
		return
	}

	priceToInt, err := strconv.Atoi(priceToStr)

	if err != nil && priceToStr != "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to parse price_to",
		})
		return
	}

	var filters []string
	params := map[string]interface{}{}
	event := r.URL.Query().Get("event")
	model := r.URL.Query().Get("model")

	if event != "" {
		filters = append(filters, "event = :event")
		params["event"] = event
	}
	if model != "" {
		filters = append(filters, "model = :model")
		params["model"] = model
	}
	if priceFromStr != "" && priceFromInt > 0 {
		filters = append(filters, "price >= :price_from")
		params["price_from"] = priceFromInt
	}
	if priceToStr != "" && priceToInt > 0 {
		filters = append(filters, "price <= :price_to")
		params["price_to"] = priceToInt
	}

	page, err := strconv.Atoi(pageStr)

	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)

	if err != nil || limit < 12 {
		limit = 12
	}

	offset := (page - 1) * limit
	params["limit"] = limit
	params["offset"] = offset
	var products []models.Product

	query := "SELECT * FROM products"
	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}
	query += " LIMIT :limit OFFSET :offset"
	rows, err := db.DB.NamedQuery(query, params)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to get products",
		})
		return
	}

	defer func(rows *sqlx.Rows) {
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
		err = rows.Scan(&product.ID, &product.Name, &product.Amount, &product.Price, &product.Description, &pictureUrls, &product.Event, &product.Model)
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
	if event := r.FormValue("event"); event != "" {
		updates = append(updates, "event = :event")
		params["event"] = event
		product.Event = event
	}
	if model := r.FormValue("model"); model != "" {
		updates = append(updates, "model = :model")
		params["model"] = model
		product.Model = model
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

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productIdStr := r.URL.Query().Get("id")
	productId, err := uuid.Parse(productIdStr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Product ID should be a valid UUID",
		})
		return
	}

	result, err := db.DB.Exec("DELETE FROM products WHERE id = $1", productId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to delete product",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Product not found",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Product deleted successfully",
	})
}
