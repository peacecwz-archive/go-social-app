package routes

import (
	"os"
	"strconv"

	"github.com/badoux/checkmail"
	"github.com/kataras/iris"
	CO "github.com/peacecwz/go-social-app/config"
	models "github.com/peacecwz/go-social-app/models"
)

// CreateNewPost route
func CreateNewPost(ctx iris.Context) {
	sessionId, _ := CO.AllSessions(ctx)
	id, err := strconv.Atoi(sessionId)
	if err != nil {
		panic(err)
	}
	post := models.Post{
		Title:     ctx.PostValueTrim("title"),
		Content:   ctx.PostValueTrim("content"),
		CreatedBy: id}

	db := CO.DB()
	db.Create(&post)

	resp := map[string]interface{}{
		"postID": post.ID,
		"mssg":   "Post Created!!",
	}
	json(ctx, resp)
}

// DeletePost route
func DeletePost(ctx iris.Context) {
	post := ctx.FormValue("post")
	db := CO.DB()
	db.Where("post_id=?", post).Delete(&(models.Post{}))

	json(ctx, map[string]interface{}{
		"mssg": "Post Deleted!!",
	})
}

// UpdatePost route
func UpdatePost(ctx iris.Context) {
	postID := ctx.PostValue("postID")
	title := ctx.PostValue("title")
	content := ctx.PostValue("content")

	db := CO.DB()
	db.Exec("UPDATE posts SET title=?, content=? WHERE postID=?", title, content, postID)

	json(ctx, map[string]interface{}{
		"mssg": "Post Updated!!",
	})
}

// UpdateProfile route
func UpdateProfile(ctx iris.Context) {
	resp := make(map[string]interface{})

	id, _ := CO.AllSessions(ctx)
	username := ctx.PostValueTrim("username")
	email := ctx.PostValueTrim("email")
	bio := ctx.PostValueTrim("bio")

	mailErr := checkmail.ValidateFormat(email)
	db := CO.DB()

	if username == "" || email == "" {
		resp["mssg"] = "Some values are missing!!"
	} else if mailErr != nil {
		resp["mssg"] = "Invalid email format!!"
	} else {
		result := db.Exec("UPDATE users SET username=?, email=?, bio=? WHERE id=?", username, email, bio, id)
		CO.Err(result.Error)

		session := CO.GetSession(ctx)
		session.Set("username", username)

		resp["mssg"] = "Profile updated!!"
		resp["success"] = true
	}

	json(ctx, resp)
}

// ChangeAvatar route
func ChangeAvatar(ctx iris.Context) {
	resp := make(map[string]interface{})
	id, _ := CO.AllSessions(ctx)

	dir, _ := os.Getwd()
	dest := dir + "/public/users/" + id + "/avatar.png"

	dErr := os.Remove(dest)
	CO.Err(dErr)

	// avatar key of post form file, but let's grab all of them,
	// the `ctx.FormFile` can be used to manually upload files per post key to the server.
	_, upErr := ctx.UploadFormFiles(dest)

	if upErr != nil {
		resp["mssg"] = "An error occured!!"
	} else {
		resp["mssg"] = "Avatar changed!!"
		resp["success"] = true
	}

	json(ctx, resp)
}

// Follow route
func Follow(ctx iris.Context) {
	sessionId, _ := CO.AllSessions(ctx)
	id, err := strconv.Atoi(sessionId)
	if err != nil {
		CO.Err(err)
	}
	user := ctx.PostValue("user")
	userId, err := strconv.Atoi(user)
	db := CO.DB()
	usersModel := db.Model(&(models.User{}))
	followsModel := db.Model(&(models.Follow{}))
	var currentUser, followerUser models.User
	usersModel.Where("id = ?", id).First(&currentUser)
	usersModel.Where("id = ?", userId).First(&followerUser)
	followsModel.Create(&(models.Follow{
		FollowBy: int(currentUser.ID),
		FollowTo: int(followerUser.ID)}))

	json(ctx, iris.Map{
		"mssg": "Followed " + currentUser.Username + "!!",
	})
}

// Unfollow route
func Unfollow(ctx iris.Context) {
	id, _ := CO.AllSessions(ctx)
	userId, err := strconv.Atoi(ctx.PostValue("user"))
	if err != nil {
		CO.Err(err)
	}
	username := CO.Get(userId, "username")

	db := CO.DB()
	db.Where("follow_by=? AND follow_to=?", id, userId).Delete(&(models.Follow{}))

	json(ctx, iris.Map{
		"mssg": "Unfollowed " + username + "!!",
	})
}

// Like post route
func Like(ctx iris.Context) {
	post := ctx.PostValue("post")
	postId, err := strconv.Atoi(post)
	if err != nil {
		CO.Err(err)
	}
	db := CO.DB()
	sessionId, _ := CO.AllSessions(ctx)
	id, err := strconv.Atoi(sessionId)
	if err != nil {
		CO.Err(err)
	}
	likesModel := db.Model(&(models.Like{}))
	likesModel.Create(&(models.Like{
		PostID: postId,
		LikeBy: id}))

	json(ctx, iris.Map{
		"mssg": "Post Liked!!",
	})
}

// Unlike post route
func Unlike(ctx iris.Context) {
	post := ctx.PostValue("post")
	id, _ := CO.AllSessions(ctx)
	db := CO.DB()
	db.Where("post_id=? AND like_by=?", post, id).Delete(&(models.Like{}))
	//CO.Err(result.Error)

	json(ctx, iris.Map{
		"mssg": "Post Unliked!!",
	})
}

// DeactivateAcc route post method
func DeactivateAcc(ctx iris.Context) {
	session := CO.GetSession(ctx)
	id, _ := CO.AllSessions(ctx)
	db := CO.DB()
	var postID int
	db.Where("view_by=?", id).Delete(&(models.ProfileView{}))
	db.Where("view_to=?", id).Delete(&(models.ProfileView{}))
	db.Where("follow_to=?", id).Delete(&(models.Follow{}))
	db.Where("follow_by=?", id).Delete(&(models.Follow{}))
	db.Where("like_by=?", id).Delete(&(models.Like{}))
	rows, err := db.Model(&(models.Post{})).Where("created_by=?", id).Select("id").Rows()
	CO.Err(err)
	for rows.Next() {
		rows.Scan(&postID)
		db.Where("post_id=?", postID).Delete(&(models.Like{}))
	}
	db.Where("created_by=?", id).Delete(&(models.Post{}))
	db.Where("id=?", id).Delete(&(models.User{}))

	dir, _ := os.Getwd()
	userPath := dir + "/public/users/" + id

	rmErr := os.RemoveAll(userPath)
	CO.Err(rmErr)

	// session.Delete("id")
	// session.Delete("username")
	// or
	session.Destroy()

	json(ctx, iris.Map{
		"mssg": "Deactivated your account!!",
	})
}
