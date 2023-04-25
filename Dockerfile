FROM docker.io/library/centos:centos7

LABEL maintainer="junsheng.wu <junsheng.wu@cetccloud.com>"

COPY bin/virtdisk-exporter /usr/local/bin/virtdisk-exporter

ENTRYPOINT [ "/usr/local/bin/virtdisk-exporter" ]

USER root

EXPOSE 9109
