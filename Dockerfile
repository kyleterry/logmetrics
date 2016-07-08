FROM golang:1.6
COPY . /go/src/github.com/kyleterry/logmetrics
WORKDIR /go/src/github.com/kyleterry/logmetrics
RUN go build
VOLUME ["/var/log/"]
CMD ["./logmetrics"]
