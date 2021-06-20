package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/adjust/rmq/v3"
	"github.com/lucaslehot/MT2022_PROJ02/app/database"
	"github.com/lucaslehot/MT2022_PROJ02/app/router"
)

type Task struct {
	order   string
	user_id int
}

const dwldPath = "./tmp"

func main() {
	port := "8080"
	newRouter := router.NewRouter()

	// Setting up Redis connection
	connection, err := rmq.OpenConnection("message_broker", "tcp", "localhost:6379", 1, nil)
	taskQueue, err := connection.OpenQueue("tasks")

	err = database.Connect()
	if err != nil {
		log.Fatalf("Impossible de se connecter à la bdd: %v", err)
	}

	log.Print("\nServer started on port " + port)

	newRouter.PathPrefix("/files/").Handler(http.StripPrefix("/files/",
		http.FileServer(http.Dir(dwldPath))))

	err = http.ListenAndServe(":"+port, newRouter)
	if err != nil {
		log.Fatalf("could not serve on port %s", port)
	}

	// Publish a task
	// /!\ Move to a controller triggered by http call /!\
	task := Task{"generate_conversions", 1}

	taskBytes, err := json.Marshal(task)
	if err != nil {
		fmt.Println(err)
	}

	err = taskQueue.PublishBytes(taskBytes)
	if err != nil {
		fmt.Println(err)
	}
}
