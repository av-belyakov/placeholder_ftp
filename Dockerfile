FROM golang:1.23.4-alpine AS temporary_image
WORKDIR /go/src/
ENV PATH /usr/local/go/bin:$PATH
RUN apk update && \
    apk add --no-cache git && \
    git clone https://github.com/av-belyakov/placeholder_ftp.git /go/src/
RUN cd cmd/ && go build -o /go/src/placeholder_ftp/placeholder_ftp .

FROM alpine
LABEL user="Artemy" application="placeholder_ftp"
RUN mkdir /opt/placeholder_ftp && \
    mkdir /opt/placeholder_ftp/config && \
    mkdir /opt/placeholder_ftp/logs
WORKDIR /opt/placeholder_ftp
COPY --from=temporary_image /go/src/placeholder_ftp /opt/placeholder_ftp/
COPY --from=temporary_image /go/src/README.md /opt/placeholder_ftp/README.md

ENTRYPOINT [ "./placeholder_ftp" ]

