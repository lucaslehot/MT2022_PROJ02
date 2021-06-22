package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/gobuffalo/packr/v2"

	"github.com/adjust/rmq/v3"
	"github.com/lucaslehot/MT2022_PROJ02/app/database"
	"github.com/lucaslehot/MT2022_PROJ02/app/models"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	// Setting up Redis connection
	connection, err := rmq.OpenConnection("message_broker", "tcp", "redis-server:6379", 1, nil)
	taskQueue, err := connection.OpenQueue("tasks")

	redisErr := database.Connect()
	if redisErr != nil {
		log.Fatalf("Impossible de se connecter à la bdd: %v", redisErr)
	}

	//-----------------------------
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("/avatars", "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")

	db := database.DbConn

	db.Model(&models.User{}).Where("id = ?", 1).Update("avatar_path", "./avatars" + string(1))

	task := models.Task{"generate_conversions", 1}

	taskBytes, err := json.Marshal(task)
	if err != nil {
		fmt.Println(err)
	}

	err = taskQueue.PublishBytes(taskBytes)
	if err != nil {
		fmt.Println(err)
	}
}

var box = packr.New("templateBox", "../views")

func RenderImageForm(w http.ResponseWriter, r *http.Request) {
	tpl := template.New("uploadImageForm.html")
	articleForm, err := box.FindString("uploadImageForm.html")
	if err != nil {
		log.Print(err)
		return
	}
	t, err := tpl.Parse(articleForm) // Parse template file.
	if err != nil {
		log.Print(err)
		return
	}
	err = t.Execute(w, nil) // merge.
	if err != nil {
		log.Print(err)
		return
	}
}
