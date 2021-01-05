FROM golang:1.15-alpine
# gcc is needed for the sqlite package
RUN apk add --no-cache gcc musl-dev

WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
CMD ["gossip"]