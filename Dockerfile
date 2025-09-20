FROM golang:1.25.1

WORKDIR /var/www/go

COPY . /var/www/go

RUN go install github.com/githubnemo/CompileDaemon@latest

ENTRYPOINT CompileDaemon --build="go build -o subscr-server ./Service/subscriptionservice/cmd/main.go" --command=./subscr-server
