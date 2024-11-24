FROM golang:1.23.2-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o verve cmd/main.go

FROM scratch

COPY --from=builder /app/verve /

CMD [ "./verve" ]
