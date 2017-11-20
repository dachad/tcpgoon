FROM scratch
MAINTAINER devops-training-bcn@googlegroups.com

ENV binary out/tcpgoon
ENV install_path /usr/local/bin
COPY ${binary} ${install_path}

ENTRYPOINT ${install_path}/${binary}
CMD ["--help"]
