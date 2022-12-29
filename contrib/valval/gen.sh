#!/bin/bash

#peers=()

mkdir -p keys

for i in {1..10}
do
  home=mage-$i

  if [ ! -d $home ]
  then
    mage init val$i --home $home --chain-id magemirror_2221-1 > /dev/null 2>&1

    rm -rf $home/data
    rm $home/config/genesis.json
    cp mage-1/config/init-data-directory.sh $home/config/init-data-directory.sh
    cp mage-1/config/priv_validator_state.json.example $home/config/priv_validator_state.json.example
  fi

  cp $home/config/priv_validator_key.json keys/priv_validator_key_$(($i-1)).json

  #peers+=($(mage tendermint show-node-id --home $home)@$home:26656)
done



