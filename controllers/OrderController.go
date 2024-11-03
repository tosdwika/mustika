package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"mustika/config"
	"mustika/models"

	"github.com/gorilla/mux"
)

// CreateOrder godoc
// @Summary Membuat pesanan baru
// @Description Menambahkan pesanan baru ke dalam database
// @Tags order
// @Accept  json
// @Produce  json
// @Param order body models.Order true "Order Data"
// @Success 201 {object} models.Order
// @Failure 400 {string} string "Invalid request payload atau status tidak valid"
// @Failure 500 {string} string "Internal Server Error"
// @Router /orders [post]
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	validStatuses := map[string]struct{}{
		"pending":   {},
		"shipped":   {},
		"delivered": {},
		"cancelled": {},
	}

	if _, valid := validStatuses[order.Status]; !valid {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	order.OrderDate = time.Now()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	if err := config.DB.Create(&order).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// GetOrders godoc
// @Summary Mendapatkan daftar pesanan
// @Description Mengambil daftar semua pesanan dengan paginasi
// @Tags order
// @Accept  json
// @Produce  json
// @Param page query int false "Halaman"
// @Param limit query int false "Jumlah item per halaman"
// @Success 200 {array} models.Order
// @Failure 404 {string} string "No data found"
// @Router /orders [get]
func GetOrders(w http.ResponseWriter, r *http.Request) {
	var orders []models.Order

	type PaginateRequest struct {
		Page  int `json:"page"`
		Limit int `json:"limit"`
	}

	var paginateRequest PaginateRequest

	paginateRequest.Page = 1
	paginateRequest.Limit = 10

	if r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&paginateRequest); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
	}

	if paginateRequest.Page <= 0 {
		paginateRequest.Page = 1
	}
	if paginateRequest.Limit <= 0 {
		paginateRequest.Limit = 10
	}

	offset := (paginateRequest.Page - 1) * paginateRequest.Limit

	if err := config.DB.Limit(paginateRequest.Limit).Offset(offset).Find(&orders).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No data found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetOrderByID godoc
// @Summary Mendapatkan detail pesanan
// @Description Mengambil informasi pesanan berdasarkan ID
// @Tags order
// @Accept  json
// @Produce  json
// @Param id path int true "Order ID"
// @Success 200 {object} models.Order
// @Failure 404 {string} string "Order not found"
// @Router /orders/{id} [get]
func GetOrderByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var order models.Order

	if err := config.DB.First(&order, id).Error; err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(order)
}

// UpdateOrder godoc
// @Summary Memperbarui data pesanan
// @Description Memperbarui informasi pesanan yang sudah ada
// @Tags order
// @Accept  json
// @Produce  json
// @Param id path int true "Order ID"
// @Param order body models.Order true "Order Data"
// @Success 200 {object} models.Order
// @Failure 400 {string} string "Invalid request payload"
// @Failure 404 {string} string "Order not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /orders/{id} [put]
func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var order models.Order

	if err := config.DB.First(&order, id).Error; err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order.UpdatedAt = time.Now()

	if err := config.DB.Save(&order).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(order)
}

// DeleteOrder godoc
// @Summary Menghapus data pesanan
// @Description Menghapus pesanan berdasarkan ID
// @Tags order
// @Accept  json
// @Produce  json
// @Param id path int true "Order ID"
// @Success 200 {string} string "Order berhasil dihapus"
// @Failure 404 {string} string "Order not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /orders/{id} [delete]
func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var order models.Order

	if err := config.DB.First(&order, id).Error; err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	if err := config.DB.Delete(&order).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Order berhasil dihapus"})
}
