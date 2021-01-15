FROM golang:1.15.6-alpine3.12 as builder

WORKDIR /src
COPY . .

RUN go mod vendor \
    && go build -o /src/jenkinsallure /src/cmd/main.go


FROM zenika/alpine-chrome:86

USER root
RUN apk add --no-cache dumb-init \
    tzdata \
    busybox-suid \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata \
    && mkdir -p /allure/etc \
    && chown -R chrome:chrome /allure/etc

USER chrome
WORKDIR /usr/src/app
COPY --chown=chrome --from=builder /src/jenkinsallure ./
RUN chmod +x /usr/src/app/jenkinsallure

VOLUME [ "/allure/etc" ]

ENTRYPOINT ["dumb-init", "--"]
CMD ["/usr/src/app/jenkinsallure", "-f", "/allure/etc/config.yml"]
