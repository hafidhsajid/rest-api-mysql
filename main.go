package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Customer struct {
	ID         string `json:"customerId"`
	Title      string `json:"Name"`
	Email      string `json:"Email"`
	Phone      string `json:"Phone"`
	Address    string `json:"Address"`
	Created_at string `json:"Created_at"`
}

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "admin:password@tcp(127.0.0.1:3306)/db")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/customer", getCustomers).Methods("GET")
	router.HandleFunc("/posts", getCustomers).Methods("GET")
	router.HandleFunc("/posts", createCustomer).Methods("POST")
	router.HandleFunc("/posts/{id}", getCustomer).Methods("GET")
	router.HandleFunc("/posts/{id}", updateCustomer).Methods("PUT")
	router.HandleFunc("/posts/{id}", deleteCustomer).Methods("DELETE")

	fmt.Println("starting web server at http://localhost:8080/")
	// http.ListenAndServe(":8080", nil)

	fmt.Println(http.ListenAndServe(":8080", router))
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var posts []Customer

	result, err := db.Query("SELECT * from customer")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var post Customer
		err := result.Scan(&post.ID, &post.Title, &post.Email, &post.Phone, &post.Address, &post.Created_at)
		if err != nil {
			panic(err.Error())
		}
		posts = append(posts, post)
	}

	json.NewEncoder(w).Encode(posts)
}
func getCustomers_1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var posts []Customer

	result, err := db.Query("SELECT * from posts")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var post Customer
		err := result.Scan(&post.ID, &post.Title)
		if err != nil {
			panic(err.Error())
		}
		posts = append(posts, post)
	}

	json.NewEncoder(w).Encode(posts)
}
func createCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stmt, err := db.Prepare("INSERT INTO posts(title) VALUES(?)")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	title := keyVal["title"]

	_, err = stmt.Exec(title)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New post was created")
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	result, err := db.Query("SELECT id, title FROM posts WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var post Customer
	for result.Next() {
		err := result.Scan(&post.ID, &post.Title)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(post)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE posts SET title = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newTitle := keyVal["title"]

	_, err = stmt.Exec(newTitle, params["id"])
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "Customer with ID = %s was updated", params["id"])
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM posts WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "Customer with ID = %s was deleted", params["id"])
}
