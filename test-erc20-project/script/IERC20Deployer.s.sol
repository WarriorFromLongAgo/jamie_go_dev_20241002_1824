// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

import {Script, console} from "forge-std/Script.sol";

contract TestUsdtERC20 is ERC20 {
    constructor(string memory name, string memory symbol, uint256 initialSupply, address owner) ERC20(name, symbol) {
        _mint(owner, initialSupply);
    }
    
    function decimals() public view virtual override returns (uint8) {
        return 6;
    }
}

contract IERC20Deployer is Script {

    function run() public {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address deployerAddress = vm.addr(deployerPrivateKey);

        vm.startBroadcast(deployerPrivateKey);

        IERC20 test_usdt_token = new TestUsdtERC20("Test_USDT", "Test_USDT", 1 * 1e4 * 1e6, deployerAddress);
        console.log("deploy test_usdt_token:", address(test_usdt_token));
        console.log("deploy deployerAddress balance:", test_usdt_token.balanceOf(deployerAddress));

        vm.stopBroadcast();
    }
}
