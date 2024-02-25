FROM golang:1.19-alpine as builder

# Setup
RUN mkdir -p /app
WORKDIR /app

# Add libraries
RUN apk add --no-cache git

# Copy & build
ADD . /app
RUN mkdir build 
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -o build/app ./main.go

# Copy into scratch container
FROM alpine
COPY --from=builder /app/build/app ./
ENTRYPOINT ["./app"]