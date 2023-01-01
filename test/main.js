const {KAVA_ENDPOINT_MGTOOL, BINANCE_CHAIN_ENDPOINT_MGTOOL, LOADED_KAVA_MNEMONIC,
    LOADED_BINANCE_CHAIN_MNEMONIC, BEP3_ASSETS } = require("./config.js");
const { setup, loadKavaDeputies } = require("./mgtool.js");
const { incomingSwap, outgoingSwap } = require("./swap.js");

var main = async () => {
    // Initialize clients compatible with mgtool
    const clients = await setup(KAVA_ENDPOINT_MGTOOL, BINANCE_CHAIN_ENDPOINT_MGTOOL,
        LOADED_KAVA_MNEMONIC, LOADED_BINANCE_CHAIN_MNEMONIC);

    // Load each Kava deputy hot wallet
    await loadKavaDeputies(clients.mageClient, BEP3_ASSETS, 100000);

    await incomingSwap(clients.mageClient, clients.bnbClient, BEP3_ASSETS, "busd", 10200005);
    // await outgoingSwap(clients.mageClient, clients.bnbClient, BEP3_ASSETS, "busd", 500005);
};

main();
