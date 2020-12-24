package api

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/teeeeeeemo/go-crud-restapi/api/controllers"
	"github.com/teeeeeeemo/go-crud-restapi/api/seed"
)

var server = controllers.Server{}

/* main.go에서 호출
   - 어플리케이션 초기화
   - DB에 dummy 데이터 추가 */
func Run() {

	var err error
	/* .env 파일에서 환경 변수를 로딩함 */
	err = godotenv.Load()

	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	/* server 초기화:
	### db 정보:
		- 드라이버
		- 계정이름
		- 비밀번호
		- 포트
		- 호스트
		- 이름 */
	server.Initialize(
		os.Getenv("DB_DRIVER"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	/* dummy 데이터 로딩 */
	seed.Load(server.DB)

	/* server 실행: 포트 번호 지정 */
	server.Run(":7878")
}
