package services

import (
	"database/sql"
	"errors"

	"cart-service/db"
	"cart-service/models"

	"github.com/google/uuid"
)

func GetCart(userID uuid.UUID) (*models.Cart, error) {
	cart := &models.Cart{
		UserID: userID,
		Items:  []models.CartItem{},
	}

	query := `
		SELECT 
			ci.id, ci.user_id, ci.product_id, ci.quantity, ci.created_at, ci.updated_at,
			p.name, p.description, p.stock, p.price, p.picture_urls, p.event, p.model
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.user_id = $1
		ORDER BY ci.created_at DESC
	`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalPrice float64
	var totalItems int

	for rows.Next() {
		var item models.CartItem
		var product models.Product

		err := rows.Scan(
			&item.ID, &item.UserID, &item.ProductID, &item.Quantity, &item.CreatedAt, &item.UpdatedAt,
			&product.Name, &product.Description, &product.Stock, &product.Price, &product.PictureUrls, &product.Event, &product.Model,
		)
		if err != nil {
			return nil, err
		}

		product.ID = item.ProductID
		item.Product = &product
		cart.Items = append(cart.Items, item)

		totalPrice += product.Price * float64(item.Quantity)
		totalItems += item.Quantity
	}

	cart.TotalPrice = totalPrice
	cart.TotalItems = totalItems

	return cart, nil
}

func AddToCart(userID, productID uuid.UUID, quantity int) error {
	var product models.Product
	err := db.DB.Get(&product, "SELECT * FROM products WHERE id = $1", productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("product not found")
		}
		return err
	}

	if product.Stock < quantity {
		return errors.New("insufficient stock")
	}

	var existingItem models.CartItem
	err = db.DB.Get(&existingItem, "SELECT * FROM cart_items WHERE user_id = $1 AND product_id = $2", userID, productID)
	
	if err == nil {
		newQuantity := existingItem.Quantity + quantity
		if product.Stock < newQuantity {
			return errors.New("insufficient stock")
		}

		_, err = db.DB.Exec("UPDATE cart_items SET quantity = $1 WHERE id = $2", newQuantity, existingItem.ID)
		return err
	} else if errors.Is(err, sql.ErrNoRows) {
		_, err = db.DB.Exec(
			"INSERT INTO cart_items (user_id, product_id, quantity) VALUES ($1, $2, $3)",
			userID, productID, quantity,
		)
		return err
	}

	return err
}

func UpdateCartItem(userID, itemID uuid.UUID, quantity int) error {
	var item models.CartItem
	err := db.DB.Get(&item, "SELECT * FROM cart_items WHERE id = $1 AND user_id = $2", itemID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("cart item not found")
		}
		return err
	}

	var product models.Product
	err = db.DB.Get(&product, "SELECT * FROM products WHERE id = $1", item.ProductID)
	if err != nil {
		return err
	}

	if product.Stock < quantity {
		return errors.New("insufficient stock")
	}

	_, err = db.DB.Exec("UPDATE cart_items SET quantity = $1 WHERE id = $2", quantity, itemID)
	return err
}

func RemoveFromCart(userID, itemID uuid.UUID) error {
	result, err := db.DB.Exec("DELETE FROM cart_items WHERE id = $1 AND user_id = $2", itemID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("cart item not found")
	}

	return nil
}

func ClearCart(userID uuid.UUID) error {
	_, err := db.DB.Exec("DELETE FROM cart_items WHERE user_id = $1", userID)
	return err
}

func ValidateCartForCheckout(userID uuid.UUID) error {
	query := `
		SELECT ci.quantity, p.stock, p.name
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.user_id = $1
	`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var cartQuantity, productStock int
		var productName string
		
		err := rows.Scan(&cartQuantity, &productStock, &productName)
		if err != nil {
			return err
		}

		if productStock < cartQuantity {
			return errors.New("insufficient stock for product: " + productName)
		}
	}

	return nil
}

func GetCartTotal(userID uuid.UUID) (float64, error) {
	var total float64
	query := `
		SELECT COALESCE(SUM(p.price * ci.quantity), 0)
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.user_id = $1
	`
	
	err := db.DB.Get(&total, query, userID)
	return total, err
}