FROM golang:alpine as build

RUN apk update && \ 
  apk add git build-base

WORKDIR /go/src/app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY . .

RUN go test . gcp-iap-auth/jwt && \
  go build -o /usr/local/bin/gcp-iap-auth

FROM scratch

COPY --from=build /usr/local/bin/gcp-iap-auth /usr/local/bin/gcp-iap-auth

EXPOSE 80 443
ENTRYPOINT ["/usr/local/bin/gcp-iap-auth"]
