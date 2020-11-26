package main
import (
	"net/http"
	"log"
	"html/template"
	_ "github.com/lib/pq"
	"database/sql"
	"os"
	"fmt"
)
type Book struct{
	Name string
	Year string
	Length string
}
const (
	DB_USER = "postgres"
	DB_PASSWORD = "0000"
	DB_NAME = "lab"
)
func dbConnect() error {
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME))
	if err != nil {
		return err
	}
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS books (book_name text,book_year text,book_length text)"); err != nil {
		return err
	}
	return nil
}
func dbAddBook(name, year, length string) error {
	sqlstmt := "INSERT INTO books VALUES ($1, $2, $3)"
	_, err := db.Exec(sqlstmt, name, year, length)
	if err != nil {
		return err
	}
	return nil
}
func dbGetBooks() ([]Book, error) {
	var books []Book
	stmt, err := db.Prepare("SELECT book_name, book_year, book_length FROM books")
	if err != nil {
		return books, err
	}
	res, err := stmt.Query()
	if err != nil {
		return books, err
	}
	var tempBook Book
	for res.Next() {
		err = res.Scan(&tempBook.Name, &tempBook.Year, &tempBook.Length)
		if err != nil {
			return books, err
		}
		books = append(books, tempBook)
	}
	return books, err
}

var db *sql.DB
func rollHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("simple_list.html")
		if err != nil {
			log.Fatal(err)
		}
		books, err := dbGetBooks()
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(w, books)
	}
}
func addBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("simple_form.html")
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		name := r.Form.Get("name")
		year := r.Form.Get("year")
		length := r.Form.Get("length")
		err := dbAddBook(name, year, length)
		if err != nil {
			log.Fatal(err)
		}
	}
}
func GetPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "4747"
		fmt.Println(port)
	}
	return ":" + port
}
func main() {
	err := dbConnect()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", rollHandler)
	http.HandleFunc("/add", addBookHandler)
	log.Fatal(http.ListenAndServe(GetPort(), nil))
}
