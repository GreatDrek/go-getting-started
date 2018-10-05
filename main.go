package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"database/sql"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

func main() {
	var err error
	port := os.Getenv("PORT")

	if port == "" {
		//log.Fatal("$PORT must be set")
		port = "5000"
	}

	http.HandleFunc("/", mainPage)
	http.HandleFunc("/users", users)
	http.HandleFunc("/db", mydb)

	db, err = sql.Open("postgres", os.Getenv("postgres://vvyrgvcitdhpbd:7b6cf6526b1839736ee8edab5a8cbac7bb9590d59656d172356584a4d88447e5@ec2-54-217-245-26.eu-west-1.compute.amazonaws.com:5432/d4qm0eh3iljgrt"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	defer db.Close()

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Error:", err)
	}
}

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	urlP := r.URL.Path
	b := []byte(urlP)
	w.Write(b[1:])
}

func users(w http.ResponseWriter, r *http.Request) {
	userSlice := []User{User{"One", "Two"}, User{"Three", "Four"}}
	js, err := json.Marshal(userSlice)
	if err != nil {
		log.Fatal("Error:", err)
	}
	w.Write(js)
}

func mydb(w http.ResponseWriter, r *http.Request) {
	//_, err := db.Exec("insert into Products (model, company, price) values ('iPhone X', $1, $2)",	"Apple", 72000)
	//if err != nil {
	//	log.Fatal("Error db:", err)
	//}
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS ticks (tick timestamp)"); err != nil {
		w.Write([]byte("Error add data base info"))
		return
	} else {
		w.Write([]byte("Add data base info"))
	}
}
