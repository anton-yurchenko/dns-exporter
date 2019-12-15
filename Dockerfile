FROM golang:1.13.1-alpine as build
WORKDIR /opt/src
COPY . .
RUN go install &&\
    go build -o /opt/release

FROM golang:1.13.1-alpine
LABEL "repository"="https://github.com/anton-yurchenko/dns-exporter"
LABEL "maintainer"="Anton Yurchenko <anton.doar@gmail.com>"
LABEL "version"="1.0.0"
RUN addgroup -S app &&\
    adduser -S app -G app
COPY --chown=app:app --from=build /opt/release /opt/dns-exporter
USER app
ENTRYPOINT [ "/opt/dns-exporter" ]