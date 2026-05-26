// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Strings.sol";

contract LiveHouseTicket is ERC721URIStorage, Ownable {
    using Strings for uint256;

    uint256 private _nextTokenId;

    struct TicketData {
        string eventName;
        string venueName;
        uint256 eventDate;
        string ticketType;
        bool isPOAP;
        bool exists;
    }

    mapping(uint256 => TicketData) private _ticketData;
    mapping(string => bool) private _usedCodes;

    event TicketMinted(uint256 tokenId, address indexed to, string code);
    event POAPClaimed(uint256 tokenId, address indexed claimer);

    constructor() ERC721("LiveHouseTicket", "LHT") Ownable(msg.sender) {}

    function mintTicket(
        address to,
        string memory code,
        string memory uri,
        TicketData memory data
    ) external onlyOwner returns (uint256) {
        require(!_usedCodes[code], "Code already used");
        require(to != address(0), "Invalid address");

        _usedCodes[code] = true;

        uint256 tokenId = _nextTokenId;
        _nextTokenId++;

        _safeMint(to, tokenId);
        _setTokenURI(tokenId, uri);

        data.exists = true;
        _ticketData[tokenId] = data;

        emit TicketMinted(tokenId, to, code);
        return tokenId;
    }

    function claimPOAP(uint256 tokenId) external {
        require(_ownerOf(tokenId) != address(0), "Token does not exist");
        require(ownerOf(tokenId) == msg.sender, "Not token owner");
        require(!_ticketData[tokenId].isPOAP, "Already claimed POAP");

        _ticketData[tokenId].isPOAP = true;

        emit POAPClaimed(tokenId, msg.sender);
    }

    function getTicketData(uint256 tokenId) external view returns (TicketData memory) {
        require(_ownerOf(tokenId) != address(0), "Token does not exist");
        return _ticketData[tokenId];
    }

    function totalSupply() external view returns (uint256) {
        return _nextTokenId;
    }
}
