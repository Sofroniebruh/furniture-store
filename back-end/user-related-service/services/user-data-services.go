package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"log"
	"net/http"
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
	Amount      int            `db:"amount"`
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
			p.id AS product_id, p.name, p.description, p.amount, p.price, p.picture_urls, p.event, p.model,
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
				INSERT INTO histories (user_id, product_id) 
				VALUES ($1, $2)
				`
	default:
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid type in request body",
		})
		return
	}

	_, err = db.DB.Exec(query, userId, product.ID)

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

	switch urlType {
	case "/user/wishlist":
		query = `
            SELECT p.* 
            FROM products p
            JOIN wishlists w ON p.id = w.product_id
            WHERE w.user_id = $1
            `
	case "/user/history":
		query = `
            SELECT p.* 
            FROM products p
            JOIN histories w ON p.id = w.product_id
            WHERE w.user_id = $1
            `
	default:
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid url type",
		})
		return
	}

	rows, err := db.DB.Query(query, userId)

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

		err = rows.Scan(&product.ID, &product.Name, &product.Amount, &product.Price, &product.Description, &pictureUrls, &product.Event, &product.Model)

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

		products = append(products, product)
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"products": products,
	})
}
