package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"mustika/config"
	"mustika/models"

	"github.com/gorilla/mux"
)

// CreateCustomer godoc
// @Summary Membuat customer baru
// @Description Menambahkan customer baru ke dalam database
// @Tags customer
// @Accept  json
// @Produce  json
// @Param customer body models.Customer true "Customer Data"
// @Success 201 {object} models.Customer
// @Failure 500 {string} string "Internal Server Error"
// @Router /customer [post]
func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	json.NewDecoder(r.Body).Decode(&customer)

	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()

	if err := config.DB.Create(&customer).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

// GetCustomers godoc
// @Summary Mendapatkan daftar customer
// @Description Mengambil daftar semua customer dengan paginasi
// @Tags customer
// @Accept  json
// @Produce  json
// @Param page query int false "Halaman"
// @Param limit query int false "Jumlah item per halaman"
// @Success 200 {array} models.Customer
// @Failure 404 {string} string "No data found"
// @Router /customers [get]
func GetCustomers(w http.ResponseWriter, r *http.Request) {
	var customers []models.Customer

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

	if err := config.DB.Limit(paginateRequest.Limit).Offset(offset).Find(&customers).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(customers) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No data found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(customers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetCustomerByID godoc
// @Summary Mendapatkan detail customer
// @Description Mengambil informasi customer berdasarkan ID
// @Tags customer
// @Accept  json
// @Produce  json
// @Param id path int true "Customer ID"
// @Success 200 {object} models.Customer
// @Failure 404 {string} string "Customer not found"
// @Router /customer/{id} [get]
func GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var customer models.Customer

	if err := config.DB.First(&customer, id).Error; err != nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(customer)
}

// UpdateCustomer godoc
// @Summary Memperbarui data customer
// @Description Memperbarui informasi customer yang sudah ada
// @Tags customer
// @Accept  json
// @Produce  json
// @Param id path int true "Customer ID"
// @Param customer body models.Customer true "Customer Data"
// @Success 200 {object} models.Customer
// @Failure 404 {string} string "Customer not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /customer/{id} [put]
func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var customer models.Customer

	if err := config.DB.First(&customer, id).Error; err != nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	json.NewDecoder(r.Body).Decode(&customer)
	customer.UpdatedAt = time.Now()

	if err := config.DB.Save(&customer).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(customer)
}

// DeleteCustomer godoc
// @Summary Menghapus data customer
// @Description Menghapus customer berdasarkan ID
// @Tags customer
// @Accept  json
// @Produce  json
// @Param id path int true "Customer ID"
// @Success 200 {string} string "Customer berhasil dihapus"
// @Failure 404 {string} string "Customer not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /customer/{id} [delete]
func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var customer models.Customer

	if err := config.DB.First(&customer, id).Error; err != nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	if err := config.DB.Delete(&customer).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Customer berhasil dihapus"})
}
