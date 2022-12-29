# kvtool

Assorted dev tools for working with the mage blockchain.

## Installation

```bash
make install
```

## Initialization: kvtool testnet

Note: The current mainnet version of mage is `v0.16.0`. To start a local testnet
with the current mainnet version use `--mage.configTemplate v0.16`. To start a
local testnet with the latest unreleased version, use
`--mage configTemplate master`

Option 1:

The `kvtool testnet bootstrap` command starts a local Mage blockchain as a
background docker container called `generated_magenode_1`. The bootstrap command
only starts the Mage blockchain and Mage REST server services.

```bash
# Start new testnet
kvtool testnet bootstrap --mage.configTemplate master
```

Option 2:

To generate a testnet for mage, binance chain, and a deputy that relays swaps between them:

```bash
# Generate a new kvtool configuration based off template files
kvtool testnet gen-config mage binance deputy --mage.configTemplate master

# Pull latest docker images. Docker must be running.
cd ./full_configs/generated && docker-compose pull

# start the testnet
kvtool testnet up

# When finished with usage, shut down the processes
kvtool testnet down
```

### Flags

Additional flags can be added when initializing a testnet to add additional
services:

`--ibc`: Run Mage testnet with an additional IBC chain

Example:

```bash
# Run Mage testnet with an additional IBC chain
kvtool testnet bootstrap --mage.configTemplate master --ibc
```

`--geth`: Run a go-ethereum node alongside the Mage testnet. The geth node is
initialized with the Mage Bridge contract and test ERC20 tokens. The Mage EVM
also includes Multicall contracts deployed. The contract addresses can be found
on the [kava-labs/kava-bridge](https://github.com/kava-labs/kava-bridge#development)
README.

Example:

```bash
# Run the testnet with a geth node in parallel
kvtool testnet bootstrap --mage.configTemplate master --geth
```

Geth node ports are **not** default, as the Mage EVM will use default JSON-RPC
ports:

Mage EVM RPC Ports:

* HTTP JSON-RPC: `8545`
* WS-RPC port: `8546`

Geth RPC Ports:

* HTTP JSON-RPC: `8555`
* WS-RPC port: `8556`

To connect to the associated Ethereum wallet with Metamask, setup a new network with the following parameters:
* New RPC URL: `http://localhost:8555`
* Chain ID: `88881` (configured from the [genesis](config/templates/geth/initstate/genesis.json#L3))
* Currency Symbol: `ETH`

Finally, connect the mining account by importing the JSON config in [this directory](config/templates/geth/initstate/.geth/keystore)
with [this password](config/templates/geth/initstate/eth-password).

## Usage: kvtool testnet

REST APIs for both blockchains are exposed on localhost:

- Mage: http://localhost:1317
- Binance Chain: http://localhost:8080

You can also interact with the blockchain using the `mage` command line. In a
new terminal window, set up an alias to `mage` on the dockerized mage node and
use it to send a query.

```bash
# Add an alias to the dockerized mage cli
alias dmage='docker exec -it generated_magenode_1 mage'

# Confirm that the alias has been added
alias mage

# For versions before v0.16.x
alias dkvcli='docker exec -it generated_magenode_1 kvcli'
```

You can test the set up and alias by executing a sample query:

```bash
dmage status
dmage q cdp params
```

To send transactions you'll need to recover a user account in the dockerized environment. Valid mnemonics for the blockchains be found in the `config/common/addresses.yaml` file.

```bash
# Recover user account
dmage keys add user --recover
# Enter mnemonic
arrive guide way exit polar print kitchen hair series custom siege afraid shrug crew fashion mind script divorce pattern trust project regular robust safe
```

Test transaction sending by transferring some coins to yourself.

```bash
# Query the recovered account's address
dmage keys show user -a
# Send yourself some coins by creating a send transaction with your address as both sender and receiver
dmage tx bank send [user-address] [user-address] 1000000umage --from user
# Enter 'y' to confirm the transaction
confirm transaction before signing and broadcasting [y/N]:

# Check transaction result by tx hash
dmage q tx [tx-hash]
```

## Shut down: kvtool testnet

When you're done make sure to shut down the kvtool testnet. Always shut down the kvtool testnets before pulling the latest image from docker, otherwise you may experience errors.

```bash
kvtool testnet down
```
