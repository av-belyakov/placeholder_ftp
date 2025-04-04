# Dockerfile собирать с аргументом --build-arg
# sudo docker build tag gitlab.cloud.gcm:5050/a.belyakov/thehivehook_go_package:test_image --build-arg VERSION=v0.3.2 .
# 
# для удаления временного образа, можно через ci/cd, можно руками 
# docker image prune -a --force --filter="label=temporary"

FROM golang:1.24.1-alpine AS packages_image
WORKDIR /go/src
COPY go.mod go.sum ./
RUN echo 'packages_image' && \
    go mod download

FROM golang:1.24.1-alpine AS build_image
LABEL temporary=''
ARG BRANCH
ARG VERSION
WORKDIR /go/
COPY --from=packages_image /go ./
RUN echo -e "build_image" && \
    rm -r ./src && \
    apk update && \
    apk add --no-cache git && \
    git clone -b ${BRANCH} https://github.com/av-belyakov/placeholder_ftp.git ./src/${VERSION}/ && \
    go build -C ./src/${VERSION}/cmd/ -o ../app

FROM alpine
LABEL author='Artemij Belyakov'
ARG STATUS
ARG VERSION
ARG USERNAME=dockeruser
ARG US_DIR=/opt/placeholder_ftp
ENV GO_PHFTP_MAIN=${STATUS}
RUN addgroup --g 1501 groupcontainer
RUN adduser -u 1501 -G groupcontainer -D ${USERNAME} --home ${US_DIR}
USER ${USERNAME}
WORKDIR ${US_DIR}
RUN mkdir ./logs
COPY --from=build_image /go/src/${VERSION}/app ./
COPY --from=build_image /go/src/${VERSION}/README.md ./
COPY --from=build_image /go/src/${VERSION}/internal/appversion/version ./  
COPY config/* ./config/

ENTRYPOINT [ "./app" ]

