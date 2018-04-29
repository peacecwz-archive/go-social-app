

#!/bin/sh
echo "Hello, Installing dependencies for Go Social Web & Mobile App"
echo "---------------------------------------------------------------"
echo "Installing Go Iris Framework (github.com/kataras/iris)"
go get -u github.com/kataras/iris
echo "Installed"

echo "---------------------------------------------------------------"
echo "Installing Go ORM (github.com/jinzhu/gorm)"
go get -u github.com/jinzhu/gorm
echo "Installed"

echo "---------------------------------------------------------------"
read -p "What's your DBMS Provider? (postgresql,mysql,mssql etc...): " dbProvider

case "$dbProvider" in
"mysql")
    go get -u github.com/go-sql-driver/mysql
    echo "Installed MySQL Driver"
    ;;
"postgresql")
    go get -u github.com/jinzhu/gorm/dialects/postgres
    echo "Installed PostgreSQL Driver"
    ;;
"mssql")
    go get -u github.com/jinzhu/gorm/dialects/mssql
    echo "Installed Microsoft SQL Server Driver"
    ;;
"sqlite")
    go get -u github.com/jinzhu/gorm/dialects/sqlite
    echo "Installed SQLite Driver"
    ;;
*)
    echo "Cannot find your Database Provider. Please install Database Provider Driver manuelly"
    ;;
esac

echo "---------------------------------------------------------------"
echo "Installed all packages successfully"