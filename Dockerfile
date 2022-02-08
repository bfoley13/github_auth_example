FROM golang:1.17rc2-alpine3.14 as build-go

ADD . /server
WORKDIR /server
ENV GOPATH ""
RUN go env -w GOPROXY=direct
RUN apk add git

RUN go mod download
RUN go build -v -o /github_auth_service ./

FROM alpine:3.13

ENV APP_CLIENT_ID ""
ENV APP_SECRET ""
ENV APP_REDIRECT_URI ""

COPY --from=build-go /github_auth_service /github_auth_service
EXPOSE 8080
CMD ["/github_auth_service"]
