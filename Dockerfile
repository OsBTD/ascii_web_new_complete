FROM golang:1.22.4-alpine

LABEL maintainer1="borntolearn314@gmail.com"
LABEL maintainer2="soufiane@gmail.com"
LABEL version="1.0"
LABEL description="Docker container for ASCII Art Web app"
LABEL author1="Oussama"
LABEL author2="Soufiane"


WORKDIR /ascii-art-web-stylize-NEW

COPY go.mod .
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]