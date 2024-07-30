FROM golang:1.22-alpine

WORKDIR /usr/src/app

COPY . .

RUN go mod tidy 

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go_final_project

CMD ["/go_final_project"] 