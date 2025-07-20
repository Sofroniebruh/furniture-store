package services

import (
	"encoding/json"
	"furniture-store-backend/db"
	"furniture-store-backend/models"
	"github.com/google/uuid"
	"log"
	"net/http"
)

func CreateColor(w http.ResponseWriter, r *http.Request) {
	var colorDto models.Color
	var color models.Color
	var exists bool
	err := json.NewDecoder(r.Body).Decode(&colorDto)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Bad Request",
		})
		return
	}

	if colorDto.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Color name is required",
		})
		return
	}

	err = db.DB.QueryRow("SELECT EXISTS(SELECT FROM colors WHERE name ILIKE $1)", colorDto.Name).Scan(&exists)

	if exists {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Color already exists",
		})
		return
	}

	err = db.DB.QueryRow(`INSERT INTO colors (name) VALUES ($1) RETURNING id`, colorDto.Name).Scan(&color.ID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	color.Name = colorDto.Name

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]models.Color{
		"created": color,
	})
}

func GetAllColors(w http.ResponseWriter, r *http.Request) {
	var rows []models.Color
	var colors []models.Color

	err := db.DB.Select(&rows, "SELECT * FROM colors")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	for _, row := range rows {
		colors = append(colors, models.Color{
			ID:   row.ID,
			Name: row.Name,
		})
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string][]models.Color{
		"colors": colors,
	})
}

func DeleteColor(w http.ResponseWriter, r *http.Request) {
	var colorDto models.Color
	var deletedColor models.Color

	err := json.NewDecoder(r.Body).Decode(&colorDto)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Bad Request",
		})
		return
	}

	if colorDto.ID == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Color id is required",
		})
		return
	}

	err = db.DB.QueryRow("DELETE FROM colors WHERE id = $1 RETURNING name", colorDto.ID).Scan(&deletedColor.Name)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to delete the color",
		})
		return
	}

	deletedColor.ID = colorDto.ID

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]models.Color{
		"deleted": deletedColor,
	})
}

func UpdateColor(w http.ResponseWriter, r *http.Request) {
	var colorDto models.Color
	var updatedColor models.Color

	err := json.NewDecoder(r.Body).Decode(&colorDto)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Bad Request",
		})
		return
	}

	if colorDto.ID == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Color id is required",
		})
		return
	}
	if colorDto.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Color name is required",
		})
		return
	}

	err = db.DB.QueryRow("UPDATE colors SET name = $1 WHERE id = $2 RETURNING name", colorDto.Name, colorDto.ID).Scan(&updatedColor.Name)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to update the color",
		})
		return
	}

	updatedColor.ID = colorDto.ID

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]models.Color{
		"updated": updatedColor,
	})
}
