package main

import (
	controller "GolangTools/controller"
	"log"
	"net/http"

	gomail "GolangTools/gomail"

	"github.com/claudiu/gocron"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/login", controller.CheckUserLogin).Methods("GET")
	router.HandleFunc("/logout", controller.Logout).Methods("GET")

	router.HandleFunc("/users", controller.InsertUser).Methods("POST")
	router.HandleFunc("/users", controller.Authenticate(controller.GetAllUsers, 1)).Methods("GET")
	router.HandleFunc("/users/{id}", controller.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", controller.Authenticate(controller.DeleteUser, 1)).Methods("DELETE")

	router.HandleFunc("/products", controller.InsertProduct).Methods("POST")
	router.HandleFunc("/products", controller.GetAllProducts).Methods("GET")
	router.HandleFunc("/products/{id}", controller.UpdateProduct).Methods("PUT")
	router.HandleFunc("/products/{id}", controller.DeleteProduct).Methods("DELETE")

	router.HandleFunc("/transactions", controller.InsertTransaction).Methods("POST")
	router.HandleFunc("/transactions", controller.GetAllTransactions).Methods("GET")
	router.HandleFunc("transactions/user", controller.GetDetailUserTransaction).Methods("GET")
	router.HandleFunc("/transactions/{id}", controller.UpdateTransaction).Methods("PUT")
	router.HandleFunc("/transactions/{id}", controller.DeleteTransaction).Methods("DELETE")

	gocron.Start()
	gocron.Every(1).Day().At("08:00").Do(gomail.SendMorningMail)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	})
	handler := corsHandler.Handler(router)
	log.Println("Starting on Port")

	err := http.ListenAndServe(":8080", handler)
	log.Fatal(err)

}
