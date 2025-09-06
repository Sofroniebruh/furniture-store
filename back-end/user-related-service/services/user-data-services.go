package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
	"user-related-service/config"
	"user-related-service/db"
	"user-related-service/models"
)

type UserDataWithCodeResponse struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
type UserDataWithEmail struct {
	Email string `json:"email"`
}

type ProductWishlistOrPurchaseHistory struct {
	ProductId string `json:"product_id"`
}

type ProductWithColorRow struct {
	ProductID   uuid.UUID      `db:"product_id"`
	Name        string         `db:"name"`
	Description string         `db:"description"`
	Stock       int            `db:"stock"`
	Price       float64        `db:"price"`
	PictureUrls pq.StringArray `db:"picture_urls"`
	Event       string         `db:"event"`
	Model       string         `db:"model"`
	ColorID     uuid.UUID      `db:"color_id"`
	ColorName   string         `db:"color_name"`
}

func GetProductsWithColors(product models.Product) (models.Product, error) {
	var rows []ProductWithColorRow

	err := db.DB.Select(&rows, `
		SELECT 
			p.id AS product_id, p.name, p.description, p.stock, p.price, p.picture_urls, p.event, p.model,
			c.id AS color_id, c.name AS color_name
		FROM products p
		JOIN product_colors pc ON p.id = pc.product_id
		JOIN colors c ON c.id = pc.color_id
		WHERE p.id = $1
	`, product.ID)

	if err != nil {
		return models.Product{}, err
	}

	if len(rows) == 0 {
		return models.Product{}, sql.ErrNoRows
	}

	for _, row := range rows {
		product.Colors = append(product.Colors, models.Color{
			ID:   row.ColorID,
			Name: row.ColorName,
		})
	}

	return product, nil
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(config.UserIdKey).(uuid.UUID)
	var user models.User

	err := db.DB.Get(&user, "SELECT * FROM users WHERE id = $1", userId)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "User not found",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]models.User{
		"user": user,
	})
}

func AddToWishListOrHistory(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(config.UserIdKey).(uuid.UUID)
	var productData ProductWishlistOrPurchaseHistory
	var product models.Product
	var query string
	var urlPath = r.URL.Path

	err := json.NewDecoder(r.Body).Decode(&productData)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	err = db.DB.Get(&product, "SELECT * FROM products WHERE id = $1", productData.ProductId)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Product not found",
		})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to get product data",
		})
		return
	}

	switch urlPath {
	case "/user/wishlist":
		query = `
				INSERT INTO wishlists (user_id, product_id) 
				VALUES ($1, $2)
				`
	case "/user/history":
		query = `
				INSERT INTO histories (user_id, product_id, created_at) 
				VALUES ($1, $2, $3)
				`
	default:
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid type in request body",
		})
		return
	}

	_, err = db.DB.Exec(query, userId, product.ID, time.Now())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to insert product data",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]models.Product{
		"product": product,
	})
}

func RemoveFromWishList(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(config.UserIdKey).(uuid.UUID)
	var productData ProductWishlistOrPurchaseHistory
	var product models.Product

	err := json.NewDecoder(r.Body).Decode(&productData)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	err = db.DB.Get(&product, "SELECT * FROM products WHERE id = $1", productData.ProductId)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Product not found",
		})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to get product data",
		})
		return
	}

	_, err = db.DB.Exec("DELETE FROM wishlists WHERE product_id = $1 AND user_id = $2", product.ID, userId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to remove from wishlist",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]models.Product{
		"product": product,
	})
}

func GetWishlistOrHistoryPerUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(config.UserIdKey).(uuid.UUID)
	var products []models.Product
	urlType := r.URL.Path
	var query string
	var totalCount int
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	page, err := strconv.Atoi(pageStr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	offset := (page - 1) * limit
	switch urlType {
	case "/user/wishlist":
		query = `
            SELECT p.* 
            FROM products p
            JOIN wishlists w ON p.id = w.product_id
            WHERE w.user_id = $1
            LIMIT $2
            OFFSET $3
            
            `
		err = db.DB.Get(&totalCount, "SELECT COUNT(*) FROM wishlists WHERE user_id = $1", userId)
	case "/user/history":
		query = `
            SELECT p.* 
            FROM products p
            JOIN histories w ON p.id = w.product_id
            WHERE w.user_id = $1
        	LIMIT $2
            OFFSET $3
            `
		err = db.DB.Get(&totalCount, "SELECT COUNT(*) FROM histories WHERE user_id = $1", userId)
	default:
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid url type",
		})
		return
	}

	rows, err := db.DB.Query(query, userId, limit, offset)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to get product data",
		})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var product models.Product
		var pictureUrls pq.StringArray

		err = rows.Scan(&product.ID, &product.Name, &product.Stock, &product.Price, &product.Description, &pictureUrls, &product.Event, &product.Model)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Failed to get user's product id",
			})
			return
		}

		product, err = GetProductsWithColors(product)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Failed to get product data with colors",
			})
			return
		}

		product.PictureUrls = pictureUrls

		products = append(products, product)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"products": products,
		"total":    totalPages,
	})
}

func AddToHistoryInternal(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		UserID    string `json:"user_id"`
		ProductID string `json:"product_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	if requestBody.UserID == "" || requestBody.ProductID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "user_id and product_id are required",
		})
		return
	}

	userID, err := uuid.Parse(requestBody.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid user_id format",
		})
		return
	}

	productUUID, err := uuid.Parse(requestBody.ProductID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid product_id format",
		})
		return
	}

	var productExists bool
	err = db.DB.Get(&productExists, "SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)", productUUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to verify product",
		})
		return
	}

	if !productExists {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Product not found",
		})
		return
	}

	query := `INSERT INTO histories (user_id, product_id, created_at) VALUES ($1, $2, NOW()) 
			   ON CONFLICT (user_id, product_id) DO UPDATE SET created_at = NOW()`
	_, err = db.DB.Exec(query, userID, productUUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to add to history",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Successfully added to history",
	})
}
