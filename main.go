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

	http.HandleFunc("/sdb", sdb)
	http.HandleFunc("/sdb2", sdb2)

	//connStr := "user=postgres password=37352410 dbname=postgres sslmode=disable"
	//db, err = sql.Open("postgres", connStr)

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
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS weather (city int, temp_lo int)"); err != nil {
		w.Write([]byte("Error add data base info"))
		return
	} else {
		w.Write([]byte("Add data base info"))
	}

	for i := 0; i < 10000; i++ {
		if _, err := db.Exec("INSERT INTO weather VALUES ($1, $2)", rand.Int31(), rand.Int31()); err != nil {
			w.Write([]byte("Error incrementing tick: %q"))
			return
		} else {
			//w.Write([]byte("Add info"))
		}
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

func sdb(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT city, temp_lo FROM weather")
	if err != nil {
		w.Write([]byte("Error reading ticks"))
		return
	}

	defer rows.Close()

	for rows.Next() {
		//var tick time.Time
		var numb int
		var numb2 int
		if err := rows.Scan(&numb, &numb2); err != nil {
			w.Write([]byte("Error scanning ticks"))
			return
		}
		if numb == 513934398 {
			w.Write([]byte(strconv.Itoa(numb) + " : " + strconv.Itoa(numb2) + "\n"))
		}

		w.Write([]byte("_"))
	}
}

func sdb2(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT city, temp_lo FROM weather WHERE city = 513934398")
	if err != nil {
		w.Write([]byte("Error reading ticks"))
		return
	}

	defer rows.Close()

	var numb int
	var numb2 int

	for rows.Next() {
		if err := rows.Scan(&numb, &numb2); err != nil {
			w.Write([]byte("Error scanning ticks" + err.Error()))
			return
		}
		w.Write([]byte("_"))
		w.Write([]byte(strconv.Itoa(numb) + " : " + strconv.Itoa(numb2) + "\n"))
	}

	//	for rows.Next() {
	//		//var tick time.Time
	//		var numb int
	//		var numb2 int
	//		if err := rows.Scan(&numb, &numb2); err != nil {
	//			w.Write([]byte("Error scanning ticks"))
	//			return
	//		}
	//		if numb == 513934398 {
	//			w.Write([]byte(strconv.Itoa(numb) + " : " + strconv.Itoa(numb2) + "\n"))
	//		}
	//	}
}

//"DROP TABLE ticks"
