FROM golang:alpine

WORKDIR /go/src/wallpaperconstructor/web/
COPY . .

ENV GOOS=js
ENV GOARCH=wasm

RUN apk update && apk add git
RUN go get -d -v ./...
RUN go build -v -o main.wasm

FROM httpd:alpine

COPY main.wasm /usr/local/apache2/htdocs/
WORKDIR /usr/local/apache2/htdocs/

ENTRYPOINT ["wallpaperconstructor"]
CMD ["img.jpeg", "out.jpeg"]