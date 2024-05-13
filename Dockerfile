FROM golang:1.22 as builder

WORKDIR /opt/psmo

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go vet -v ./...
RUN go build -o /opt/psmo/bin/ ./cmd/...

FROM gcr.io/distroless/static-debian12 as psmo-api

COPY --from=builder /opt/psmo/bin/ /

CMD ["/api"]
