# Multi-stage build setup to reduce final image size

# Stage 1: Create a build image to compile our code
FROM golang:alpine as builder
RUN mkdir /build
WORKDIR /build
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main .


# Stage 2: Create final image containing only the go executable
FROM scratch
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]