FROM scratch
MAINTAINER devops-training-bcn@googlegroups.com
COPY out/tcpgoon /
ENTRYPOINT ["/tcpgoon"]
CMD ["--help"]
