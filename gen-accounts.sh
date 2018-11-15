#!/bin/bash
for i in {1..10}
do
    ~/geth-linux-amd64 account new --password <(echo password)
done
