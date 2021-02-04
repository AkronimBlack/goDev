package templates

/*MainTemplate stub for cmd/main.go*/
func MainTemplate() []byte {
	return []byte(`package main

func main() {

}
`)
}

/*MainTestTemplate stub for cmd/main.go*/
func MainTestTemplate() []byte {
	return []byte(`package main_test

func main() {

}
`)
}

/*DockerComposeTemplate stub for generic docker-compose.yml*/
func DockerComposeTemplate() []byte {
	return []byte(`
version: '3.3'

services:
   {{.Name}}:
      container_name: {{.Name}}
      build: ./docker/Dockerfile.dev
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
   {{.Name}}:
     name: {{.Name}}_network`)
}

/*DockerfileDevTemplate stub for generic docker-compose.yml*/
func DockerfileDevTemplate() []byte {
	return []byte(`FROM golang:alpine
RUN apk update && apk upgrade && apk add bash
WORKDIR /app
COPY ./ /app
RUN go mod download
ENTRYPOINT go run main.go
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
	`)
}

/*GoSumTemplate stub for generic docker-compose.yml*/
func GoSumTemplate() []byte {
	return []byte(`
	`)
}
