package controllers

import (
	gomail "GolangTools/gomail"
	model "GolangTools/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	// "github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	var response model.ErrorResponse

	err := r.ParseForm()
	if err != nil {
		response.Status = 400
		response.Message = "Error Parsing Data"
		w.WriteHeader(400)
		log.Println(err.Error())
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	if email == "" {
		response.Status = 400
		response.Message = "Please Input Email"
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if password == "" {
		response.Status = 400
		response.Message = "Please Input Password"
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	rows, err := db.Query("SELECT email, password FROM users WHERE email= ?", email)

	if err != nil {
		response.Status = 400
		response.Message = err.Error()
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	var user model.User
	var users []model.User

	for rows.Next() {
		if err := rows.Scan(&user.Email, &user.Password); err != nil {
			log.Println(err.Error())
		} else {
			users = append(users, user)
		}
	}

	if users[0].Password == password {
		generateToken(w, user.ID, user.Email, user.UserType)
		response.Status = 200
		response.Message = "Login Success"
		w.Header().Set("Content-Type", "application/json")
	} else {
		response.Status = 400
		response.Message = "Login Failed"
		w.WriteHeader(400)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	var response model.UsersResponse

	query := "SELECT * FROM users"
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

	var user model.User
	var users []model.User

	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.Email, &user.Password, &user.UserType); err != nil {
			log.Println(err.Error())
		} else {
			users = append(users, user)
		}
	}

	if len(users) != 0 {
		response.Status = 200
		response.Message = "Success Get Data"
		response.Data = users
	} else if response.Message == "" {
		response.Status = 400
		response.Message = "Data Not Found"
		w.WriteHeader(400)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
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
	userId := vars["id"]
	query, errQuery := db.Exec(`DELETE FROM users WHERE id = ?;`, userId)
	RowsAffected, _ := query.RowsAffected()

	if RowsAffected == 0 {
		response.Status = 400
		response.Message = "User not found"
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

func InsertUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	var response model.UserResponse

	if err != nil {
		response.Status = 400
		response.Message = "Error Parsing Data"
		w.WriteHeader(400)
		log.Println(err.Error())
		return
	}

	var user model.User

	user.Name = r.Form.Get("name")
	user.Age, _ = strconv.Atoi(r.Form.Get("age"))
	user.Address = r.Form.Get("address")
	user.Email = r.Form.Get("email")
	user.Password = r.Form.Get("password")

	res, errQuery := db.Exec("INSERT INTO users (name, age, address, email, password) VALUES (?,?,?,?,?)", user.Name, user.Age, user.Address, user.Email, user.Password)
	id, _ := res.LastInsertId()

	if errQuery == nil {
		gomail.SendMail(user.Email, user.Name)
		go gomail.SendPromoMail(user.Email, user.Name) //go routine(mengirim 2 email asynchronous)
		response.Status = 200
		response.Message = "Success"
		user.ID = int(id)
		response.Data = user
	} else {
		response.Status = 400
		response.Message = "Error Insert Data"
		w.WriteHeader(400)
		log.Println(errQuery.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	var response model.UserResponse

	if err != nil {
		response.Status = 400
		response.Message = "Error Parsing Data"
		w.WriteHeader(400)
		log.Println(err.Error())
		return
	}

	vars := mux.Vars(r)
	userId := vars["id"]

	var user model.User
	user.Name = r.Form.Get("name")
	user.Age, _ = strconv.Atoi(r.Form.Get("age"))
	user.Address = r.Form.Get("address")
	user.Email = r.Form.Get("email")
	user.Password = r.Form.Get("password")

	rows, _ := db.Query("SELECT * FROM users WHERE id = ?", userId)
	var prevDatas []model.User
	var prevData model.User

	for rows.Next() {
		if err := rows.Scan(&prevData.ID, &prevData.Name, &prevData.Age, &prevData.Address, &prevData.Email, &prevData.Password, &prevData.UserType); err != nil {
			log.Println(err.Error())
		} else {
			prevDatas = append(prevDatas, prevData)
		}
	}

	if len(prevDatas) > 0 {
		if user.Name == "" {
			user.Name = prevDatas[0].Name
		}
		if user.Age == 0 {
			user.Age = prevDatas[0].Age
		}
		if user.Address == "" {
			user.Address = prevDatas[0].Address
		}
		if user.Email == "" {
			user.Email = prevDatas[0].Email
		}
		if user.Password == "" {
			user.Password = prevDatas[0].Password
		}

		_, errQuery := db.Exec(`UPDATE users SET name = ?, age = ?, address = ?, email = ?, password = ? WHERE id = ?;`, user.Name, user.Age, user.Address, user.Email, user.Password, userId)

		if errQuery == nil {
			response.Status = 200
			response.Message = "Success Update Data"
			id, _ := strconv.Atoi(userId)
			user.ID = id
			response.Data = user
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

//Auth
func CheckUserLogin(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	var response model.ErrorResponse
	err := r.ParseForm()

	if err != nil {
		response.Status = 400
		response.Message = "Error Parsing Data"
		w.WriteHeader(400)
		log.Println(err.Error())
		return
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	rows, _ := db.Query("SELECT * FROM users WHERE Email = ? AND password = ?", email, password)

	fmt.Println(email)
	var user model.User
	var users []model.User

	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.Email, &user.Password, &user.UserType); err != nil {
			log.Println(err.Error())
		} else {
			users = append(users, user)
		}
	}

	if len(users) == 0 {
		response.Status = 200
		response.Message = "Login Failed"

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		generateToken(w, user.ID, user.Email, user.UserType)

		response.Status = 200
		response.Message = "Login Success"

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		// //set redis
		// rdb := redis.NewClient(&redis.Options{
		//     Addr:     "localhost:6379",
		//     Password: "", // no password set
		//     DB:       0,  // use default DB
		// })
		// SetRedis(rdb, "emailUser", user.Email, 0)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// reset token yang dikirim, token yang lama ditimpa
	resetUserToken(w)

	var response model.UserResponse
	response.Status = 200
	response.Message = "Logout Success"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendUnAuthorizedResponse(w http.ResponseWriter) {
	var response model.ErrorResponse

	response.Status = 401
	response.Message = "Unauthorized Access"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
