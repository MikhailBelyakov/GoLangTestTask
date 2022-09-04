FROM golang:1.13
ADD . /go/src
WORKDIR /go/src

RUN go build -o bin/main  main.go
RUN git clone https://github.com/vishnubob/wait-for-it.git

CMD ./wait-for-it/wait-for-it.sh --host=db --port=3306 --timeout=60 -- /go/src/bin/main

EXPOSE 8080

