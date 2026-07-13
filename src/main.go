package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/yasseraitnasser/omni-association/src/auth"
	"github.com/yasseraitnasser/omni-association/src/database"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc
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

func Method(m string) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != m {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			f(w, r)
		}
	}
}

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/login", Chain(auth.Login, Method("POST")))
	router.HandleFunc("/", home)

	var err error
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = database.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer database.DB.Close()

	database.AddAdminUser()

	port := os.Getenv("SERVER_PORT")
	host := os.Getenv("SERVER_HOST")
	domain := host + ":" + port
	log.Printf("Listening on: %s\n", domain)
	http.ListenAndServe(domain, router)
}
