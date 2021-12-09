FROM golang:1.17 as build
WORKDIR /go/src/github.com/rcarrata/rck-auth/
ADD cmd /go/src/github.com/rcarrata/rck-auth/cmd
ADD pkg /go/src/github.com/rcarrata/rck-auth/pkg
ADD go.mod go.sum /go/src/github.com/rcarrata/rck-auth/
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rck-auth cmd/rck-auth/main.go

FROM scratch
COPY --from=0 /go/src/github.com/rcarrata/rck-auth/rck-auth .
EXPOSE 8080
CMD ["/rck-auth"]