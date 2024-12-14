require("@nomicfoundation/hardhat-toolbox");

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.19",
  networks: {
    lifeNetwork: {
      url: "http://127.0.0.1:9650/ext/bc/DCuTeqpQJppqJd97vq1ViWtVxwddrb7cCb9ULAx3pQm5ECaYf/rpc",
      chainId: 9999,
      accounts: [
        "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027" // ewoq private key
      ]
    }
  }
};