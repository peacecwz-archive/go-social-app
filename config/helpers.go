package config

import (
	"fmt"
	"strconv"

	"github.com/peacecwz/go-social-app/models"
)

// THIS FILE CONTAINS ALL THE METHODS WHICH WILL BE USED IN TEMPLATES/VIEWS

// Get function to get anything of user with ID
func Get(id int, what string) string {
	db := DB()
	usersModel := db.Model(&(models.User{}))
	var RET string
	fmt.Println(what)
	usersModel.Where("id = ?", id).Select(what).Row().Scan(&RET)
	return RET
}

// IsFollowing route
func IsFollowing(by string, to string) bool {
	db := DB()
	followsModel := db.Model(&(models.Follow{}))
	var followCount int
	followsModel.Where("follow_by=? AND follow_to=?", by, to).Count(&followCount)
	if followCount == 0 {
		return false
	}
	return true
}

// UsernameDecider Helper
func UsernameDecider(user int, session string) string {
	username := Get(user, "username")
	sessionUserId, err := strconv.Atoi(session)
	if err != nil {
		Err(err)
	}
	sesUsername := Get(sessionUserId, "username")
	if username == sesUsername {
		return "You"
	}
	return username
}

// NoOfFollowers helper
func NoOfFollowers(user int) int {
	db := DB()
	followsModel := db.Model(&(models.Follow{}))
	var followersCount int
	followsModel.Where("follow_to=?", user).Count(&followersCount)
	//db.QueryRow("SELECT COUNT(followID) AS followersCount FROM  WHERE ", user).Scan()
	return followersCount
}

// LikedOrNot helper
func LikedOrNot(post int, user interface{}) bool {
	db := DB()
	likesModel := db.Model(&(models.Like{}))
	var likeCount int
	likesModel.Where("like_by=? AND post_id=?", user, post).Select("id").Count(&likeCount)
	//db.QueryRow("SELECT COUNT(likeID) AS likeCount FROM likes WHERE ", user, post).Scan()
	if likeCount == 0 {
		return false
	}
	return true
}
