package config

// THIS FILE CONTAINS ALL THE METHODS WHICH WILL BE USED IN TEMPLATES/VIEWS

// Get function to get anything of user with ID
func Get(id interface{}, what string) string {
	db := DB()
	var RET string
	db.Table("users").Where("id = ?", id).Select(what).Row().Scan(&RET)
	return RET
}

// IsFollowing route
func IsFollowing(by string, to string) bool {
	db := DB()
	var followCount int
	db.Table("follow").Where("followBy=? AND followTo=?", by, to).Select("followID").Count(&followCount)
	//db.QueryRow("SELECT COUNT(followID) AS followCount FROM  WHERE  LIMIT 1", by, to).Scan(&followCount)
	if followCount == 0 {
		return false
	}
	return true
}

// UsernameDecider Helper
func UsernameDecider(user int, session string) string {
	username := Get(user, "username")
	sesUsername := Get(session, "username")
	if username == sesUsername {
		return "You"
	}
	return username
}

// NoOfFollowers helper
func NoOfFollowers(user int) int {
	db := DB()
	var followersCount int
	db.Table("follow").Where("followTo=?", user).Select("followID").Count(&followersCount)
	//db.QueryRow("SELECT COUNT(followID) AS followersCount FROM  WHERE ", user).Scan()
	return followersCount
}

// LikedOrNot helper
func LikedOrNot(post int, user interface{}) bool {
	db := DB()
	var likeCount int
	db.Table("likes").Where("likeBy=? AND postID=?", user, post).Select("likeID").Count(&likeCount)
	//db.QueryRow("SELECT COUNT(likeID) AS likeCount FROM likes WHERE ", user, post).Scan()
	if likeCount == 0 {
		return false
	}
	return true
}
