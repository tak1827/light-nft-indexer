//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";

contract MockNFT is ERC721 {
    uint256 private _tokenCount = 0;

    constructor() ERC721("MockEBL", "MEBL") {}

    function mint(address to, uint256 tokenId) public {
        _tokenCount += 1;
        _mint(to, tokenId);
    }

    function transferFrom(
        address from,
        address to,
        uint256 tokenId
    ) public override {
        // NOTE: omit approval
        // require(_isApprovedOrOwner(_msgSender(), tokenId), "ERC721: transfer caller is not owner nor approved");
        _transfer(from, to, tokenId);
    }

    function totalSupply() external view returns (uint256) {
        return _tokenCount;
    }
}
