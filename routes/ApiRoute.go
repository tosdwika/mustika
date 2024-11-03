package routes

import (
	"mustika/controllers"
	"mustika/middlewares"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/register", controllers.Register).Methods("POST")
	r.HandleFunc("/login", controllers.Login).Methods("POST")

	customerRouter := r.PathPrefix("/customers").Subrouter()
	customerRouter.Use(middlewares.JWTMiddleware)

	customerRouter.HandleFunc("", controllers.CreateCustomer).Methods("POST")        // Tambah customer
	customerRouter.HandleFunc("", controllers.GetCustomers).Methods("GET")           // Ambil semua customer
	customerRouter.HandleFunc("/{id}", controllers.GetCustomerByID).Methods("GET")   // Ambil customer berdasarkan ID
	customerRouter.HandleFunc("/{id}", controllers.UpdateCustomer).Methods("PUT")    // Update customer
	customerRouter.HandleFunc("/{id}", controllers.DeleteCustomer).Methods("DELETE") // Hapus customer

	orderRouter := r.PathPrefix("/orders").Subrouter()
	orderRouter.Use(middlewares.JWTMiddleware)
	orderRouter.HandleFunc("", controllers.CreateOrder).Methods("POST")        // Tambah order
	orderRouter.HandleFunc("", controllers.GetOrders).Methods("GET")           // Ambil semua orders
	orderRouter.HandleFunc("/{id}", controllers.GetOrderByID).Methods("GET")   // Ambil order berdasarkan ID
	orderRouter.HandleFunc("/{id}", controllers.UpdateOrder).Methods("PUT")    // Update order
	orderRouter.HandleFunc("/{id}", controllers.DeleteOrder).Methods("DELETE") // Hapus order

}
