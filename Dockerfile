FROM golang:1.19-alpine

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

# Build
RUN go build -o /cmd/server ./cmd/server/main.go

# Run
CMD [ "/cmd/server" ]