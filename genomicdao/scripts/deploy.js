const hre = require("hardhat");

async function main() {
  const [deployer] = await ethers.getSigners();
  console.log("Deploying contracts with account:", deployer.address);

  // Deploy GeneNFT
  const GeneNFT = await ethers.getContractFactory("GeneNFT");
  const nft = await GeneNFT.deploy();
  await nft.waitForDeployment();
  console.log("GeneNFT deployed to:", nft.target);

  // Deploy PCSP Token
  const PCSPToken = await ethers.getContractFactory("PostCovidStrokePrevention");
  const token = await PCSPToken.deploy();
  await token.waitForDeployment();
  console.log("PCSP Token deployed to:", token.target);

  // Deploy Controller
  const Controller = await ethers.getContractFactory("Controller");
  const controller = await Controller.deploy(nft.target, token.target);
  await controller.waitForDeployment();
  console.log("Controller deployed to:", controller.target);

  // Transfer ownership to Controller
  await nft.transferOwnership(controller.target);
  await token.transferOwnership(controller.target);
  console.log("Ownership transferred to Controller");

  // Log all addresses for future reference
  console.log("\nContract Addresses:");
  console.log("-------------------");
  console.log("GeneNFT:", nft.target);
  console.log("PCSP Token:", token.target);
  console.log("Controller:", controller.target);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});