FROM alpine:3.6
# TODO: maybe we can use scratch again
MAINTAINER devops-training-bcn@googlegroups.com

ENV binary tcpgoon
ENV install_path /usr/local/bin
COPY ${binary} ${install_path}

ENTRYPOINT ${install_path}/${binary}
CMD ["--help"]
