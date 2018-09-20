#!/bin/bash
for i in {1..10}
do
   /Users/hadv/work/github/hadv/go-ethereum/build/bin/geth --datadir ./ account new --password ./passwd
done
