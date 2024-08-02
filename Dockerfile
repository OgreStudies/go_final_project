FROM golang:1.22-alpine

EXPOSE 7540

WORKDIR /usr/src/app

COPY . .

RUN go mod download

RUN go build -o /go_final_project ./cmd/main.go

CMD ["/go_final_project"]