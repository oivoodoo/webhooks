FROM golang:1.12-alpine

COPY . ${GOPATH}/pkg/gitlab.com/oivoodoo/webhooks

WORKDIR ${GOPATH}/pkg/gitlab.com/oivoodoo/webhooks

EXPOSE 5000
