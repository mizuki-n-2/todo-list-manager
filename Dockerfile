FROM golang:1.20-alpine

WORKDIR /usr/src/app

#COPY go.mod  go.sum ./   ## go.sumファイルが生成されたらこっちを使う
COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

CMD ["app"]