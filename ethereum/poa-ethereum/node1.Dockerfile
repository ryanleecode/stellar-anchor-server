FROM ethereum/client-go:v1.9.1
RUN apk add bind-tools
WORKDIR /home
COPY ./genesis.json ./
COPY password.txt ./
COPY run-node1.sh ./
RUN chmod +x ./run-node1.sh
RUN mkdir node
COPY node1 node
RUN geth --datadir node init genesis.json
ENTRYPOINT [ "sh", "run-node1.sh" ]