FROM golang:alpine

WORKDIR /go/src/wallpaperconstructor
COPY . .

RUN apk update && apk add git
RUN go get -d -v ./...
RUN go install -v ./...

ENTRYPOINT ["wallpaperconstructor"]
CMD ["img.jpeg", "out.jpeg"]