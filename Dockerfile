#first stage - builder
FROM golang:1.24  as builder
WORKDIR /app
ARG GITLAB_ACCESS_TOKEN
RUN git config --global url."https://oauth2:${GITLAB_ACCESS_TOKEN}@git.larsa.io/".insteadOf "https://git.larsa.io/"
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./
#RUN go get -insecure git.larsa.org/mahdawi/microservices-commons.git
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

#second stage
FROM alpine:3.14
RUN apk add --no-cache \
        wkhtmltopdf \
        xvfb \
        ttf-dejavu ttf-droid ttf-freefont ttf-liberation \
    ;
WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN chown -R appuser:appgroup /app

COPY --from=builder /app/main .
COPY --from=builder  --chown=appuser:appgroup  /app/internal/modules/seed/data/ /app/internal/modules/seed/data/
COPY --from=builder /app/templates/ ./templates

USER appuser
CMD ["./main"]