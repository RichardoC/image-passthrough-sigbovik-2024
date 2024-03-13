# ARG GO_VERSION=1
FROM golang:1.22-alpine3.19 as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /run-app .


FROM alpine:3.19

COPY --from=builder /run-app /usr/local/bin/

# USER 37100

CMD ["run-app"]
