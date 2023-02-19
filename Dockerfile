FROM golang:1.19.5

COPY go.mod go.sum /go/src/github.com/PSNAppz/Fold-ELK/
WORKDIR /go/src/github.com/PSNAppz/Fold-ELK
RUN mv .env.example .env
ENV $(cat /path/to/.env | xargs)
RUN go mod download

COPY . /go/src/github.com/PSNAppz/Fold-ELK
RUN go build -o /usr/bin/fold-elk github.com/PSNAppz/Fold-ELK/cmd/api

EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/fold-elk"]

