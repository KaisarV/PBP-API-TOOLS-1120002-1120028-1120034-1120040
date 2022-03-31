package controllers

import (
	model "GolangTools/model"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	db := connect()

	defer db.Close()
	var response model.ProductsResponse

	query := "SELECT * FROM products"
	id := r.URL.Query()["id"]
	if id != nil {
		query += " WHERE id = " + id[0]
	}

	rows, err := db.Query(query)

	if err != nil {
		response.Status = 400
		response.Message = err.Error()
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	var product model.Product
	var products []model.Product

	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			log.Println(err.Error())
		} else {
			products = append(products, product)
		}
	}

	if len(products) != 0 {
		response.Status = 200
		response.Message = "Success Get Data"
		response.Data = products
	} else {
		response.Status = 400
		response.Message = "Data Not Found"
		w.WriteHeader(400)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	var response model.ErrorResponse

	if err != nil {
		response.Status = 400
		response.Message = "Error Parsing Data"
		w.WriteHeader(400)
		log.Println(err.Error())
		return
	}

	vars := mux.Vars(r)
	productId := vars["id"]
	query, errQuery := db.Exec(`DELETE FROM transactions WHERE ProductId = ?;`, productId)
	query, errQuery = db.Exec(`DELETE FROM products WHERE id = ?;`, productId)
	RowsAffected, err := query.RowsAffected()

	if RowsAffected == 0 {
		response.Status = 400
		response.Message = "Product not found"
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if errQuery == nil {
		response.Status = 200
		response.Message = "Success Delete Data"
		w.WriteHeader(200)
	} else {
		response.Status = 400
		response.Message = "Error Delete Data"
		w.WriteHeader(400)
		log.Println(errQuery.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	var response model.ProductResponse
	err := r.ParseForm()

	if err != nil {
		response.Status = 400
		response.Message = "Error Parsing Data"
		w.WriteHeader(400)
		log.Println(err.Error())
		return
	}

	var product model.Product
	product.Name = r.Form.Get("name")
	product.Price, _ = strconv.Atoi(r.Form.Get("price"))

	log.Println(product.Name)
	log.Println(product.Price)

	res, errQuery := db.Exec("INSERT INTO products (name, price) VALUES (?,?)", product.Name, product.Price)

	id, err := res.LastInsertId()

	if errQuery == nil {
		response.Status = 200
		response.Message = "Success"
		product.ID = int(id)
		response.Data = product
	} else {
		response.Status = 400
		response.Message = "Error Insert Data"
		w.WriteHeader(400)
		log.Println(errQuery.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {

	db := connect()
	defer db.Close()

	err := r.ParseForm()
	var response model.ProductResponse

	if err != nil {
		response.Status = 400
		response.Message = "Error Parsing Data"
		w.WriteHeader(400)
		log.Println(err.Error())
		return
	}

	vars := mux.Vars(r)
	productId := vars["id"]

	var product model.Product
	product.Name = r.Form.Get("name")
	product.Price, _ = strconv.Atoi(r.Form.Get("price"))

	rows, _ := db.Query(`SELECT * FROM products WHERE ID = ?;`, productId)
	var prevDatas []model.Product
	var prevData model.Product

	for rows.Next() {
		if err := rows.Scan(&prevData.ID, &prevData.Name, &prevData.Price); err != nil {
			log.Println(err.Error())
		} else {
			prevDatas = append(prevDatas, prevData)
		}
	}

	if len(prevDatas) > 0 {
		if product.Name == "" {
			product.Name = prevDatas[0].Name
		}
		if product.Price == 0 {
			product.Price = prevDatas[0].Price
		}

		_, errQuery := db.Exec(`UPDATE products SET name = ?, price = ? WHERE id = ?;`, product.Name, product.Price, productId)

		if errQuery == nil {
			response.Status = 200
			response.Message = "Success Update Data"
			id, _ := strconv.Atoi(productId)
			product.ID = id
			response.Data = product
			w.WriteHeader(200)
		} else {
			response.Status = 400
			response.Message = "Error Update Data"
			w.WriteHeader(400)
			log.Println(errQuery)
		}
	} else {
		response.Status = 400
		response.Message = "Data Not Found"
		w.WriteHeader(400)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
