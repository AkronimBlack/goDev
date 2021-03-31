package templates

/*MainTemplate stub for cmd/main.go*/
func MainTemplate() []byte {
	return []byte(`package main
  
import (
    "github.com/joho/godotenv"
    "log"
)

func main() {
    err := godotenv.Load()
    if err != nil {
      log.Fatal("Error loading .env file")
    }
}
`)
}

/*MainTestTemplate stub for cmd/main.go*/
func MainTestTemplate() []byte {
	return []byte(`package main_test
  
import (
  "github.com/joho/godotenv"
  "log"
)

func main() {
    err := godotenv.Load(".env.test")
    if err != nil {
      log.Fatal("Error loading .env.test file")
    }
}
`)
}

/*DockerComposeTemplate stub for generic docker-compose.yml*/
func DockerComposeTemplate() []byte {
	return []byte(`version: '3.3'

services:
   {{.Name}}:
      container_name: {{.Name}}
      build: ./
      ports:
        - 8080:8080
      volumes:
        - ./:/app
      depends_on:
        - {{.Name}}_db
      networks:
        - {{.Name}}_network


   {{.Name}}_db:
      image: mysql:5.7
      volumes:
        - {{.Name}}_db_data:/var/lib/mysql
      restart: always
      environment:
        MYSQL_ROOT_PASSWORD: secret
        MYSQL_DATABASE: {{.Name}}
        MYSQL_USER: user
        MYSQL_PASSWORD: user
      ports: 
        - 3306:3306
      networks:
        - {{.Name}}_network

volumes:
   {{.Name}}_db_data: {}
networks:
   {{.Name}}_network:`)
}

/*DockerfileDevTemplate stub for generic docker-compose.yml*/
func DockerfileDevTemplate() []byte {
	return []byte(`FROM golang:alpine
RUN apk update && apk upgrade && apk add bash
WORKDIR /app
COPY ./ /app
RUN go mod download
ENTRYPOINT go run cmd/{{.Name}}/main.go
	`)
}

/*DockerfileTemplate stub for generic docker-compose.yml*/
func DockerfileTemplate() []byte {
	return []byte(`FROM golang AS builder
LABEL maintainer="{{.Maintainer}}"
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o {{.Name}}
FROM alpine
COPY --from=builder /app/{{.Name}} .
EXPOSE 8080
ENTRYPOINT ["./{{.Name}}"]
	`)
}

/*GoModTemplate stub for generic docker-compose.yml*/
func GoModTemplate() []byte {
	return []byte(`module {{.FullName}}

go 1.15
require (
    github.com/joho/godotenv v1.3.0
)`)
}

/*GoSumTemplate stub for generic docker-compose.yml*/
func GoSumTemplate() []byte {
	return []byte(`
	`)
}

/*GinTemplate stub for generic gin main.go file*/
func GinTemplate() []byte {
	return []byte(`package main

import (
  "fmt"
  "io"
  "log"
  "os"
  "time"

  "github.com/gin-contrib/cors"
  "github.com/gin-gonic/gin"
  "github.com/joho/godotenv"
)

var router *gin.Engine

func main() {
  router = gin.New()

  router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
    return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
      param.ClientIP,
      param.TimeStamp.Format(time.RFC1123),
      param.Method,
      param.Path,
      param.Request.Proto,
      param.StatusCode,
      param.Latency,
      param.Request.UserAgent(),
      param.ErrorMessage,
    )
  }))
  logFile, err := os.OpenFile("logs/{{.Name}}.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

  buildDependencies()

  gin.DefaultWriter = io.MultiWriter(os.Stdout, logFile)
  router.Use(gin.Recovery())

  config := cors.Config{
    AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
    AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
    AllowAllOrigins:  true,
  }
  registerRoutes()
  router.Use(cors.New(config))
  err = godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
  log.Fatal(router.Run(":8080"))
}

func registerRoutes() {}

func buildDependencies(){}
`)
}

func EnvTemplate() []byte {
	return []byte(`HOST=0.0.0.0
PORT=8080`)
}

func GitIgnoreTemplate() []byte {
	return []byte(`.vscode
.idea
/logs
/vendor
.env`)
}

func MigrateTemplate() []byte {
	return []byte(`package repositories

func Migrate(conn *DatabaseConnection) error {
  var err error
  return err
}`)
}

func ConnectionTemplate() []byte {
	return []byte(`package repositories

  import (
    "fmt"
    "log"
  
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
  )

func NewDatabaseConnection(driver, user, password, hostname, port, database string, debug bool) *DatabaseConnection {
  dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, hostname, port, database)

  if driver == "postgres" {
    dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", hostname, port, user, database, password)
  }

  db, err := gorm.Open(mysql.Open(dsn))

  if err != nil {
    log.Println("Unable to connect to database")
    log.Panic(err.Error())
  }

  if debug {
    return NewDatabaseConnectionWithDB(db.Debug())
  }
  return NewDatabaseConnectionWithDB(db)
}

// NewDatabaseConnection constructor for DatabaseConnection
func NewDatabaseConnectionWithDB(db *gorm.DB) *DatabaseConnection {
  return &DatabaseConnection{
    db: db,
  }
}

// DatabaseConnection ....
type DatabaseConnection struct {
  db *gorm.DB
}

//GetConnection returns new gorm.DB connection. 
func (r *DatabaseConnection) GetConnection() *gorm.DB {
  //return r.db.Session(&gorm.Session{FullSaveAssociations: true})
  return r.db
}

//AddWith append preloads to the query
func (r *DatabaseConnection) AddWith(db *gorm.DB, with []string) *gorm.DB {
  for _, w := range with {
    db = db.Preload(w)
  }
  return db
}

//GetConnectionWithPreload get db connection with preloads
func (r *DatabaseConnection) GetConnectionWithPreload(with []string) *gorm.DB {
  if with == nil {
    with = make([]string, 0)
  }
  db := r.GetConnection()
  for _, w := range with {
    db = db.Preload(w)
  }
  return db
}`)
}
