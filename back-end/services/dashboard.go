package services

import (
	"database/sql"
	"encoding/json"
	"furniture-store-backend/db"
	"furniture-store-backend/models"
	"log"
	"net/http"
	"strconv"
)

var productRequest struct {
	Name        string  `json:"name"`
	Amount      int     `json:"amount"`
	Price       float64 `json:"price"`
	PictureUrl  string  `json:"pictureUrl"`
	Description string  `json:"description"`
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&productRequest)
	var product models.Product

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to parse request body",
		})
		return
	}

	if productRequest.Price <= 0 || productRequest.Amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Price AND/OR amount must be greater than zero",
		})
		return
	}

	if productRequest.PictureUrl == "" || productRequest.Name == "" || productRequest.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Image URL, Name and Description cannot be empty",
		})
		return
	}

	err = db.DB.QueryRow(`
				INSERT INTO products (name, amount, price, picture_url, description)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id`,
		productRequest.Name,
		productRequest.Amount,
		productRequest.Price,
		productRequest.PictureUrl,
		productRequest.Description).Scan(&product.ID)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to create a product",
		})
		return
	}

	product.Name = productRequest.Name
	product.Amount = productRequest.Amount
	product.Price = productRequest.Price
	product.PictureUrl = productRequest.PictureUrl
	product.Description = productRequest.Description

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
		err = rows.Scan(&product.ID, &product.Name, &product.Amount, &product.Price, &product.PictureUrl, &product.Description)

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
