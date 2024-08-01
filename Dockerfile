FROM golang:1.22-alpine

WORKDIR /usr/src/app

COPY . .

RUN go mod download

RUN go build -o /go_final_project ./cmd/TODOServer/todo_server.go

CMD ["/go_final_project"]