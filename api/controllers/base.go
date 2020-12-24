package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql database driver
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres database driver

	"github.com/teeeeeeemo/go-crud-restapi/api/models"
)

/* server 구조체 */
type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

/* 초기화
### db 정보:
	- 드라이버
	- 계정이름
	- 비밀번호
	- 포트
	- 호스트
	- 이름
*/
func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error

	/* 드라이버: mysql */
	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", Dbdriver)
		}
	}

	/* 드라이버: postgres */
	if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=$s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	}

	/* DB 마이그레이션 */
	server.DB.Debug().AutoMigrate(&models.User{}, &models.Post{})

	/* Router 초기화 */
	server.Router = mux.NewRouter()
	server.initializeRoutes()
}

/* server 구동 */
func (server *Server) Run(addr string) {
	// fmt.Println("Listening to port 7878")
	fmt.Printf("Listening to port %s\n", addr)
	/* listen 할 주소 정보와 handler를 인자로 서버를 요청 대기 상태로 만듬 */
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
