package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Data struct {
	PageTitle string
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("src/templates/layouts/layout.html")
	if err != nil {
		log.Fatal("Error parsing layout.html")
	}
	data := Data{
		PageTitle: "omni-association",
	}
	tmpl.Execute(w, data)
}

func requestSegment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]
	page := vars["page"]

	fmt.Fprintf(w, "you've requested the book: %s on page %s\n", title, page)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/books/{title}/page/{page}", requestSegment)
	router.HandleFunc("/", home)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	domain := host + ":" + port

	http.ListenAndServe(domain, router)
}
