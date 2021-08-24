FROM ubuntu:20.04

RUN apt-get update && \
    apt-get -y upgrade && \
    apt-get -y install curl jq file vim

VOLUME [ /nibirud ]
WORKDIR /nibirud
EXPOSE 26656 26657
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start"]
STOPSIGNAL SIGTERM

COPY wrapper.sh /usr/bin/wrapper.sh
