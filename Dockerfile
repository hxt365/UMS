FROM golang:1.16.6

RUN mkdir /code

WORKDIR /code

COPY . /code/

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --build="go build main.go" --command=./main

EXPOSE 8000