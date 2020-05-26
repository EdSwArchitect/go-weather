
FROM ubuntu:19.10

RUN apt-get update

RUN mkdir -p /etc/ssl/certs 

COPY ./ca-certificates.crt /etc/ssl/certs

RUN mkdir -p /opt/playground /data /static-data

WORKDIR /opt/playground

COPY ./go-weather go-weather
COPY ./resources/config.json /static-data/config.json

EXPOSE 8080
EXPOSE 18080
VOLUME /data

ENTRYPOINT ["/opt/playground/go-weather", "-configFile", "/data/config.json"]

