FROM alpine:3.6
MAINTAINER devops-training-bcn@googlegroups.com

ENV binary tcpgoon
ENV install_path /usr/local/bin
COPY ${binary} ${install_path}

ENTRYPOINT ["sh", "-c", "${install_path}/${binary}"]
CMD ["--help"]
