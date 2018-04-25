package routes

import (
	"strconv"
	"time"

	CO "github.com/peacecwz/go-social-app/config"

	"github.com/kataras/iris"
)

// Index route
func Index(ctx iris.Context) {
	loggedIn(ctx, "/welcome")

	id, _ := CO.AllSessions(ctx)
	db := CO.DB()
	var (
		postID    int
		title     string
		content   string
		createdBy int
		createdAt string
	)
	feeds := []interface{}{}
	rows, qErr := db.Raw("SELECT posts.postID, posts.title, posts.content, posts.createdBy, posts.createdAt from posts, follow WHERE follow.followBy=? AND follow.followTo = posts.createdBy ORDER BY posts.postID DESC", id).Rows()
	CO.Err(qErr)

	for rows.Next() {
		rows.Scan(&postID, &title, &content, &createdBy, &createdAt)
		feed := map[string]interface{}{
			"postID":    postID,
			"title":     title,
			"content":   content,
			"createdBy": createdBy,
			"createdAt": createdAt,
		}
		feeds = append(feeds, feed)
	}

	renderTemplate(ctx, "index", iris.Map{
		"title":   "Home",
		"session": ses(ctx),
		"posts":   feeds,
		"GET":     CO.Get,
	})
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

	user := ctx.Params().Get("id")
	sesID, _ := CO.AllSessions(ctx)
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

	me := CO.MeOrNot(ctx, user) // Check if its me or not
	var noMssg string           // Mssg to be displayed when user has no posts

	if me == true {
		noMssg = "You have no posts. Go ahead and create one!!"
	} else {
		noMssg = username + " has no posts!!"

		// VIEW PROFILE
		if sesID != "" {
			result := db.Exec("INSERT INTO profile_views(viewBy, viewTo, viewTime) VALUES(?, ?, ?)", sesID, user, time.Now())

			CO.Err(result.Error)
		}

	}

	// USER DETAILS
	db.Raw("SELECT COUNT(id) AS userCount, id AS userID, username, email, bio FROM users WHERE id=?", user).Row().Scan(&userCount, &userID, &username, &email, &bio)

	invalid(ctx, userCount)

	// POSTS
	result := db.Raw("SELECT * FROM posts WHERE createdBy=? ORDER BY postID DESC")
	CO.Err(result.Error)
	rows, gErr := result.Rows()
	CO.Err(gErr)

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
	db.Table("follow").Where("followTo=?", user).Select("followID").Count(&followers)
	db.Table("follow").Where("followBy=?", user).Select("followID").Count(&followings)
	db.Table("profile_views").Where("viewTo=?", user).Select("viewID").Count(&pViews)
	//db.Raw("SELECT COUNT() AS followers FROM follow WHERE ").Scan()                                         // FOLLOWERS
	//db.QueryRow("SELECT COUNT(followID) AS followers FROM follow WHERE followBy=?", user).Scan() // FOLLOWINGS
	//db.QueryRow("SELECT COUNT() AS pViews FROM  WHERE ", user).Scan() // PROFILE VIEWS

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
	var (
		id       int
		username string
		email    string
	)
	explore := []interface{}{}

	result := db.Raw("SELECT id, username, email FROM users WHERE id <> ? ORDER BY RAND() LIMIT 10", user)
	rows, err := result.Rows()
	CO.Err(err)

	for rows.Next() {
		rows.Scan(&id, &username, &email)
		exp := map[string]interface{}{
			"id":       id,
			"username": username,
			"email":    email,
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
	var (
		postCount int
		postID    int
		title     string
		content   string
		createdBy int
		createdAt string
	)
	var likesCount int

	// post details
	db.Table("posts").Where("postID=?", param).Select("COUNT(postID) AS postCount, postID, title, content, createdBy, createdAt").Row().Scan(&postCount, &postID, &title, &content, &createdBy, &createdAt)
	//db.QueryRow("SELECT  FROM  WHERE ")
	invalid(ctx, postCount)

	// likes
	db.Table("likes").Where("postID=?", param).Count(&likesCount)

	renderTemplate(ctx, "view_post", iris.Map{
		"title":   "View Post",
		"session": ses(ctx),
		"post": iris.Map{
			"postID":    postID,
			"title":     title,
			"content":   content,
			"createdBy": createdBy,
			"createdAt": createdAt,
		},
		"postCreatedBy": strconv.Itoa(createdBy),
		"lon":           CO.LikedOrNot,
		"likes":         likesCount,
	})
}

// EditPost route
func EditPost(ctx iris.Context) {
	loggedIn(ctx, "")

	post := ctx.Params().Get("id")
	db := CO.DB()
	var (
		postCount int
		postID    int
		title     string
		content   string
	)
	db.Raw("SELECT COUNT(postID) AS postCount, postID, title, content FROM posts WHERE postID=?", post).Row().Scan(&postCount, &postID, &title, &content)
	invalid(ctx, postCount)

	renderTemplate(ctx, "edit_post", iris.Map{
		"title":   "Edit Post",
		"session": ses(ctx),
		"post": iris.Map{
			"postID":  postID,
			"title":   title,
			"content": content,
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

	user := ctx.Params().Get("id")
	username := CO.Get(user, "username")
	db := CO.DB()
	var followBy int
	followers := []interface{}{}
	me := CO.MeOrNot(ctx, user)
	var noMssg string

	result := db.Raw("SELECT followBy FROM follow WHERE followTo=? ORDER BY followID DESC", user)
	CO.Err(result.Error)
	rows, err := result.Rows()
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

	user := ctx.Params().Get("id")
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
	var likeBy int
	likes := []interface{}{}

	db.Raw("SELECT COUNT(postID) AS postCount FROM posts WHERE postID=?", post).Row().Scan(&postCount)
	invalid(ctx, postCount)

	result := db.Raw("SELECT likeBy FROM likes WHERE postID=?", post)
	rows, err := result.Rows()
	CO.Err(err)

	for rows.Next() {
		rows.Scan(&likeBy)
		l := map[string]interface{}{
			"likeBy": likeBy,
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
