FROM debian:stretch-slim

WORKDIR /

COPY _output/bin/Scheduler-framework /usr/local/bin

CMD ["Scheduler-framework"]
