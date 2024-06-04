FROM golang AS builder

RUN adduser --disabled-password --gecos "" --home "/nonexistent" --shell "/sbin/nologin" --no-create-home "app"

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .


FROM scratch

WORKDIR /
COPY --from=builder /app/app .
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

ENV PORT 8080
ENV HELLO_MSG Hello
USER app:app

CMD ["/app"]