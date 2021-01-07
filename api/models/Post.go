package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

/* post 구조체 */
type Post struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id,string"` // 아이디
	Title     string    `gorm:"size:255;not null;" json:"title"`             // 제목
	Content   string    `gorm:"size:255;not null;" json:"content"`           // 내용
	Author    User      `json:"author"`                                      // 작성자
	AuthorID  uint32    `gorm:"not null" json:"author_id,string"`            // 작성자 아이디
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"` // 생성시간
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"` // 수정시간
}

/* post 준비 메서드 */
func (p *Post) Prepare() {

	p.ID = 0
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

}

/* post 유효성 검사 메서드 */
func (p *Post) Validate() error {

	if p.Title == "" {
		return errors.New("Required Title")
	}
	if p.Content == "" {
		return errors.New("Required Content")
	}
	if p.AuthorID < 1 {
		return errors.New("Required Author")
	}
	return nil

}

/* post 저장 메서드 */
func (p *Post) SavePost(db *gorm.DB) (*Post, error) {

	var err error
	/* DB 삽입: post */
	err = db.Debug().Model(&Post{}).Create(&p).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		/* DB 조회: user by id */
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}

	return p, nil

}

// TODO pagination
/* post 목록 조회 메서드 */
func (p *Post) FindAllPosts(db *gorm.DB) (*[]Post, error) {

	var err error
	posts := []Post{}
	/* DB 조회: post 목록, 100개 */
	err = db.Debug().Model(&Post{}).Limit(100).Find(&posts).Error
	if err != nil {
		return &[]Post{}, err
	}
	if len(posts) > 0 {
		for i := range posts {
			/* DB 조회: user by id */
			err := db.Debug().Model(&User{}).Where("id = ?", posts[i].AuthorID).Take(&posts[i].Author).Error
			if err != nil {
				return &[]Post{}, err
			}
		}
	}
	return &posts, nil

}

/* post 상세 조회 메서드 */
func (p *Post) FindPostByID(db *gorm.DB, pid uint64) (*Post, error) {

	var err error
	/* DB 조회: post by id */
	err = db.Debug().Model(&Post{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		/* DB 조회: user by id */
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil

}

/* post 수정 메서드 */
func (p *Post) UpdateAPost(db *gorm.DB) (*Post, error) {

	var err error

	/* DB 수정: post by id */
	err = db.Debug().Model(&Post{}).Where("id = ?", p.ID).Updates(Post{Title: p.Title, Content: p.Content, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		/* DB 조회: user by id */
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil

}

/* post 삭제 메서드 */
func (p *Post) DeleteAPost(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	/* DB 삭제: post by id and author_id(user id) */
	db = db.Debug().Model(&Post{}).Where("id = ? and author_id = ?", pid, uid).Take(&Post{}).Delete(&Post{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Post not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil

}
