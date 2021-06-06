FROM golang:1.16.5 as builder

ENV GO111MODULE=on

ARG BEAVER_VERSION=1.2.3

RUN mkdir -p $GOPATH/src/github.com/clivern

RUN git clone -b master https://github.com/Clivern/Beaver.git $GOPATH/src/github.com/clivern/beaver

WORKDIR $GOPATH/src/github.com/clivern/beaver

RUN git checkout tags/$BEAVER_VERSION

RUN go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/beaver .

RUN mkdir -p /go/logs/beaver
RUN mkdir -p /go/configs/beaver

# Build a small image
FROM alpine:3.12
RUN apk --no-cache add ca-certificates

# Copy our static executable
COPY --from=builder /go/bin/beaver /go/bin/beaver
COPY --from=builder /go/logs/beaver /go/logs/beaver
COPY --from=builder /go/configs/beaver /go/configs/beaver

WORKDIR /go/bin/

EXPOSE 8080

CMD ["./beaver"]