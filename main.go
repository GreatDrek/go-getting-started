package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

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
	http.HandleFunc("/infodb", infomydb)
	http.HandleFunc("/deletedb", deletemydb)

	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
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

	if _, err := db.Exec("INSERT INTO ticks VALUES (now())"); err != nil {
		w.Write([]byte("Error incrementing tick: %q"))
		return
	}
}

func infomydb(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT tick FROM ticks")
	if err != nil {
		w.Write([]byte("Error reading ticks"))
		return
	}

	defer rows.Close()
	for rows.Next() {
		var tick time.Time
		if err := rows.Scan(&tick); err != nil {
			w.Write([]byte("Error scanning ticks"))
			return
		}
		w.Write([]byte(tick.String() + "\n"))
	}
}

func deletemydb(w http.ResponseWriter, r *http.Request) {
	result, err := db.Exec("DELETE FROM ticks")
	if err != nil {
		w.Write([]byte("Error delete"))
	}
}
