FROM alpine
RUN apk add --no-cache ca-certificates
ADD ./slackbot /
CMD ["/slackbot"]
