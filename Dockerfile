#
# k8s-used-api-versions.dockerfile
#
# k8s-used-api-versions checks the status of the used API versions 
# and export them in a Prometheus metrics format
#
# @authors k8s-used-api-versions Maintainers

FROM golang:1.17 as build
ENV GOPATH=/go
ENV PATH="$PATH:$GOPATH/bin"
WORKDIR /app
COPY . /app
RUN go test -v ./pkg/...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager main.go

FROM alpine:3.12
ARG USERNAME=manager
ARG USERID=1003
ARG GROUPID=1003

RUN addgroup -g ${GROUPID} -S manager && \
    adduser -u ${USERID} -S manager -G manager

WORKDIR /
USER manager
COPY --from=build /app/manager .
COPY --from=build /app/config/versions.yaml ./config/versions.yaml
ENTRYPOINT ["/manager"]