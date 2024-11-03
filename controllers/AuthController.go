package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"mustika/config"
	"mustika/models"
	"mustika/utils"

	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// @Summary Registrasi pengguna baru
// @Description Membuat pengguna baru dengan menyimpan username dan password terenkripsi
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body models.User true "User Data"
// @Success 201 {object} models.User
// @Failure 500 {string} string "Internal Server Error"
// @Router /register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := config.DB.Create(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login godoc
// @Summary Login pengguna
// @Description Mengautentikasi pengguna dan menghasilkan token JWT jika berhasil
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body models.User true "User Data"
// @Success 200 {object} map[string]string "token JWT"
// @Failure 401 {string} string "User not found atau password salah"
// @Failure 500 {string} string "Internal Server Error"
// @Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	var dbUser models.User
	if err := config.DB.Where("username = ?", user.Username).First(&dbUser).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(int(dbUser.ID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
