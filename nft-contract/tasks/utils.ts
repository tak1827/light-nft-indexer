import { Contract, Wallet, randomBytes } from "ethers";
import { HardhatRuntimeEnvironment } from "hardhat/types";

export const ZERO_ADDRESS = "0x0000000000000000000000000000000000000000";

type GeneratorFunction<T> = () => T;

const generateData = <T>(
  size: number,
  generator: GeneratorFunction<T>
): T[] => {
  const result = [];
  for (let i = 0; i < size; i++) {
    result.push(generator());
  }
  return result;
};

export const generateAddresses = (size: number): string[] =>
  generateData<string>(size, () =>
    Wallet.createRandom().address.toLocaleLowerCase()
  );

export const generateSerial = (size: number, base: number): number[] =>
  generateData<number>(
    size,
    (
      (_base) => () =>
        _base++
    )(base)
  );

export const generateRandomBytes = (size: number, l: number): Uint8Array[] =>
  generateData<Uint8Array>(size, () => randomBytes(l));

export const loadContract = async (
  hre: HardhatRuntimeEnvironment,
  name: string
): Promise<Contract> => {
  const { deployments, ethers } = hre;
  const contract = await deployments.get(name);
  return await ethers.getContractAt(name, contract.address);
};
