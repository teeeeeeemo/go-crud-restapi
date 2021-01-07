package controllers

import (
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/teeeeeeemo/go-crud-restapi/api/middlewares"
	_ "github.com/teeeeeeemo/go-crud-restapi/docs"
)

func (s *Server) initializeRoutes() {

	/* ---------- Home Route ---------- */
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	/* ---------- Login Route ---------- */
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	/* ---------- Users routes ---------- */
	/* Create User */
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	/* Get User List */
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	/* Show User Details */
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	/* Update User */
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	/* Delete User */
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.DeleteUser))).Methods("DELETE")

	/* ---------- Posts routes ---------- */
	/* Create Post */
	s.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreatePost))).Methods("POST")
	/* Get Post List */
	s.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetPosts))).Methods("GET")
	/* Show Post Details */
	s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetPost))).Methods("GET")
	/* Update Post */
	s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdatePost))).Methods("PUT")

	/* ---------- Swagger ---------- */
	s.Router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

}
