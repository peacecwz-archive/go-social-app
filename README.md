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


# Usage

1. Open PHPMyAdmin or any other db tool, create a db & import `db.sql`.

2. Install all the dependencies.
```bash
# with npm
npm install

# or with yarn
yarn
```

3. Create a `.env` file & insert the following code. Replace values with yours!!
```javascript
PORT=YOUR PORT [STRING]
SESSION_HASH_SECRET=SESSION_HASH_KEY [STRING]
SESSION_BLOCK_SECRET=SESSION_BLOCK_KEY [STRING]
DB_USER=DB_USER [STRING]
DB_PASSWORD=DB PASSWORD [STRING]
DB=DB YOU JUST CREATE [STRING]
```

> See the default `.env` file to see an example of it

4. My root folder name is `github.com/peacecwz/go-social-app`, if yours is different then go ahead & change it as it used for imports. It can be done easily by searching the whole project.

5. Now run the app.
```javascript
npm start [OR] yarn start
```

6. Run the app in browser.
```javascript
localhost:[PORT]
```

7. Enjoy 💖💖!!


# Requirements
1. Make sure you keep this project inside `src/` of your Golang project folder ($GOPATH).
2. Following packages should be installed.
    1. [iris](https://github.com/kataras/iris)
    2. [checkmail](https://github.com/badoux/checkmail)
    3. [MySQL driver](https://github.com/go-sql-driver/mysql)
    4. [bcrypt](https://golang.org/x/crypto/bcrypt)
    5. [sessions](https://github.com/gorilla/sessions)
    6. [godotenv](https://github.com/joho/godotenv)


# Screenshots

Coming soon

# Contribution


# Other Projects

