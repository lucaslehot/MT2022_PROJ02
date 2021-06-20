package main

import (
	"encoding/json"
	"image"
	"image/jpeg"
	"log"
	"os"
	"time"

	"github.com/adjust/rmq/v3"
	"github.com/lucaslehot/MT2022_PROJ02/worker/database"
	"github.com/lucaslehot/MT2022_PROJ02/worker/models"
	"github.com/nfnt/resize"
)

type Task struct {
	Order  string
	UserId int
}

func main() {
	err := database.Connect()
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}

	// Setting up Redis connection
	connection, err := rmq.OpenConnection("message_broker", "tcp", "localhost:6379", 1, nil)
	taskQueue, err := connection.OpenQueue("tasks")

	// CREATE CONSUMER FUNCTION
	err := taskQueue.StartConsuming(10, time.Second)
	name, err := taskQueue.AddConsumerFunc(func(delivery rmq.Delivery) {
		var task Task
		if err = json.Unmarshal([]byte(delivery.Payload()), &task); err != nil {
			// handle json error
			if err := delivery.Reject(); err != nil {
				// handle reject error
			}
			return
		}

		// perform task
		img := getAvatar(task.UserId)

		// 3 - generate image conversions
		// 4 - store conversions in volume

		log.Printf("performing task %s", task)
		if err := delivery.Ack(); err != nil {
			// handle ack error
		}

	})
}

func getAvatar(userId int) {
	// 1 - get avatar url from db
	db := database.DbConn
	var user models.User
	err := db.Where("id = ?", userId).Find(&user)
	if err != nil {
		log.Fatal(err)
	}

	// 2 - retrieve image form volume
	// open avatar
	file, err := os.Open(user.avatar_path)
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	return img
}

func generateConversion(img image.Image) {
	// resize to width 1000 using NearestNeighbor resampling
	// and preserve aspect ratio
	m := resize.Resize(1000, 0, img, resize.NearestNeighbor)

	out, err := os.Create("test_resized.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
}
