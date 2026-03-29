# syntax=docker/dockerfile:1
FROM golang:1.24.1

WORKDIR /app
# копируем наш статический бинарник
# COPY cmd/main.go .

# RUN go build -o cnc_manager cmd/main.go 

COPY . .
RUN go mod download

RUN go build -o cnc_manager ./cmd
# указываем команду запуска
CMD ["./cnc_manager"]