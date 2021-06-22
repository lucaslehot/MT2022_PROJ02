package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/adjust/rmq/v3"
	"github.com/lucaslehot/MT2022_PROJ02/app/database"
	"github.com/lucaslehot/MT2022_PROJ02/app/models"
	"github.com/nfnt/resize"
)

func main() {
	err := database.Connect()
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}

	// Setting up Redis connection
	connection, err := rmq.OpenConnection("message_broker", "tcp", "redis-server:6379", 1, nil)
	taskQueue, err := connection.OpenQueue("tasks")

	fmt.Printf("queue connected: %v", taskQueue)

	// CREATE CONSUMER FUNCTION

	forever := make(chan bool)
	go func() {
		consumeErr := taskQueue.StartConsuming(10, time.Second) // donc la il récupères un truc avec task queue et en haut c'est le même task queue donc c'est bien fait quand même

		if consumeErr != nil {
			log.Fatalf("could not connect to db: %v", consumeErr)
		}

		taskQueue.AddConsumerFunc("log", func(delivery rmq.Delivery) {
			var task models.Task
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
			generateConversion(img)
			// 4 - store conversions in volume
			storeImage(img)

			log.Printf("Task performed", task)
			if err := delivery.Ack(); err != nil {
				// handle ack error
			}
		})
	}()
	log.Printf("Watcher running")
	<-forever
}

func getAvatar(userId int) image.Image {
	// 1 - get avatar url from db
	db := database.DbConn
	var user models.User
	db.Where("id = ?", userId).Find(&user)

	// 2 - retrieve image form volume
	// open avatar
	file, err1 := os.Open(user.AvatarPath)
	if err1 != nil {
		log.Fatal(err1)
	}

	// decode jpeg into image.Image
	img, err2 := jpeg.Decode(file)
	if err2 != nil {
		log.Fatal(err2)
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

func storeImage(img image.Image) {
	tempFile, err := ioutil.TempFile("./avatars", "upload-*-conversion.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// write this byte array to our temporary file
	err = jpeg.Encode(tempFile, img, nil)
	if err != nil {
		// Handle error
	}

	// return that we have successfully uploaded our file!
	fmt.Printf("Successfully Uploaded Conversion\n")
}
