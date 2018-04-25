package routes

import (
	"os"
	"time"

	CO "github.com/peacecwz/go-social-app/config"

	"github.com/badoux/checkmail"
	"github.com/kataras/iris"
)

// CreateNewPost route
func CreateNewPost(ctx iris.Context) {

	title := ctx.PostValueTrim("title")
	content := ctx.PostValueTrim("content")
	id, _ := CO.AllSessions(ctx)

	db := CO.DB()

	rs := db.Exec("INSERT INTO posts(title, content, createdBy, createdAt) VALUES (?, ?, ?, ?)", title, content, id, time.Now())
	CO.Err(rs.Error)
	var insertID int

	rs.Select("postID").Last(&insertID)

	resp := map[string]interface{}{
		"postID": insertID,
		"mssg":   "Post Created!!",
	}
	json(ctx, resp)
}

// DeletePost route
func DeletePost(ctx iris.Context) {
	post := ctx.FormValue("post")
	db := CO.DB()

	rs := db.Exec("DELETE FROM posts WHERE postID=?", post)
	CO.Err(rs.Error)

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
	id, _ := CO.AllSessions(ctx)
	user := ctx.PostValue("user")
	username := CO.Get(user, "username")

	db := CO.DB()
	result := db.Exec("INSERT INTO follow(followBy, followTo, followTime) VALUES(?, ?, ?)", id, user, time.Now())

	CO.Err(result.Error)

	json(ctx, iris.Map{
		"mssg": "Followed " + username + "!!",
	})
}

// Unfollow route
func Unfollow(ctx iris.Context) {
	id, _ := CO.AllSessions(ctx)
	user := ctx.PostValue("user")
	username := CO.Get(user, "username")

	db := CO.DB()
	result := db.Exec("DELETE FROM follow WHERE followBy=? AND followTo=?", id, user)
	CO.Err(result.Error)

	json(ctx, iris.Map{
		"mssg": "Unfollowed " + username + "!!",
	})
}

// Like post route
func Like(ctx iris.Context) {
	post := ctx.PostValue("post")
	db := CO.DB()
	id, _ := CO.AllSessions(ctx)

	result := db.Exec("INSERT INTO likes(postID, likeBy, likeTime) VALUES (?, ?, ?)", post, id, time.Now())
	CO.Err(result.Error)

	json(ctx, iris.Map{
		"mssg": "Post Liked!!",
	})
}

// Unlike post route
func Unlike(ctx iris.Context) {
	post := ctx.PostValue("post")
	id, _ := CO.AllSessions(ctx)
	db := CO.DB()

	result := db.Exec("DELETE FROM likes WHERE postID=? AND likeBy=?", post, id)
	CO.Err(result.Error)

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

	db.Exec("DELETE FROM profile_views WHERE viewBy=?", id)
	db.Exec("DELETE FROM profile_views WHERE viewTo=?", id)
	db.Exec("DELETE FROM follow WHERE followBy=?", id)
	db.Exec("DELETE FROM follow WHERE followTo=?", id)
	db.Exec("DELETE FROM likes WHERE likeBy=?", id)
	rows, err := db.Table("posts").Where("createdBy=?", id).Select("postID").Rows()
	CO.Err(err)
	for rows.Next() {
		rows.Scan(&postID)
		db.Exec("DELETE FROM likes WHERE postID=?", postID)
	}

	db.Exec("DELETE FROM posts WHERE createdBy=?", id)
	db.Exec("DELETE FROM users WHERE id=?", id)

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
