FROM golang:1.13-alpine
WORKDIR /go/src/editor-js-link-resolver
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
ENV HOST=0.0.0.0 PORT=9000 ALLOW_ORIGIN=*

CMD ["editor-js-link-resolver"]

