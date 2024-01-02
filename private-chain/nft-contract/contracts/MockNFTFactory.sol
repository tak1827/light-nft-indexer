//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./MockNFT.sol";

contract MockNFTFactory {
    event NFTCreated(address indexed nft, address creator);

    constructor() {}

    function create() public {
        IERC721 token = new MockNFT();
        emit NFTCreated(address(token), msg.sender);
    }
}
