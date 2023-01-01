const Kava = require('@mage-labs/javascript-sdk');
const BnbApiClient = require('@binance-chain/javascript-sdk');
const { sleep } = require("./helpers.js");

const setup = async (mageEndpoint, binanceEndpoint, mageMnemonic, binanceMnemonic) => {
    // Start new Kava client
    mageClient = new Kava.KavaClient(mageEndpoint);
    mageClient.setWallet(mageMnemonic);
    mageClient.setBroadcastMode("async");
    await mageClient.initChain();

    // Start Binance Chain client
    const bnbClient = await new BnbApiClient.BncClient(binanceEndpoint);
    bnbClient.chooseNetwork('mainnet');
    const privateKey = BnbApiClient.crypto.getPrivateKeyFromMnemonic(binanceMnemonic);
    bnbClient.setPrivateKey(privateKey);
    bnbClient.chainId = "Binance-Chain-Tigris"
    await bnbClient.initChain()

    // Override the default transaction broadcast endpoint with the tendermint RPC endpoint on the binance-chain node (port 26658 in mgtool)
    bnbClient.setBroadcastDelegate(async(signedTx) => {
      const signedBz = signedTx.serialize()
      console.log(signedBz)
      const opts = {
        params: {tx: "0x" + signedBz},
        headers: {
          "content-type": "application/json",
        },
      }
      var result
      try {
         result = await bnbClient._httpClient.request(
          "get",
          `http://localhost:26658/broadcast_tx_commit`,
          null,
          opts
        )
      } catch (error) {
        console.log(error)
      }

      return result
    })

    return { mageClient: mageClient, bnbClient: bnbClient };
}

const loadKavaDeputies = async (mageClient, assets, amount) => {
    var counter = 0;
    for (var denom in assets) {
        const assetInfo = assets[denom];
        const coins = Kava.utils.formatCoins(amount * assetInfo.conversionFactor, assetInfo.mageDenom);
        const txHash = await mageClient.transfer(assetInfo.mageDeputyHotWallet, coins);
        console.log("Load", assetInfo.mageDenom, "deputy:", txHash);

        counter++;
        counter < Object.keys(assets).length ? await sleep(7000) : await sleep(3000);
    }
    console.log("")
}
const loadBinanceChainDeputies = async (bnbClient, assets, amount) => {
 for (var denom in assets) {
    const assetInfo = assets[denom];
    const res = await bnbClient.transfer(bnbClient.getClientKeyAddress(),
      assetInfo.binanceChainDeputyHotWallet, amount, assetInfo.binanceChainDenom);

    if (res && res.status == 200) {
      console.log('\tLoad', assetInfo.binanceChainDenom, 'deputy:', res.result.result.hash, '\n');
    } else {
      console.log('Tx error:', res);
      return;
    }

    await sleep(2000);
  }
}

module.exports = {
    setup,
    loadKavaDeputies,
    loadBinanceChainDeputies
}
