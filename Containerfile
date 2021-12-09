FROM golang:1.17 as build
WORKDIR /jwt-practice
COPY main.go /jwt-practice
COPY go.mod .
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o ./app .

FROM scratch
COPY --from=0 /jwt-practice/app .
EXPOSE 8080
CMD ["/app"]