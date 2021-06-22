package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lucaslehot/MT2022_PROJ02/app/controllers"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
	Route{
		Name:        "Upload file",
		Method:      "POST",
		Pattern:     `/image/upload`,
		HandlerFunc: controllers.UploadFile,
	},

	Route{
		Name:        "Display image form",
		Method:      "GET",
		Pattern:     "/image/new",
		HandlerFunc: controllers.RenderImageForm,
	},

	// User management routes

	Route{
		Name:        "Create user",
		Method:      "POST",
		Pattern:     "/user/create",
		HandlerFunc: controllers.CreateUser,
	},

	Route{
		Name:        "Read user",
		Method:      "GET",
		Pattern:     `/user/{ID}`,
		HandlerFunc: controllers.ReadUser,
	},

	Route{
		Name:        "Update user",
		Method:      "POST",
		Pattern:     "/user/update",
		HandlerFunc: controllers.UpdateUser,
	},

	Route{
		Name:        "Delete user",
		Method:      "POST",
		Pattern:     `/delete/delete`,
		HandlerFunc: controllers.DeleteUser,
	},
}
