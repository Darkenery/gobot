FROM golang:latest as builder

WORKDIR /go/src/github.com/darkenery/gobot
COPY . .

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure -v

RUN GOOS=linux GOARCH=arm go build cmd/main.go
RUN cp config.yaml.dist config.yaml

FROM easypi/alpine-arm:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/src/github.com/darkenery/gobot/main .
COPY --from=builder /go/src/github.com/darkenery/gobot/config.yaml .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]


