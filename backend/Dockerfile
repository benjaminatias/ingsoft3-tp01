# Etapa 1: compilación
FROM golang:1.25-alpine AS build

WORKDIR /app

# Primero los módulos, para aprovechar la caché de Docker.
COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

# Binario estático para poder utilizar una imagen final pequeña.
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api ./cmd/api

# Etapa 2: tests (parte del build: ya tiene Go y el código compilado).
# El pipeline construye esta etapa y la corre; la imagen final NO se lleva nada de esto.
FROM build AS test
COPY tests ./tests
COPY scripts/cobertura.sh /usr/local/bin/cobertura.sh
ENTRYPOINT ["sh", "/usr/local/bin/cobertura.sh"]

# Etapa 3: imagen final
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Solamente se copia el binario compilado.
COPY --from=build /app/api /app/api

EXPOSE 8080

CMD ["/app/api"]
