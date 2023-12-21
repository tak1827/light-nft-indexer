#!/bin/sh

echo $PATH

# write the genesis block
if [ ! -d "/root/.ethereum/geth" ]; then  
  geth init $GENESIS
fi

# sync only mode if the `ETHERBASE` is not set
if [ -n "$ETHERBASE" ]; then
  OPTS="$OPTS --mine --miner.etherbase $ETHERBASE --miner.gaslimit 30000000"
  OPTS="$OPTS --keystore $KEYSTORE --unlock $ETHERBASE --password /dev/null --allow-insecure-unlock"
fi

exec geth \
  --keystore $KEYSTORE \
  --syncmode full --gcmode archive --networkid $NETWORK_ID \
  --http --http.addr 0.0.0.0 --http.vhosts '*' --http.corsdomain '*' --http.api net,eth,web3,txpool,debug,admin \
  --ws --ws.addr 0.0.0.0 --ws.origins '*' --ws.api net,eth,web3,txpool,debug,admin \
  $OPTS $@
