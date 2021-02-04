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
	