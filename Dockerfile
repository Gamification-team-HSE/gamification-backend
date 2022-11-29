FROM golang:1.19-alpine

ARG DB_PASSWORD
ARG JWT_SECRET
ARG SUPER_ADMIN_EMAIL
ARG SMTP_PASSWORD

ENV DB_PASSWORD $DB_PASSWORD
ENV SMTP_PASSWORD $SMTP_PASSWORD
ENV SUPER_ADMIN_EMAIL $SUPER_ADMIN_EMAIL
ENV JWT_SECRET $JWT_SECRET

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./ ./

COPY /migrations/*.sql /app/migrations

# Build
RUN go build -o /cmd/server ./cmd/server/main.go

# Run
CMD [ "/cmd/server" ]