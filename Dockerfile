FROM golang:latest as builder
LABEL maintainer="Utkarsh Bhardwaj (Passeriform) <bhardwajutkars.ub@gmail.com>"
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./web/server.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/server .
EXPOSE 8080
CMD ["./server"]