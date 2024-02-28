package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func seedAccount(store Storage, fname, lname, pw string) *Account {
	acc, err := NewAccount(fname, lname, pw)
	if err != nil {
		log.Fatal(err)
	}

	if err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	fmt.Println("new account => ", acc.Number)

	return acc
}

func seedAccounts(s Storage) {
	seedAccount(s, "anthony", "GG", "hunter88888")
}

var logger = logrus.New()

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"uri":    r.RequestURI,
			"remote": r.RemoteAddr,
		}).Info("HTTP request")
		fmt.Println("HELLO")
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
func main() {
	seed := flag.Bool("seed", true, "seed the db")
	flag.Parse()

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Println("seeding the database")
		seedAccounts(store)
	}

	mux.NewRouter().Handle("/", LoggingMiddleware(http.HandlerFunc(handler)))
	server := NewAPIServer(":3000", store)
	server.Run()
}
