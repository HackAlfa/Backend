FROM golang:1.22.3-alpine3.20
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY cache ./cache/
COPY server ./server/
COPY ml ./ml/
RUN CGO_ENABLED=0 GOOS=linux go build -o /backend
RUN chmod +x /backend

EXPOSE 8080

CMD ["/backend"]