FROM golang:1.25.0

WORKDIR /app

COPY go.mod go.sum ./

COPY . ./ 

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /my_app main.go

EXPOSE 7540

CMD ["/my_app"]