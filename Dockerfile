# --- build the Vue3 web UI ---
FROM node:22-alpine AS web
WORKDIR /web
COPY web/package.json ./
RUN npm install
COPY web/ ./
RUN npm run build

# --- build the Go server (pure-Go SQLite, no CGO) ---
FROM golang:1.24-alpine AS server
WORKDIR /src
COPY server/ ./
RUN go mod tidy && CGO_ENABLED=0 go build -o /promptops .

# --- final runtime image ---
FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=server /promptops /app/promptops
COPY --from=web /web/dist /app/web/dist
EXPOSE 8080
CMD ["/app/promptops"]
