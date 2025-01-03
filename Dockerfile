FROM ubuntu:latest

RUN apt update && apt install -y curl net-tools python3

WORKDIR /app

COPY ./bin/styx ./bin/styx

EXPOSE 8080

COPY ./scripts/entrypoint.sh /entrypoint.sh

RUN chmod +x ./bin/styx /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
