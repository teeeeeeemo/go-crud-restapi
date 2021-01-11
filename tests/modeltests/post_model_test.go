package modeltests

import (
	"log"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/teeeeeeemo/go-crud-restapi/api/models"
	"gopkg.in/go-playground/assert.v1"
)

/* Test: post 목록 조회 */
func TestFindAllPosts(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table %v\n", err)
	}

	_, _, err = seedUsersAndPosts()
	if err != nil {
		log.Fatalf("Error seeding user and post table %v\n", err)
	}

	posts, err := postInstance.FindAllPosts(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the posts: %v\n", err)
		return
	}

	assert.Equal(t, len(*posts), 2)

}

/* Test: post 저장 */
func TestSavePost(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error user and post refreshing table %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	newPost := models.Post{
		ID:       1,
		Title:    "제목제목제목",
		Content:  "내용내용내용",
		AuthorID: user.ID,
	}

	savedPost, err := newPost.SavePost(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the post: %v\n", err)
		return
	}

	assert.Equal(t, newPost.ID, savedPost.ID)
	assert.Equal(t, newPost.Title, savedPost.Title)
	assert.Equal(t, newPost.Content, savedPost.Content)
	assert.Equal(t, newPost.AuthorID, savedPost.AuthorID)

}

/* Test: post 상세 조회 by id */
func TestGetPostByID(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}

	post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}

	foundPost, err := postInstance.FindPostByID(server.DB, post.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}

	assert.Equal(t, foundPost.ID, post.ID)
	assert.Equal(t, foundPost.Title, post.Title)
	assert.Equal(t, foundPost.Content, post.Content)

}

/* Test: post 수정 */
func TestUpdateAPost(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table %v\n", err)
	}

	post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}

	postUpdate := models.Post{
		ID:       1,
		Title:    "수정된제목제목제목",
		Content:  "수정된내용내용내용",
		AuthorID: post.AuthorID,
	}

	updatedPost, err := postUpdate.UpdateAPost(server.DB)
	if err != nil {
		t.Errorf("this is the error updating the post: %v\n", err)
		return
	}

	assert.Equal(t, updatedPost.ID, postUpdate.ID)
	assert.Equal(t, updatedPost.Title, postUpdate.Title)
	assert.Equal(t, updatedPost.Content, postUpdate.Content)
	assert.Equal(t, updatedPost.AuthorID, postUpdate.AuthorID)

}

/* Test: post 삭제 */
func TestDeleteAPost(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}

	post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}

	isDeleted, err := postInstance.DeleteAPost(server.DB, post.ID, post.AuthorID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}

	// one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	// Can be done this way too
	assert.Equal(t, isDeleted, int64(1))

}
