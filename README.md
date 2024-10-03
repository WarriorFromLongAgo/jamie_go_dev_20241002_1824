# jamie_go_dev_20241002_1824

# go-project
## env
```
go version: 1.22.5

go get -u gorm.io/driver/mysql
go get -u gorm.io/gorm
go get -u github.com/ethereum/go-ethereum
go get -u github.com/gin-gonic/gin

```

## run
```
index = scan block

api = api server

```


# test-erc20-project
## env

```
openzeppelin: v5.0.2
forge 0.2.0 (fdfaafd 2024-07-30T00:24:49.704507500Z)

```

## run(use anvil)
```
cd test-erc20-project

$env:PRIVATE_KEY = "0x2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6"
$env:RPC_URL = "http://127.0.0.1:8545"
forge script ./script/IERC20Deployer.s.sol:IERC20Deployer --rpc-url $env:RPC_URL --private-key $env:PRIVATE_KEY --broadcast -vvvvv

forge build

abigen --abi=./test_erc20.abi --pkg=testerc20 --out=TestErc20Project.go

```

