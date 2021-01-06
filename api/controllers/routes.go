package controllers

import (
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/teeeeeeemo/go-crud-restapi/api/middlewares"
	_ "github.com/teeeeeeemo/go-crud-restapi/docs"
)

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	// Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.DeleteUser))).Methods("DELETE")

	// Posts routes
	s.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreatePost))).Methods("POST")
	s.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetPosts))).Methods("GET")

	// Swagger
	s.Router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

}
