package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/teeeeeeemo/go-crud-restapi/api/auth"
	"github.com/teeeeeeemo/go-crud-restapi/api/models"
	"github.com/teeeeeeemo/go-crud-restapi/api/responses"
	"github.com/teeeeeeemo/go-crud-restapi/api/utils/formaterror"
)

/* user 생성 메서드 */
// @Summary Create a User
// @Description Create a user with the input payload
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "Create user"
// @Success 200 {object} models.User
// @Router /users [post]
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

/* user 목록 조회 메서드 */
// @Summary Get User List
// @Description 사용자 목록 조회
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
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

/* user 상세 조회 메서드 */
// @Summary Show User Details
// @Description 사용자 상세 조회
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "user id"
// @Success 200 {object} models.User
// @Router /users/{id} [get]
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

/* user 수정 메서드 */
// @Summary Update a User
// @Description Update a user with the input payload
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "user id"
// @Param user body models.User true "Update user"
// @Success 200 {object} models.User
// @Router /users/{id} [put]
func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {

	/* Vars returns the route variables for the current requests, if any.
	Vars는 현재 요청에 대한 경로 변수 반환함 */
	vars := mux.Vars(r)
	/* string -> uint 변환 */
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	/* 요청의 바디를 모두 읽기 */
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	/* user 객체 할당 */
	user := models.User{}
	/* body의 데이터를 user로 언마샬링 */
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	/* 토큰 아이디 추출 */
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	/* 토큰 아이디 검증 */
	if tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	/* user 준비 */
	user.Prepare()
	/* user 유효성 검사 */
	err = user.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	/* user 수정 */
	updatedUser, err := user.UpdateAUser(server.DB, uint32(uid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, updatedUser)

}

/* user 삭제 메서드 */
// @Summary Delete a User
// @Description 사용자 삭제
// @Tags users
// @Accept json
// @Produce json
// @Param id query string true "user id"
// @Success 204
// @Router /users/{id} [delete]
func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {

	/* Vars returns the route variables for the current request, if any.
	Vars는 현재 요청에 대한 경로 변수 반환함 */
	vars := mux.Vars(r)

	/* user 객체 할당 */
	user := models.User{}

	/* string -> uint 변환 */
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	/* 토큰 아이디 추출 */
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	/* 토큰 아이디 검증 */
	if tokenID != 0 && tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	/* user 삭제 */
	_, err = user.DeleteAUser(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.JSON(w, http.StatusNoContent, "")

}
