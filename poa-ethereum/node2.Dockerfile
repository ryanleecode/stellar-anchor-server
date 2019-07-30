FROM ethereum/client-go:v1.9.1
RUN apk add bind-tools
WORKDIR /home
COPY ./genesis.json ./
COPY password.txt ./
COPY run-node2.sh ./
RUN chmod +x ./run-node2.sh
RUN mkdir node
COPY node2 node
RUN geth --datadir node init genesis.json
ENTRYPOINT [ "sh", "run-node2.sh" ]
