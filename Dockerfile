FROM alpine:3.17
RUN apk add --no-cache ca-certificates
COPY vtv .
EXPOSE 8043
ENTRYPOINT [ "./vtv" ]