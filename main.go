package main

import (
	"net/http"

	"github.com/HarukiIdo/go_sample/db"
	"github.com/HarukiIdo/go_sample/router"
)

func main() {

	db, err := db.NewDB("/")
	if err != nil {
		return
	}
	mux := router.NewRouter(db)

	http.ListenAndServe(":8080", mux)
}
