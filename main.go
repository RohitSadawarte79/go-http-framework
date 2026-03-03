package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age,omitempty"`
}

func NewUser(Id int, FirstName string, LastName string, Age int) User {
	return User{
		Id:        Id,
		FirstName: FirstName,
		LastName:  LastName,
		Age:       Age,
	}
}

type Users struct {
	UsersList []User
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	users := Users{
		UsersList: make([]User, 0),
	}

	idStr := r.URL.Query().Get("id")
	if idStr != "" {

		id, err := strconv.Atoi(idStr)

		if err != nil {
			http.Error(w, "Some error occured", http.StatusBadRequest)
			return
		}

		user := User{
			Id:        id,
			FirstName: "Some",
			LastName:  "Name",
		}

		users.UsersList = append(users.UsersList, user)

	} else {
		user1 := NewUser(1, "Alice", "Smith", 30)
		user2 := NewUser(2, "Bob", "Jones", 25)
		users.UsersList = append(users.UsersList, user1)
		users.UsersList = append(users.UsersList, user2)

	}

	JSON(w, http.StatusOK, users)

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("Creating user: %+v\n", user)

	JSON(w, http.StatusCreated, user)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	users := Users{
		UsersList: make([]User, 0),
	}

	user := User{
		Id:        id,
		FirstName: "Some",
		LastName:  "Name",
	}

	users.UsersList = append(users.UsersList, user)

	JSON(w, http.StatusOK, users)

}

func main() {
	router := NewRouter()

	corsConfig := CORSConfig{
		AllowedOrigins: map[string]bool{
			"http://localhost:3000": true,
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}

	corsMiddleware := NewCORS(corsConfig)

	router.HandleFunc("GET", "/user", GetUser)
	router.HandleFunc("POST", "/user", CreateUser)
	router.HandleFunc("GET", "/user/:id", GetUserByID)

	stack := Chain(corsMiddleware, Recovery, Logger, RequestId)(router)

	fmt.Println("Listening on port 8080:", "http://localhost:8080")
	err := http.ListenAndServe(":8080", stack)

	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
