// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "@openzeppelin/contracts/utils/Counters.sol";
import "./NFT.sol";
import "./Token.sol";

contract Controller {
    using Counters for Counters.Counter;

    //
    // STATE VARIABLES
    //
    Counters.Counter private _sessionIdCounter;
    GeneNFT public geneNFT;
    PostCovidStrokePrevention public pcspToken;

    struct UploadSession {
        uint256 id;
        address user;
        string proof;
        bool confirmed;
    }

    struct DataDoc {
        string id;
        string hashContent;
    }

    mapping(uint256 => UploadSession) sessions;
    mapping(string => DataDoc) docs;
    mapping(string => bool) docSubmits;
    mapping(uint256 => string) nftDocs;


    //
    // EVENTS
    //
    event UploadData(string docId, uint256 sessionId);
    event GeneNFTMinted(uint256 tokenId, string docId);
    event PCSPRewarded(address user, uint256 amount);

    constructor(address nftAddress, address pcspAddress) {
        geneNFT = GeneNFT(nftAddress);
        pcspToken = PostCovidStrokePrevention(pcspAddress);
    }


    modifier docNotSubmited(string memory docId) {
        require(!docSubmits[docId], "Doc already been submitted");
        _;
    }

    function uploadData(string memory docId) public docNotSubmited(docId) returns (uint256) {
        // to start an uploading gene data session. The doc id is used to identify a unique gene profile. Also should check if the doc id has been submited to the system before. This method return the session id

        // get current session id, and update current session data 
        uint256 sessionId = _sessionIdCounter.current();
        sessions[sessionId] = UploadSession({
            id: sessionId,
            user: msg.sender,
            proof: "",
            confirmed: false
        });
        
        // update doc submited flag
        docSubmits[docId] = true;

        // increment session id counter
        _sessionIdCounter.increment();

        // emit event
        emit UploadData(docId, sessionId);

        return sessionId;
    }

    function confirm(
        string memory docId,
        string memory contentHash,
        string memory proof,
        uint256 sessionId,
        uint256 riskScore
    ) public {
        // The proof here is used to verify that the result is returned from a valid computation on the gene data. For simplicity, we will skip the proof verification in this implementation. The gene data's owner will receive a NFT as a ownership certicate for his/her gene profile.
        require(bytes(docs[docId].id).length == 0, "Doc already been submitted");

        require(getSession(sessionId).user == msg.sender, "Invalid session owner");
        require(!getSession(sessionId).confirmed, "Session is ended");


        // verify proof
        require(_verifyProof(proof), "Invalid proof");
        
       
        // update doc content
        docs[docId] = DataDoc({
            id: docId,
            hashContent: contentHash
        });

        // Mint NFT 
        uint256 tokenId = geneNFT.safeMint(msg.sender);
        nftDocs[tokenId] = docId;

        // Reward PCSP token based on risk stroke
        uint256 rewardAmount = pcspToken.reward(msg.sender, riskScore);

        // Close session
        sessions[sessionId].confirmed = true;
        sessions[sessionId].proof = proof;

        // emit events
        emit GeneNFTMinted(tokenId, docId);
        emit PCSPRewarded(msg.sender, rewardAmount);
    }

    function getSession(uint256 sessionId) public view returns(UploadSession memory) {
        return sessions[sessionId];
    }

    function getDoc(string memory docId) public view returns(DataDoc memory) {
        return docs[docId];
    }


    /// @custom:dev-note Proof verification is skipped for development
    function _verifyProof(string memory /* _proof */) internal pure returns (bool) {
        return true;
    }
}
