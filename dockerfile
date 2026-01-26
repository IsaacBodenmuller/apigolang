FROM golang:1.25

WORKDIR /go/src/app

# Certificados base do sistema
RUN apt-get update && apt-get install -y ca-certificates

# Certificado corporativo
COPY corp-ca.crt /usr/local/share/ca-certificates/corp-ca.crt
RUN update-ca-certificates

# Código
COPY . .

# Build
RUN go build -o main cmd/main.go

EXPOSE 8000
CMD ["./main"]