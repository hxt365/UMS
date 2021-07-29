FROM golang:1.16.6

RUN mkdir /code

WORKDIR /code

COPY . /code/

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT ./wait-for-it.sh ${DATABASE_HOST}:${DATABASE_PORT} && CompileDaemon --build="go build main.go" --command=./main

EXPOSE 8000