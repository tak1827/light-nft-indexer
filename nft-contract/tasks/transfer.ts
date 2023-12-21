import { task } from "hardhat/config";
import * as types from "hardhat/internal/core/params/argumentTypes";
import { loadContract, generateAddresses, generateSerial } from "./utils";

task("transfer", "Transfer nfts")
  .addParam("total", "The total number of nfts to be transferred", 3, types.int)
  .addParam("idFrom", "The token id start from", 0, types.int, true)

  .setAction(async (args, hre) => {
    const { getNamedAccounts } = hre;
    const { deployer } = await getNamedAccounts();

    const nft = await loadContract(hre, "MockNFT");
    console.log(`contract: ${nft.target}, from: ${deployer}`);

    const supply = await nft.totalSupply();
    const initalId = args.idFrom ? args.idFrom : Number(supply);

    // mint
    const tokenIds = generateSerial(args.total, initalId);
    const owners = generateAddresses(args.total);
    let i = 0;
    for (const tokenId of tokenIds) {
      const tx = await nft.mint(owners[i], tokenId, { from: deployer });
      console.log(`[${i + 1}/${args.total}] minted ${tokenId} to ${owners[i]}`);
      i++;
      await tx.wait();
    }

    // transfer
    const tos = generateAddresses(args.total);
    i = 0;
    for (const tokenId of tokenIds) {
      const tx = await nft.transferFrom(owners[i], tos[i], tokenId, {
        from: deployer,
      });
      console.log(
        `[${i + 1}/${args.total}] transferred ${tokenId} to ${tos[i]}`
      );
      i++;
      await tx.wait();
    }
  });
