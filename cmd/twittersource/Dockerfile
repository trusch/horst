FROM alpine

RUN apk add --no-cache ca-certificates
ADD twittersource /bin/twittersource

ENTRYPOINT ["/bin/twittersource"]
