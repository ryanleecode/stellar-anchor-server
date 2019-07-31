FROM ubuntu:18.04

RUN apt-get update
RUN apt-get install -y software-properties-common
RUN add-apt-repository -y ppa:ethereum/ethereum
RUN apt-get update
RUN apt-get install -y ethereum
WORKDIR /home
COPY genesis.json ./
COPY boot.key ./
EXPOSE 30310
CMD [ "bootnode", "-nodekey", "boot.key", "-verbosity", "9", "-addr", "0.0.0.0:30310" ]