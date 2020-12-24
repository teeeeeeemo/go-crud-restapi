package controllers

import (
	"net/http"

	"github.com/teeeeeeemo/go-crud-restapi/api/responses"
)

/* Home 응답 */
func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To Hyun-Go-Rest-API")
}
