FROM golang:alpine AS builder

ARG GO_COMMANDS

# Added certificates to be able to make http request
RUN apk update && apk add --no-cache git curl bash make ca-certificates && rm -rf /var/cache/apk/*
RUN mkdir /usr/share/ca-certificates/extra  && \
cp -R /etc/ssl/certs /usr/local/share/ca-certificates/extra
RUN update-ca-certificates --fresh

# Set necessary environmet variables needed for our image
ENV GO111MODULE= \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN mkdir build

# Set the Current Working Directory inside the container
WORKDIR /build

# create the microservice into the image and copy all from the selected one
COPY . .

# clean up all dependecies and vendor files
RUN ${GO_COMMANDS}

# build and generate the binary
RUN go build -a -o backendify ./main.go

# create a dist folder
RUN mkdir dist

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/backendify ./


# Build a small image
FROM scratch

COPY --from=builder /dist/backendify /
COPY --from=builder /usr/local/share/ca-certificates/extra /etc/ssl/certs/

# Command to run the binary
CMD ["./backendify"]
