package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"products-service/db"
	"products-service/models"
)

func AddStockHistoryEntry(productID uuid.UUID, historyType models.StockHistoryType, quantity, previousStock, newStock int, reason string) error {
	_, err := db.DB.Exec(`
		INSERT INTO stock_history (product_id, type, quantity, previous_stock, new_stock, reason)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, productID, historyType, quantity, previousStock, newStock, reason)

	return err
}

func UpdateProductStock(w http.ResponseWriter, r *http.Request) {
	var request models.StockHistoryRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	if request.Quantity == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Quantity cannot be zero",
		})
		return
	}

	var product models.Product
	err = db.DB.Get(&product, "SELECT id, name, stock, price, picture_urls, description, event, model FROM products WHERE id = $1", request.ProductID)

	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Product not found",
		})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to get product",
		})
		return
	}

	previousStock := product.Stock
	newStock := previousStock + request.Quantity

	if newStock < 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Insufficient stock",
		})
		return
	}

	var historyType models.StockHistoryType
	if request.Quantity > 0 {
		historyType = models.StockHistoryTypeIn
	} else {
		historyType = models.StockHistoryTypeOut
	}

	tx, err := db.DB.Beginx()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to start transaction",
		})
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE products SET stock = $1 WHERE id = $2", newStock, request.ProductID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to update product stock",
		})
		return
	}

	_, err = tx.Exec(`
		INSERT INTO stock_history (product_id, type, quantity, previous_stock, new_stock, reason)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, request.ProductID, historyType, abs(request.Quantity), previousStock, newStock, request.Reason)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to create stock history entry",
		})
		return
	}

	if err = tx.Commit(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to commit transaction",
		})
		return
	}

	product.Stock = newStock

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Stock updated successfully",
		"product": product,
	})
}

func GetStockHistory(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	productIDStr := r.URL.Query().Get("product_id")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	var baseQuery string
	var countQuery string
	var args []interface{}

	if productIDStr != "" {
		productID, err := uuid.Parse(productIDStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid product ID",
			})
			return
		}

		baseQuery = `
			SELECT sh.id, sh.product_id, sh.type, sh.quantity, sh.previous_stock, 
				   sh.new_stock, sh.reason, sh.created_at,
				   p.name, p.stock, p.price, p.picture_urls, p.description, p.event, p.model
			FROM stock_history sh
			JOIN products p ON sh.product_id = p.id
			WHERE sh.product_id = $1
			ORDER BY sh.created_at DESC
			LIMIT $2 OFFSET $3
		`
		countQuery = "SELECT COUNT(*) FROM stock_history WHERE product_id = $1"
		args = []interface{}{productID, limit, offset}
	} else {
		baseQuery = `
			SELECT sh.id, sh.product_id, sh.type, sh.quantity, sh.previous_stock, 
				   sh.new_stock, sh.reason, sh.created_at,
				   p.name, p.stock, p.price, p.picture_urls, p.description, p.event, p.model
			FROM stock_history sh
			JOIN products p ON sh.product_id = p.id
			ORDER BY sh.created_at DESC
			LIMIT $1 OFFSET $2
		`
		countQuery = "SELECT COUNT(*) FROM stock_history"
		args = []interface{}{limit, offset}
	}

	var totalCount int
	if productIDStr != "" {
		productID, _ := uuid.Parse(productIDStr)
		err = db.DB.Get(&totalCount, countQuery, productID)
	} else {
		err = db.DB.Get(&totalCount, countQuery)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to get total count",
		})
		return
	}

	rows, err := db.DB.Query(baseQuery, args...)
	if err != nil {
		log.Printf("Query error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to get stock history",
		})
		return
	}
	defer rows.Close()

	var stockHistory []models.StockHistory

	for rows.Next() {
		var history models.StockHistory
		var product models.Product

		err = rows.Scan(
			&history.ID, &history.ProductID, &history.Type, &history.Quantity,
			&history.PreviousStock, &history.NewStock, &history.Reason, &history.CreatedAt,
			&product.Name, &product.Stock, &product.Price, &product.PictureUrls,
			&product.Description, &product.Event, &product.Model,
		)

		if err != nil {
			log.Printf("Scan error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Failed to scan stock history data",
			})
			return
		}

		product.ID = history.ProductID
		history.Product = &product
		stockHistory = append(stockHistory, history)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"stock_history": stockHistory,
		"total_pages":   totalPages,
		"current_page":  page,
		"total_count":   totalCount,
	})
}

func GetProductStockHistory(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "id")
	productID, err := uuid.Parse(productIDStr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid product ID",
		})
		return
	}

	r = r.WithContext(r.Context())
	query := r.URL.Query()
	query.Set("product_id", productID.String())
	r.URL.RawQuery = query.Encode()

	GetStockHistory(w, r)
}
