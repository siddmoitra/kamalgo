# The build stage
FROM golang:1.22 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY app.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo app.go

# The final stage
FROM scratch
LABEL org.opencontainers.image.source=https://github.com/siddmoitra/kamalgo
WORKDIR /app
COPY --from=build /app/app .
WORKDIR /app
EXPOSE 80
CMD ["./app"]