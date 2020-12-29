package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/teeeeeeemo/go-crud-restapi/api/models"
	"github.com/teeeeeeemo/go-crud-restapi/api/responses"
	"github.com/teeeeeeemo/go-crud-restapi/api/utils/formaterror"
)

/* user 생성 */
func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {

	/* 요청의 바디를 모두 읽기 */
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	/* user 객체 할당 */
	user := models.User{}
	/* body의 데이터를 user로 언마샬링 */
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	/* user 준비 */
	user.Prepare()
	/* user 유효성 검사 */
	err = user.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	/* user 저장 */
	userCreated, err := user.SaveUser(server.DB)
	if err != nil {
		formatedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formatedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	responses.JSON(w, http.StatusCreated, userCreated)

}

/* user 목록 조회 */
func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {

	/* user 객체 할당 */
	user := models.User{}

	/* user 목록 조회 */
	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, users)
}

/* user 상세 조회 */
func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {

	/* Vars returns the route variables for the current request, if any.
	Vars는 현재 요청에 대한 경로 변수 반환함 */
	vars := mux.Vars(r)
	/* string -> uint 변환 */
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	/* user 객체 할당 */
	user := models.User{}
	/* id로 user 조회 */
	userGotten, err := user.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, userGotten)
}
