#GOLANG development toolbox

A toolbox for golang development and testing


## Code generation

```
dev-tools make:app -n=github.com/AkronimBlack/project
```

The previous command generates scaffolding for a REST api project. The following is made:

```
{app_name}	
	|-api
	|   |- openapi
	|	  |- proto
	|-application
 	|-cmd
	|	   |-{app_name}
	|		     |- main.go
	|		     |- main_test.go
	|-docker
	|    |- Dockerfile
	|	   |- Dockerfile.dev
	|-domain
	|-infrastructure
	|    |-transport
	|    |   |- http  
	|	   |   |- grpc
	|	   |   |- amqp
	|	   |-repositories  
	|-docker-compose.yml
	|-README.md
	|-.env
	|-.env.test
```


```{app_name}``` is the extracted from ```-n=github.com/AkronimBlack/project```. If you follow the recommended form of package naming as is shown in the example
```{app_name}``` is taken as the third part if broken by ```/```. in the example used here ```app_name = project``` 

## Docker

The scaffolding include some basic docker files:

### docker/Dockerfile

```
FROM golang AS builder
LABEL maintainer="AkronimBlack"
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o test
FROM alpine
COPY --from=builder /app/test .
EXPOSE 8080
ENTRYPOINT ["./test"]
	
```

### docker/Dockerfile.dev

```
FROM golang:alpine
RUN apk update && apk upgrade && apk add bash
WORKDIR /app
COPY ./ /app
RUN go mod download
ENTRYPOINT go run main.go
```

### docker-compose.yml
```

version: '3.3'

services:
   test:
      container_name: test
      build: ./docker/Dockerfile.dev
      ports:
        - 8080:8080
      volumes:
        - ./:/app
      depends_on:
        - test_db
      networks:
        - test_network


   test_db:
      image: mysql:5.7
      volumes:
        - test_db_data:/var/lib/mysql
      restart: always
      environment:
        MYSQL_ROOT_PASSWORD: secret
        MYSQL_DATABASE: test
        MYSQL_USER: user
        MYSQL_PASSWORD: user
      ports: 
        - 3306:3306
      networks:
        - test_network

volumes:
   test_db_data: {}
networks:
   test:
     name: test_network
```