package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/teeeeeeemo/go-crud-restapi/api/models"
)

/* dummy 데이터용 user 슬라이스 선언 및 초기화 */ // slice: slice는 배열과 같지만, 길이가 고정되어 있지 않으며 동적으로 크기가 늘어난다.
var users = []models.User{
	models.User{
		Nickname: "Vayne",
		Email:    "vayne@lol.com",
		Password: "password",
	},
	models.User{
		Nickname: "Teemo",
		Email:    "teemo@lol.com",
		Password: "password",
	},
}

/* dummy 데이터용 post 슬라이스 선언 및 초기화 */
var posts = []models.Post{
	models.Post{
		Title:   "Title 1",
		Content: "Content 1",
	},
	models.Post{
		Title:   "Title",
		Content: "Content 2",
	},
}

/* dummy 데이터 로딩 */
func Load(db *gorm.DB) {

	/* table 존재하면 drop */
	err := db.Debug().DropTableIfExists(&models.Post{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	} else {
		log.Println("drop table 완료")
	}

	/* table 자동 마이그레이션 */
	err = db.Debug().AutoMigrate(&models.User{}, &models.Post{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	} else {
		log.Println("migrate table 완료")
	}

	/* table 외래키 추가 */
	err = db.Debug().Model(&models.Post{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	} else {
		log.Println("add foreign key 완료")
	}

	/* dummy 데이터 추가: users, posts */
	// for i, _ := range users {
	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		} else {
			log.Println("seed users table 완료")
		}
		posts[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		} else {
			log.Println("seed posts table 완료")
		}
	}
}
