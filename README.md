# Social Network Web & Mobile App
Social Network Application, Created Iris Go Framework + React Native

# Quick Links
1. [Installation](#installation)
2. [Usage](#usage)
3. [Requirements](#requirements)
4. [Screenshots](#screenshots)
5. [Contribution](#contribution)
6. [Other Projects](#other-projects)

# Installation

### Easy Install

1. Run install.sh (Linux & Mac) or install.bat (Windows)

### Manuel

1. Install Iris and Go ORM Packages manually

```bash
go get -u github.com/kataras/iris
go get -u github.com/jinzhu/gorm
```

2. Install Database Driver

    * For MySQL
    ```bash
    go get -u github.com/go-sql-driver/mysql
    ```
    * For PostgreSQL
    ```bash
    go get -u github.com/jinzhu/gorm/dialects/postgres
    ```
    * For Microsoft SQL Server
    ```bash
    go get -u github.com/jinzhu/gorm/dialects/mssql
    ```
    * For Sqlite
    ```bash
    go get -u github.com/jinzhu/gorm/dialects/sqlite
    ```
    
3. Install other packages

```bash
go get -u golang.org/x/crypto/bcrypt
go get -u github.com/gorilla/sessions
go get -u github.com/joho/godotenv
go get -u github.com/badoux/checkmail
```
    
4. Install npm packages

```bash
npm install
```
    
or

```bash
yarn install
```
    
# Usage

1. Create database for project

2. Edit .env file

```
PORT="Application Port"
DB_PORT="Database Port"
DB_HOST="Database Host"
DB_USER="Database Username"
DB_PASSWORD="Database Passowrd"
DB="Database Name"
DB_TYPE="Database Provider Name"
SESSION_HASH_SECRET="Session Hash Secret"
SESSION_BLOCK_SECRET="Session Block Secret"
```

3. Run application

```bash
go run main.go
```

4. Enjoy!

# Requirements
1. Make sure you keep this project inside `src/` of your Golang project folder ($GOPATH).
2. Following packages should be installed.
    1. [iris](https://github.com/kataras/iris)
    2. [checkmail](https://github.com/badoux/checkmail)
    3. [Gorm](https://github.com/jinzhu/gorm)
    4. [bcrypt](https://golang.org/x/crypto/bcrypt)
    5. [sessions](https://github.com/gorilla/sessions)
    6. [godotenv](https://github.com/joho/godotenv)


# Screenshots

Coming soon

# Contribution

* If you want to contribute to codes, create pull request
* If you find any bugs or error, create an issue

## License

This project is licensed under the MIT License

# Other Projects

## Aktuel Listesi

![Aktuel Listesi](https://pbs.twimg.com/media/Dbe07xRXUAEbOSh.jpg:large)

Aktuel Listesi is smart shopping center product follower. It's mobile application project with backend. Developed with ASP.NET Core Web API and React Native

Backend: [https://github.com/peacecwz/aktuel-listesi](https://github.com/peacecwz/aktuel-listesi)

Mobile App: [https://github.com/peacecwz/aktuel-listesi-app](https://github.com/peacecwz/aktuel-listesi-app)
