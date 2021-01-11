package modeltests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/teeeeeeemo/go-crud-restapi/api/controllers"
	"github.com/teeeeeeemo/go-crud-restapi/api/models"
)

var server = controllers.Server{}
var userInstance = models.User{}
var postInstance = models.Post{}

/* Test Main */
func TestMain(m *testing.M) {

	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}

	Database()

	os.Exit(m.Run())

}

/* DB */
func Database() {

	var err error
	TestDbDriver := os.Getenv("TestDbDriver")

	/* mysql */
	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			os.Getenv("TestDbUser"),
			os.Getenv("TestDbPassword"),
			os.Getenv("TestDbHost"),
			os.Getenv("TestDbPort"),
			os.Getenv("TestDbName"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("this is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}

	/* postgres */
	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
			os.Getenv("TestDbHost"),
			os.Getenv("TestDbPort"),
			os.Getenv("TestDbUser"),
			os.Getenv("TestDbName"),
			os.Getenv("TestDbPassword"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("this is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}

}

/* 테이블 초기화 - user */
func refreshUserTable() error {

	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed table")
	return nil

}

/* 테이블 초기화 - post */
func refreshPostTable() error {

	err := server.DB.DropTableIfExists(&models.Post{}).Error
	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.Post{}).Error
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed table")
	return nil

}

/* dummy 데이터 user - 단건 */
func seedOneUser() (models.User, error) {

	refreshUserTable()

	user := models.User{
		Nickname: "Tang",
		Email:    "Tang@gmail.com",
		Password: "password",
	}

	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}

	return user, nil

}

/* dummy 데이터 user - 다건 */
func seedUsers() error {

	users := []models.User{
		models.User{
			Nickname: "Moondo",
			Email:    "moondo@lol.com",
			Password: "password",
		},
		models.User{
			Nickname: "Soraka",
			Email:    "soraka@lol.com",
			Password: "password",
		},
	}

	for i := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}

	return nil

}

/* 테이블 초기화 - user, post */
func refreshUserAndPostTable() error {

	err := server.DB.DropTableIfExists(&models.User{}, &models.Post{}).Error
	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.User{}, &models.Post{}).Error
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed tables")
	return nil

}

/* dummy 데이터 - user 단건, post 단건 */
func seedOneUserAndOnePost() (models.Post, error) {

	err := refreshUserAndPostTable()
	if err != nil {
		return models.Post{}, err
	}

	user := models.User{
		Nickname: "Ahri",
		Email:    "ahri@lol.com",
		Password: "password",
	}

	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.Post{}, err
	}

	post := models.Post{
		Title:    "This is the title - ahri",
		Content:  "This is the content - ahri",
		AuthorID: user.ID,
	}

	err = server.DB.Model(&models.Post{}).Create(&post).Error
	if err != nil {
		return models.Post{}, err
	}

	return post, nil

}

/* dummy 데이터 - user 다건, post 다건 */
func seedUsersAndPosts() ([]models.User, []models.Post, error) {

	var err error
	if err != nil {
		return []models.User{}, []models.Post{}, err
	}

	var users = []models.User{
		models.User{
			Nickname: "Malphite",
			Email:    "malphite@lol.com",
			Password: "password",
		},
		models.User{
			Nickname: "Amumu",
			Email:    "amumu@lol.com",
			Password: "password",
		},
	}

	var posts = []models.Post{
		models.Post{
			Title:   "Malphite Title",
			Content: "Malphite Content",
		},
		models.Post{
			Title:   "Amumu Title",
			Content: "Amumu Content",
		},
	}

	for i := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		posts[i].AuthorID = users[i].ID

		err = server.DB.Model(&models.Post{}).Create(&posts[i]).Error
	}

	return users, posts, nil

}
