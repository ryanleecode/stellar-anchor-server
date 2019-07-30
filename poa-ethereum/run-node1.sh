#!/bin/sh

export BOOTNODE_IP=$(dig +short blockchain-bootnode)
geth --datadir node/ --syncmode 'full' --port 30311 --rpc --rpcaddr '0.0.0.0' --rpccorsdomain '*' --rpcport 8501 --rpcapi 'personal,db,eth,net,web3,txpool,miner' --bootnodes "enode://2b66d070a2914509a515462fcd4bf3076710017246916d03da19976da68d9c1ff25e6b5724b36cb4ca94329da9aecdf386798dada2380391ccc2abd36a354b07@$BOOTNODE_IP:30310" --networkid 1010 --gasprice '0' --unlock 'cdf2899f05075d1abdb7e6476dbfea65be92be9b' --password password.txt --mine --allow-insecure-unlock
