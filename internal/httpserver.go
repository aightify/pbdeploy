package main

import (
	"log"
	"net/http"
)

func getProduct() {
	log.Print("GEtting the product")
}

func createProduct() {
	log.Print("Creating the product")
}

func main() {
	http.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {

		// if r.Method == "POST" {
		// 	createProduct()
		// 	w.Write([]byte("Product created"))
		// 	return
		// }

		if r.Method == "GET" {
			w.Write([]byte("Product listing"))
			return
		}
		// w.Write([]byte("Method not alowed"))
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	log.Print("Starting server @ port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
