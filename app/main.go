package main

import (
	"log"
	"net/http"

	"github.com/lucaslehot/MT2022_PROJ02/app/database"
	"github.com/lucaslehot/MT2022_PROJ02/app/router"
)

func main() {
	port := "8080"
	newRouter := router.NewRouter()
	if err := database.Connect(); err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}
	log.Print("\nServer started on port " + port)

	newRouter.PathPrefix("/").Handler(http.FileServer(http.Dir("./views")))

	err := http.ListenAndServe(":"+port, newRouter)
	if err != nil {
		log.Fatalf("could not serve on port %s", port)
	}
}
