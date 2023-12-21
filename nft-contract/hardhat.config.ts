import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";
import "hardhat-deploy";

import "./tasks/transfer";

const config: HardhatUserConfig = {
  networks: {
    hardhat: {
    },
    localhost: {
      url: "http://127.0.0.1:8545"
    },
  },
  namedAccounts: {
		deployer: 0
	},
  solidity: "0.8.20",
};

export default config;
