package controllers

import "github.com/teeeeeeemo/go-crud-restapi/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route

	// Users routes

}
