package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"

	//"time"
	"strconv"

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
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS weather (city varchar(80), temp_lo int)"); err != nil {
		w.Write([]byte("Error add data base info"))
		return
	} else {
		w.Write([]byte("Add data base info"))
	}

	if _, err := db.Exec("INSERT INTO weather VALUES ('San Francisco', $1)", rand.Int()); err != nil {
		w.Write([]byte("Error incrementing tick: %q"))
		return
	} else {
		w.Write([]byte("Add info"))
	}
}

type weather struct {
	city    string
	temp_lo int
}

func (w weather) String() string {
	return w.city + " : " + strconv.Itoa(w.temp_lo)
}

func infomydb(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM weather")
	if err != nil {
		w.Write([]byte("Error reading ticks"))
		return
	}

	defer rows.Close()
	for rows.Next() {
		//var tick time.Time
		var weatherTest weather
		if err := rows.Scan(&weatherTest.city, &weatherTest.temp_lo); err != nil {
			w.Write([]byte("Error scanning ticks"))
			return
		}
		w.Write([]byte(weatherTest.String() + "\n"))
	}
}

func deletemydb(w http.ResponseWriter, r *http.Request) {
	_, err := db.Exec("DROP TABLE weather")
	if err != nil {
		w.Write([]byte("Error delete"))
	}
}

//"DROP TABLE ticks"
