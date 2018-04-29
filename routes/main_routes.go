package routes

import (
	"fmt"
	"strconv"

	"github.com/peacecwz/go-social-app/models"

	CO "github.com/peacecwz/go-social-app/config"

	"github.com/kataras/iris"
)

// Index route
func Index(ctx iris.Context) {
	loggedIn(ctx, "/welcome")

	sessionId, _ := CO.AllSessions(ctx)
	if sessionId != "" {

		id, err := strconv.ParseInt(sessionId, 10, 64)
		CO.Err(err)
		db := CO.DB()
		posts := []interface{}{}
		rows, qErr := db.Raw("SELECT posts.\"id\", posts.\"title\", posts.\"content\", posts.\"created_by\", posts.\"created_at\" from \"public\".posts, \"public\".follows WHERE follows.\"follow_by\"=? AND follows.\"follow_to\" = posts.\"created_by\" ORDER BY posts.\"id\" DESC", id).Rows()
		CO.Err(qErr)
		var (
			postId     int
			title      string
			content    string
			created_by int
			created_at string
		)
		for rows.Next() {
			rows.Scan(&postId, &title, &content, &created_by, &created_at)
			post := map[string]interface{}{
				"postId":     postId,
				"title":      title,
				"content":    content,
				"created_by": created_by,
				"created_at": created_at,
			}
			posts = append(posts, post)
		}
		fmt.Println(posts)
		renderTemplate(ctx, "index", iris.Map{
			"title":   "Home",
			"session": ses(ctx),
			"posts":   posts,
			"GET":     CO.Get,
		})
	}

}

// Welcome route
func Welcome(ctx iris.Context) {
	notLoggedIn(ctx)
	renderTemplate(ctx, "welcome", iris.Map{
		"title": "Welcome",
	})
}

// NotFound route
func NotFound(ctx iris.Context) {
	renderTemplate(ctx, "404", iris.Map{
		"title":   "Oops!! Error",
		"session": ses(ctx),
	})
}

// Profile Page
func Profile(ctx iris.Context) {
	loggedIn(ctx, "")

	userId, err := strconv.Atoi(ctx.Params().Get("id"))
	if err != nil {
		CO.Err(err)
	}
	sesID, _ := CO.AllSessions(ctx)
	currentUserId, err := strconv.Atoi(sesID)
	if err != nil {
		CO.Err(err)
	}
	db := CO.DB()

	// VARS FOR USER DETAILS
	var (
		userCount int
		userID    int
		username  string
		email     string
		bio       string
	)

	// VARS FOR POSTS
	var (
		postID    int
		title     string
		content   string
		createdBy int
		createdAt string
	)
	posts := []interface{}{}

	var (
		followers  int //for followers
		followings int //for followings
		pViews     int // for profile views
	)

	me := CO.MeOrNot(ctx, currentUserId) // Check if its me or not
	var noMssg string                    // Mssg to be displayed when user has no posts

	if me == true {
		noMssg = "You have no posts. Go ahead and create one!!"
	} else {
		noMssg = username + " has no posts!!"

		// VIEW PROFILE
		if sesID != "" {
			profileViewsModel := db.Model(&(models.ProfileView{}))
			profileViewsModel.Create(&(models.ProfileView{
				ViewBy: currentUserId,
				ViewTo: userId}))
			db.Model(&(models.ProfileView{})).Create(&(models.ProfileView{
				ViewBy: currentUserId,
				ViewTo: userId}))

		}

	}

	// USER DETAILS
	db.Raw("SELECT COUNT(id) AS userCount, id AS userID, username, email, bio FROM users WHERE id=?", userId).Row().Scan(&userCount, &userID, &username, &email, &bio)

	invalid(ctx, userCount)

	// POSTS
	rows, err := db.Model(&(models.Post{})).Where("created_by = ?", userId).Rows()
	CO.Err(err)

	for rows.Next() {
		rows.Scan(&postID, &title, &content, &createdBy, &createdAt)
		post := map[string]interface{}{
			"postID":    postID,
			"title":     title,
			"content":   content,
			"createdBy": createdBy,
			"createdAt": createdAt,
		}
		posts = append(posts, post)
	}
	db.Model(&(models.Follow{})).Where("follow_to=?", userId).Select("id").Count(&followers)
	db.Model(&(models.Follow{})).Where("follow_by=?", userId).Select("id").Count(&followings)
	db.Model(&(models.ProfileView{})).Where("view_to=?", userId).Select("id").Count(&pViews)

	renderTemplate(ctx, "profile", iris.Map{
		"title":   "@" + username,
		"session": ses(ctx),
		"user": iris.Map{
			"id":       strconv.Itoa(userID),
			"username": username,
			"email":    email,
			"bio":      bio,
		},
		"posts":      posts,
		"followers":  followers,
		"followings": followings,
		"views":      pViews,
		"no_mssg":    noMssg,
		"GET":        CO.Get,
		"isF":        CO.IsFollowing,
	})

}

// Explore route
func Explore(ctx iris.Context) {
	loggedIn(ctx, "")
	user, _ := CO.AllSessions(ctx)
	db := CO.DB()
	explore := []interface{}{}
	usersModel := db.Model(&(models.User{}))
	rows, err := usersModel.Where("id <> ?", user).Select("id, username, email").Limit(10).Rows()
	CO.Err(err)

	for rows.Next() {
		var user models.User
		rows.Scan(&user)
		exp := map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		}
		explore = append(explore, exp)
	}

	renderTemplate(ctx, "explore", iris.Map{
		"title":   "Explore",
		"session": ses(ctx),
		"users":   explore,
		"GET":     CO.Get,
		"noF":     CO.NoOfFollowers,
		"UD":      CO.UsernameDecider,
	})
}

// CreatePost route
func CreatePost(ctx iris.Context) {
	loggedIn(ctx, "")
	renderTemplate(ctx, "create_post", iris.Map{
		"title":   "Create Post",
		"session": ses(ctx),
	})
}

// ViewPost route
func ViewPost(ctx iris.Context) {
	loggedIn(ctx, "")

	param := ctx.Params().Get("id")
	db := CO.DB()
	var post models.Post
	postsModel := db.Model(&(models.Post{}))
	likesModel := db.Model(&(models.Like{}))
	postsModel.Where("id = ?", param).First(&post)
	if post.ID == 0 {
		invalid(ctx, 0)
	}
	var likesCount int
	likesModel.Where("post_id = ?", param).Count(&likesCount)
	// likes
	//db.Table("likes").Where("postID=?", param).Count()
	renderTemplate(ctx, "view_post", iris.Map{
		"title":   "View Post",
		"session": ses(ctx),
		"post": iris.Map{
			"postID":    post.ID,
			"title":     post.Title,
			"content":   post.Content,
			"createdBy": post.CreatedBy,
			"createdAt": post.CreatedAt,
		},
		"postCreatedBy": strconv.Itoa(post.CreatedBy),
		"lon":           CO.LikedOrNot,
		"likes":         likesCount,
	})
}

// EditPost route
func EditPost(ctx iris.Context) {
	loggedIn(ctx, "")

	postID := ctx.Params().Get("id")
	db := CO.DB()
	postsModel := db.Model(&(models.Post{}))
	var post models.Post
	postsModel.Where("id = ?", postID).Select(&post)
	if post.ID == 0 {
		invalid(ctx, 0)
	}

	renderTemplate(ctx, "edit_post", iris.Map{
		"title":   "Edit Post",
		"session": ses(ctx),
		"post": iris.Map{
			"postID":  postID,
			"title":   post.Title,
			"content": post.Content,
		},
	})
}

// EditProfile route
func EditProfile(ctx iris.Context) {
	loggedIn(ctx, "")

	db := CO.DB()
	id, _ := CO.AllSessions(ctx)
	var (
		email  string
		bio    string
		joined string
	)
	db.Raw("SELECT email, bio, joined FROM users WHERE id=?", id).Row().Scan(&email, &bio, &joined)
	renderTemplate(ctx, "edit_profile", iris.Map{
		"title":   "Edit Profile",
		"session": ses(ctx),
		"email":   email,
		"bio":     bio,
		"joined":  joined,
	})
}

// Followers route
func Followers(ctx iris.Context) {
	loggedIn(ctx, "")

	user, err := strconv.Atoi(ctx.Params().Get("id"))
	if err != nil {
		CO.Err(err)
	}
	username := CO.Get(user, "username")
	db := CO.DB()
	var followBy int
	followers := []interface{}{}
	me := CO.MeOrNot(ctx, user)
	var noMssg string
	rows, err := db.Model(&(models.Follow{})).Where("follow_to=?", user).Order("id", true).Select("follow_by").Rows()
	CO.Err(err)
	for rows.Next() {
		rows.Scan(&followBy)
		f := map[string]interface{}{
			"followBy": followBy,
		}
		followers = append(followers, f)
	}

	if me == true {
		noMssg = "You"
	} else {
		noMssg = username
	}

	renderTemplate(ctx, "followers", iris.Map{
		"title":     username + "'s Followers",
		"session":   ses(ctx),
		"followers": followers,
		"no_mssg":   noMssg + " have no followers!!",
		"GET":       CO.Get,
		"UD":        CO.UsernameDecider,
		"noF":       CO.NoOfFollowers,
	})
}

// Followings route
func Followings(ctx iris.Context) {
	loggedIn(ctx, "")

	user, err := strconv.Atoi(ctx.Params().Get("id"))
	if err != nil {
		CO.Err(err)
	}
	username := CO.Get(user, "username")
	db := CO.DB()
	var followTo int
	followings := []interface{}{}
	me := CO.MeOrNot(ctx, user)
	var noMssg string

	result := db.Raw("SELECT followTo FROM follow WHERE followBy=? ORDER BY followID DESC", user)
	rows, fErr := result.Rows()
	CO.Err(fErr)

	for rows.Next() {
		rows.Scan(&followTo)
		f := map[string]interface{}{
			"followTo": followTo,
		}
		followings = append(followings, f)
	}

	if me == true {
		noMssg = "You"
	} else {
		noMssg = username
	}

	renderTemplate(ctx, "followings", iris.Map{
		"title":      username + "'s Followings",
		"session":    ses(ctx),
		"followings": followings,
		"no_mssg":    noMssg + " have no followings!!",
		"GET":        CO.Get,
		"UD":         CO.UsernameDecider,
		"noF":        CO.NoOfFollowers,
	})
}

// Likes route
func Likes(ctx iris.Context) {
	loggedIn(ctx, "")

	post := ctx.Params().Get("id")
	db := CO.DB()
	var postCount int
	var like_by int
	likes := []interface{}{}
	db.Model(&(models.Post{})).Where("post_id=?", post).Count(&postCount)
	invalid(ctx, postCount)

	rows, err := db.Model(&(models.Like{})).Where("post_id=?", post).Select("like_by").Rows()
	CO.Err(err)

	for rows.Next() {
		rows.Scan(&like_by)
		l := map[string]interface{}{
			"like_by": like_by,
		}
		likes = append(likes, l)
	}

	renderTemplate(ctx, "likes", iris.Map{
		"title":   "Likes",
		"session": ses(ctx),
		"likes":   likes,
		"GET":     CO.Get,
		"UD":      CO.UsernameDecider,
		"noF":     CO.NoOfFollowers,
	})
}

// Deactivate route
func Deactivate(ctx iris.Context) {
	renderTemplate(ctx, "deactivate", iris.Map{
		"title":   "Deactivate your acount",
		"session": ses(ctx),
	})
}
