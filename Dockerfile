FROM golang:1.15.6-alpine3.12 as builder

WORKDIR /app
COPY . .

RUN go build -o /app/build/jenkinsallure /app/cmd/main.go 

FROM chromedp/headless-shell:88.0.4324.87 as prod

WORKDIR /app

COPY --from=builder /app/build/jenkinsallure .
RUN apt install dumb-init && chmod +x /app/jenkinsallure

VOLUME ["/config"]

ENTRYPOINT ["dumb-init", "--"]

CMD ["/app/jenkinsallure", "-h"]
