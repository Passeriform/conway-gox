FROM golang:latest as builder
LABEL maintainer="Utkarsh Bhardwaj (Passeriform) <bhardwajutkars.ub@gmail.com>"
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -installsuffix cgo -ldflags '-extldflags "-static"' -o server ./web/server.go

FROM scratch
WORKDIR /app
ENV ENVIRONMENT PRODUCTION
ENV PORT 8080
COPY --from=builder /build/server .
COPY --from=builder /build/web/static ./static
COPY --from=builder /build/web/templates ./templates
EXPOSE 8080
CMD ["./server"]