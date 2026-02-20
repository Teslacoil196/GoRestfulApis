FROM golang:1.26.0-trixie AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go
#RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.1/migrate.linux-amd64.tar.gz | tar xvz 
RUN ls -al
      

FROM alpine:latest 
WORKDIR /app
COPY --from=builder app/main .
#THE 'E' fixed it -- imagine the markplier meme now 
COPY --from=builder app/migrate ./migratE 
COPY app.env .
COPY --chmod=755 start.sh .
#COPY start.sh .
COPY wait-for.sh .
COPY db/migrate ./migrate

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]