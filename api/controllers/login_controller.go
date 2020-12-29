package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/teeeeeeemo/go-crud-restapi/api/auth"
	"github.com/teeeeeeemo/go-crud-restapi/api/models"
	"github.com/teeeeeeemo/go-crud-restapi/api/responses"
	"github.com/teeeeeeemo/go-crud-restapi/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

/* 로그인 */
func (server *Server) Login(w http.ResponseWriter, r *http.Request) {

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

	/* user 준비 */
	user.Prepare()
	/* user 유효성 검사 */
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	/* email, password로 signin */
	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, token)
}

/* email, password로 인증하여 토큰 생성 */
func (server *Server) SignIn(email, password string) (string, error) {

	var err error

	/* user 객체 할당 */
	user := models.User{}

	/* DB 조회: email에 해당하는 user */
	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}

	/* 비밀번호 검증 */
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	/* 토큰 생성하여 리턴 */
	return auth.CreateToken(user.ID)

}
