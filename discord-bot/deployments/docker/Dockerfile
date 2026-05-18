FROM golang:1.26.3

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o pcep-disgo ./cmd/app/main.go

CMD ["./pcep-disgo"]