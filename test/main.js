const {MAGE_ENDPOINT_KVTOOL, BINANCE_CHAIN_ENDPOINT_KVTOOL, LOADED_MAGE_MNEMONIC,
    LOADED_BINANCE_CHAIN_MNEMONIC, BEP3_ASSETS } = require("./config.js");
const { setup, loadMageDeputies } = require("./kvtool.js");
const { incomingSwap, outgoingSwap } = require("./swap.js");

var main = async () => {
    // Initialize clients compatible with kvtool
    const clients = await setup(MAGE_ENDPOINT_KVTOOL, BINANCE_CHAIN_ENDPOINT_KVTOOL,
        LOADED_MAGE_MNEMONIC, LOADED_BINANCE_CHAIN_MNEMONIC);

    // Load each Mage deputy hot wallet
    await loadMageDeputies(clients.mageClient, BEP3_ASSETS, 100000);

    await incomingSwap(clients.mageClient, clients.bnbClient, BEP3_ASSETS, "busd", 10200005);
    // await outgoingSwap(clients.mageClient, clients.bnbClient, BEP3_ASSETS, "busd", 500005);
};

main();
