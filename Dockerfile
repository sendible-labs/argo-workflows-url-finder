FROM golang:alpine AS builder

WORKDIR $GOPATH/src/app/argo-workflows-url-finder

COPY url-finder .

RUN go get -d -v ./...
RUN go install -v ./...

RUN CGO_ENABLED=0 go build -o /go/bin/argo-workflows-url-finder

FROM scratch

COPY --from=builder /go/bin/argo-workflows-url-finder /go/bin/argo-workflows-url-finder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/go/bin/argo-workflows-url-finder"]