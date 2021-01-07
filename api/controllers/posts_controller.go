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

/* post 생성 메서드 */
// @Summary Create a Post
// @Description Create a post with the input payload
// @Tags posts
// @Accept json
// @Produce json
// @Param post body models.Post true "Create post"
// @Sucess 201 {object} models.Post
// @Router /posts [post]
func (server *Server) CreatePost(w http.ResponseWriter, r *http.Request) {

	/* 요청의 바디를 모두 읽기 */
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	/* post 객체 할당 */
	post := models.Post{}
	/* body의 데이터를 post로 언마샬링 */
	err = json.Unmarshal(body, &post)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	/* post 준비 */
	post.Prepare()
	/* post 유효성 검사 */
	err = post.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	/* 토큰 아이디 추출 */
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	/* 토큰 아이디 검증 */
	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	/* post 저장 */
	postCreated, err := post.SavePost(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, postCreated.ID))
	responses.JSON(w, http.StatusCreated, postCreated)
}

/* post 목록 조회 메서드 */
// @Summary Get Post List
// @Description 포스트 목록 조회
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {array} models.Post
// @Router /posts [get]
func (server *Server) GetPosts(w http.ResponseWriter, r *http.Request) {

	/* post 객체 할당 */
	post := models.Post{}

	/* post 목록 조회 */
	posts, err := post.FindAllPosts(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, posts)

}

/* post 상세 조회 메서드 */
// @Summary Show Post Details
// @Description 포스트 상세 조회
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "post id"
// @Success 200 {object} models.Post
// @Router /posts/{id} [get]
func (server *Server) GetPost(w http.ResponseWriter, r *http.Request) {

	/* Vars returns the route variables for the current request, if any.
	Vars는 현재 요청에 대한 경로 변수 반환함 */
	vars := mux.Vars(r)
	/* string -> uint 변환 */
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	/* post 객체 할당 */
	post := models.Post{}

	/* id로 post 조회 */
	postReceived, err := post.FindPostByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, postReceived)

}

/* post 수정 메서드 */
// @Summary Update a Post
// @Description Update a post with the input payload
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "post id"
// @Param post body models.Post true "Update post"
// @Success 200 {object} models.Post
// @Router /posts/{id} [put]
func (server *Server) UpdatePost(w http.ResponseWriter, r *http.Request) {

	/* Vars returns the route variables for the current request, if any.
	Vars는 현재 요청에 대한 경로 변수 반환함 */
	vars := mux.Vars(r)

	/* post id 유효성 검사 */
	/* string -> uint 변환 */
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	/* 토큰 아이디 추출 */
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	/* post 존재 여부 확인 */
	/* post 객체 할당 */
	post := models.Post{}
	/* id로 post 조회 */
	postReceived, err := post.FindPostByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	/* user id와 post 작성자 id가 다를 경우 */
	if uid != postReceived.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	/* 요청의 바디를 모두 읽기 */
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	/* 요청 데이터로 post 업데이트 */
	/* post 객체 할당 */
	postUpdate := models.Post{}
	/* body의 데이터를 post로 언마샬링 */
	err = json.Unmarshal(body, &postUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	/* post 준비 */
	postUpdate.Prepare()
	/* post 유효성 검사 */
	err = postUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	/* 요청 데이터의 user id와 토큰에서 추출한 user id 동일 여부 체크 */
	if uid != postUpdate.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	/* 수정할 post의 id 세팅 */
	postUpdate.ID = post.ID // pid로 해도 같잖아

	/* post 수정 */
	postUpdated, err := postUpdate.UpdateAPost(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, postUpdated)

}
