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
