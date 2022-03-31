package controllers

import (
	model "GolangTools/model"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	db := connect()

	defer db.Close()
	var response model.TransactionsResponse

	query := "SELECT * FROM transactions"
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

	var transaction model.Transaction
	var transactions []model.Transaction

	for rows.Next() {
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.ProductId, &transaction.Quantity); err != nil {
			log.Println(err.Error())
		} else {
			transactions = append(transactions, transaction)
		}
	}

	if len(transactions) != 0 {
		response.Status = 200
		response.Message = "Succes Get Data"
		response.Data = transactions
	} else if response.Message == "" {
		response.Status = 400
		response.Message = "Data Not Found"
		w.WriteHeader(400)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
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
	transactionId := vars["id"]
	query, errQuery := db.Exec(`DELETE FROM transactions WHERE id = ?;`, transactionId)
	RowsAffected, _ := query.RowsAffected()

	if RowsAffected == 0 {
		response.Status = 400
		response.Message = "Transaction not found"
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
		response.Message = "Failed Delete Data"
		w.WriteHeader(400)
		log.Println(errQuery.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	var response model.TransactionResponse
	err := r.ParseForm()

	if err != nil {
		response.Status = 400
		response.Message = "Error Parsing Data"
		w.WriteHeader(400)
		log.Println(err.Error())
		return
	}

	var transaction model.Transaction

	transaction.UserID, _ = strconv.Atoi(r.Form.Get("userid"))
	transaction.ProductId, _ = strconv.Atoi(r.Form.Get("productid"))
	transaction.Quantity, _ = strconv.Atoi(r.Form.Get("qty"))

	rows, _ := db.Query("SELECT * FROM products WHERE id=?", transaction.ProductId)

	if err != nil {
		response.Status = 400
		response.Message = err.Error()
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	i := 0
	for rows.Next() {
		i++
	}

	if i == 0 {
		_, err = db.Exec("INSERT INTO products (id) VALUES (?)", transaction.ProductId)

		if err != nil {
			response.Status = 400
			response.Message = err.Error()
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}
	res, errQuery := db.Exec("INSERT INTO transactions (userid, productid, quantity) VALUES (?,?,?)", transaction.UserID, transaction.ProductId, transaction.Quantity)

	id, _ := res.LastInsertId()

	if errQuery == nil {
		response.Status = 200
		response.Message = "Success"
		transaction.ID = int(id)
		response.Data = transaction
	} else {
		response.Status = 400
		response.Message = "Error Insert Data"
		w.WriteHeader(400)
		log.Println(errQuery.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	var response model.TransactionResponse

	if err != nil {
		response.Status = 400
		response.Message = "Error Parsing Data"
		w.WriteHeader(400)
		log.Println(err.Error())
		return
	}

	vars := mux.Vars(r)
	transactionId := vars["id"]

	var transaction model.Transaction
	transaction.UserID, _ = strconv.Atoi(r.Form.Get("userid"))
	transaction.ProductId, _ = strconv.Atoi(r.Form.Get("productid"))
	transaction.Quantity, _ = strconv.Atoi(r.Form.Get("qty"))

	rows, _ := db.Query("SELECT * FROM transactions WHERE id = ?", transactionId)
	var prevDatas []model.Transaction
	var prevData model.Transaction

	for rows.Next() {
		if err := rows.Scan(&prevData.ID, &prevData.UserID, &prevData.ProductId, &prevData.Quantity); err != nil {
			log.Println(err.Error())
		} else {
			prevDatas = append(prevDatas, prevData)
		}
	}

	if len(prevDatas) > 0 {
		if transaction.UserID == 0 {
			transaction.UserID = prevDatas[0].UserID
		}
		if transaction.ProductId == 0 {
			transaction.ProductId = prevDatas[0].ProductId
		}
		if transaction.Quantity == 0 {
			transaction.Quantity = prevDatas[0].Quantity
		}

		_, errQuery := db.Exec(`UPDATE transactions SET UserID = ?, ProductID = ?, Quantity = ? WHERE id = ?;`, transaction.UserID, transaction.ProductId, transaction.Quantity, transactionId)

		if errQuery == nil {
			response.Status = 200
			response.Message = "Success Update Data"
			id, _ := strconv.Atoi(transactionId)
			transaction.ID = id
			response.Data = transaction
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

func GetDetailUserTransaction(w http.ResponseWriter, r *http.Request) {

	db := connect()
	defer db.Close()

	var transactionDetails []model.TransactionDetail

	query := "SELECT t.ID , u.ID, u.Name, u.Age, u.Address, p.ID, p.Name, p.Price, t.Quantity FROM transactions t JOIN users u ON t.UserId = u.ID JOIN products p ON p.ID = t.ProductID"

	id := r.URL.Query()["id"]
	if id != nil {
		query += " WHERE u.id = " + id[0]
	}

	rows, err := db.Query(query)

	var response model.TransactionDetailsResponse

	if err != nil {
		response.Status = 400
		response.Message = err.Error()
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	var transactionDetail model.TransactionDetail
	var user model.User
	var product model.Product

	for rows.Next() {
		if err := rows.Scan(&transactionDetail.ID, &user.ID, &user.Name, &user.Age, &user.Address, &product.ID, &product.Name, &product.Price, &transactionDetail.Quantity); err != nil {
			log.Println(err.Error())
		} else {
			transactionDetail.User = user
			transactionDetail.Product = product
			transactionDetails = append(transactionDetails, transactionDetail)
		}
	}

	if len(transactionDetails) != 0 {
		response.Status = 200
		response.Message = "Success Get Data"
		response.Data = transactionDetails
	} else {
		response.Status = 400
		response.Message = "Error Get Data"
		w.WriteHeader(400)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
