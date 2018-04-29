package routes

import (
	"fmt"
	"os"
	"strconv"

	"github.com/peacecwz/go-social-app/models"

	CO "github.com/peacecwz/go-social-app/config"

	"github.com/badoux/checkmail"
	"github.com/kataras/iris"
	"golang.org/x/crypto/bcrypt"
)

// Signup route
func Signup(ctx iris.Context) {
	notLoggedIn(ctx)
	renderTemplate(ctx, "signup", iris.Map{
		"title": "Signup For Free",
	})
}

// Login route
func Login(ctx iris.Context) {
	notLoggedIn(ctx)
	renderTemplate(ctx, "login", iris.Map{
		"title": "Login To Continue",
	})
}

// Logout route
func Logout(ctx iris.Context) {
	loggedIn(ctx, "")
	session := CO.GetSession(ctx)
	session.Destroy()
	// or
	// session.Delete("id")
	// session.Delete("username")
	// if just delete the values.
	ctx.Redirect("/login")
}

// UserSignup function to register user
func UserSignup(ctx iris.Context) {
	fmt.Println("UserSignup")
	resp := make(map[string]interface{})
	username := ctx.PostValueTrim("username")
	email := ctx.PostValueTrim("email")
	password := ctx.PostValueTrim("password")
	passwordAgain := ctx.PostValueTrim("password_again")

	mailErr := checkmail.ValidateFormat(email)

	CO.Err(mailErr)

	db := CO.DB()
	var userCount int
	model := db.Model(&(models.User{}))
	model.Where("username=? OR email=?", username, email).Count(&userCount)

	if username == "" || email == "" || password == "" || passwordAgain == "" {
		resp["mssg"] = "Some values are missing!!"
	} else if len(username) < 4 || len(username) > 32 {
		resp["mssg"] = "Username should be between 4 and 32"
	} else if mailErr != nil {
		resp["mssg"] = "Invalid email format!!"
	} else if password != passwordAgain {
		resp["mssg"] = "Passwords don't match"
	} else if userCount > 0 {
		resp["mssg"] = "Username or Email already exists!!"
	} else {
		user := &(models.User{
			Username: username,
			Email:    email,
			Password: string(hash(password))})
		model.Create(user)

		insStr := strconv.FormatUint(uint64(user.ID), 10)

		dir, _ := os.Getwd()
		userPath := dir + "/public/users/" + insStr
		if CO.FileExists(userPath) == false {
			dirErr := os.Mkdir(userPath, 0655)
			CO.Err(dirErr)
		}

		if CO.FileExists(userPath+"/avatar.png") != false {
			linkErr := os.Link(dir+"/public/images/golang.png", userPath+"/avatar.png")
			CO.Err(linkErr)
		}

		session := CO.GetSession(ctx)
		session.Set("id", insStr)
		session.Set("username", username)

		resp["success"] = true
		resp["mssg"] = "Hello, " + username + "!!"

	}
	json(ctx, resp)
}

// UserLogin function to log user in
func UserLogin(ctx iris.Context) {
	resp := make(map[string]interface{})

	rusername := ctx.PostValueTrim("username")
	rpassword := ctx.PostValueTrim("password")

	db := CO.DB()

	userModel := db.Model(&(models.User{}))
	var user models.User
	userModel.Where("username = ?", rusername).Find(&user)

	encErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(rpassword))

	if rusername == "" || rpassword == "" {
		resp["mssg"] = "Some values are missing!!"
	} else if user.ID == 0 {
		resp["mssg"] = "Invalid username!!"
	} else if encErr != nil {
		resp["mssg"] = "Invalid password!!"
	} else {
		session := CO.GetSession(ctx)
		session.Set("id", strconv.Itoa(int(user.ID)))
		session.Set("username", user.Username)
		resp["mssg"] = "Hello, " + user.Username + "!!"
		resp["success"] = true
	}
	json(ctx, resp)
}
