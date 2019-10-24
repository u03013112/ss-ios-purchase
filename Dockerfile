FROM alpine

RUN apk add --no-cache \ 
    tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo Asia/Shanghai > /etc/timezone && \
    apk del tzdata

COPY --from=golang:1.11.2-alpine3.8 /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY build/ss-ios /usr/local/bin/ss-ios

CMD [ "/usr/local/bin/ss-ios" ]