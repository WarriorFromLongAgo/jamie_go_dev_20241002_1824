package eth

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	globalconst "go-project/common"
)

func TestEthClient_LatestFinalizedBlockHeader(t *testing.T) {
	ctx := context.Background()
	ethClient, err := DialEthClient(ctx, "http://127.0.0.1:8545")
	if err != nil {
		t.Fatalf("Failed to create EthClient: %v", err)
	}
	defer ethClient.Close()

	header, err := ethClient.LatestFinalizedBlockHeader()
	if err != nil {
		t.Fatalf("Failed to get latest finalized block header: %v", err)
	}
	t.Logf("Latest finalized block number: %d", header.Number.Uint64())
}

func TestEthClient_BlockHeaderByBlockHash(t *testing.T) {
	ctx := context.Background()
	ethClient, err := DialEthClient(ctx, "http://127.0.0.1:8545")
	if err != nil {
		t.Fatalf("Failed to create EthClient: %v", err)
	}
	defer ethClient.Close()

	latestHeader, err := ethClient.LatestFinalizedBlockHeader()
	if err != nil {
		t.Fatalf("Failed to get latest finalized block header: %v", err)
	}

	header, err := ethClient.BlockHeaderByBlockHash(latestHeader.Hash())
	if err != nil {
		t.Fatalf("Failed to get block header by hash: %v", err)
	}
	t.Logf("Block header number: %d", header.Number.Uint64())
}

func TestEthClient_BlockHeaderListByRange(t *testing.T) {
	ctx := context.Background()
	ethClient, err := DialEthClient(ctx, "http://127.0.0.1:8545")
	if err != nil {
		t.Fatalf("Failed to create EthClient: %v", err)
	}
	defer ethClient.Close()

	latestHeader, err := ethClient.LatestFinalizedBlockHeader()
	if err != nil {
		t.Fatalf("Failed to get latest finalized block header: %v", err)
	}

	startBlock := new(big.Int).Sub(latestHeader.Number, big.NewInt(10))
	headers, err := ethClient.BlockHeaderListByRange(startBlock, latestHeader.Number)
	if err != nil {
		t.Fatalf("Failed to get block header list: %v", err)
	}
	t.Logf("Retrieved %d block headers", len(headers))
}

func TestEthClient_BlockByNumber(t *testing.T) {
	ctx := context.Background()

	t.Log("正在连接 EthClient...")
	ethClient, err := DialEthClient(ctx, "http://127.0.0.1:8545")
	if err != nil {
		t.Fatalf("创建 EthClient 失败: %v", err)
	}
	if ethClient == nil {
		t.Fatalf("EthClient 为 nil，尽管没有报错")
	}
	defer ethClient.Close()
	t.Log("EthClient 创建成功")

	// 选择一个区块号
	blockHeader, _ := ethClient.LatestFinalizedBlockHeader()
	blockNumber := blockHeader.Number
	//blockNumber := big.NewInt(740) // 替换为您想测试的实际区块号
	// 获取区块
	block, err := ethClient.BlockByNumberV3(ctx, blockNumber)
	if err != nil {
		t.Fatalf("获取区块失败: %v", err)
	}
	blockJson, _ := json.Marshal(block)
	fmt.Printf("scanBlocks BlockByNumberV3 block number %d \n", blockNumber)
	fmt.Printf("scanBlocks BlockByNumberV3 block number %d \n", block.Number())
	fmt.Printf("scanBlocks BlockByNumberV3 block %s \n", blockJson)

	// 获取区块中的交易
	transactions := block.Transactions()
	transactionsJson, _ := json.Marshal(transactions)
	fmt.Printf("scanBlocks BlockByNumberV3 transactionsJson %s \n", transactionsJson)

	// 打印交易数量
	t.Logf("区块 %d 中包含 %d 笔交易", blockNumber, len(transactions))

	// 遍历并打印每笔交易的哈希
	for i, tx := range transactions {
		t.Logf("交易 %d: %s", i, tx.Hash().Hex())
	}
}

func TestEthClient_TxByTxHash(t *testing.T) {
	ctx := context.Background()
	ethClient, err := DialEthClient(ctx, "http://127.0.0.1:8545")
	if err != nil {
		t.Fatalf("Failed to create EthClient: %v", err)
	}
	defer ethClient.Close()

	// 这里需要一个有效的交易哈希,您可能需要先发送一个交易或者从区块链上获取一个有效的交易哈希
	txHash := common.HexToHash("0xd1085a5feae0dd6ced58e7facf75cce3cfd1f07c6138e38c071bc92ee3e50ea5")
	tx, err := ethClient.TxByTxHash(txHash)
	if err != nil {
		t.Fatalf("Failed to get transaction by hash: %v", err)
	}
	marshal, _ := json.Marshal(tx)
	t.Logf("Transaction: %s", string(marshal))
}

func TestEthClient_TxReceiptByTxHash(t *testing.T) {
	ctx := context.Background()
	ethClient, err := DialEthClient(ctx, "http://127.0.0.1:8545")
	if err != nil {
		t.Fatalf("Failed to create EthClient: %v", err)
	}
	defer ethClient.Close()

	// 这里需要一个有效的交易哈希,您可能需要先发送一个交易或者从区块链上获取一个有效的交易哈希
	txHash := common.HexToHash("0x1234567890123456789012345678901234567890123456789012345678901234")
	receipt, err := ethClient.TxReceiptByTxHash(txHash)
	if err != nil {
		t.Fatalf("Failed to get transaction receipt: %v", err)
	}
	t.Logf("Transaction status: %d", receipt.Status)
}

//func TestEthClient_TxCountByAddress(t *testing.T) {
//	ctx := context.Background()
//	ethClient, err := DialEthClient(ctx, "http://127.0.0.1:8545")
//	if err != nil {
//		t.Fatalf("Failed to create EthClient: %v", err)
//	}
//	defer ethClient.Close()
//
//	address := common.HexToAddress("0xa0Ee7A142d267C1f36714E4a8F75612F20a79720")
//	count, err := ethClient.TxCountByAddress(address)
//	if err != nil {
//		t.Fatalf("Failed to get transaction count: %v", err)
//	}
//	t.Logf("Transaction count for address: %d", count)
//}
//
//func TestEthClient_SuggestGasPrice(t *testing.T) {
//	ctx := context.Background()
//	ethClient, err := DialEthClient(ctx, "http://127.0.0.1:8545")
//	if err != nil {
//		t.Fatalf("Failed to create EthClient: %v", err)
//	}
//	defer ethClient.Close()
//
//	gasPrice, err := ethClient.SuggestGasPrice()
//	if err != nil {
//		t.Fatalf("Failed to get suggested gas price: %v", err)
//	}
//	t.Logf("Suggested gas price: %s", gasPrice.String())
//}
//
//func TestEthClient_SuggestGasTipCap(t *testing.T) {
//	ctx := context.Background()
//	ethClient, err := DialEthClient(ctx, "http://127.0.0.1:8545")
//	if err != nil {
//		t.Fatalf("Failed to create EthClient: %v", err)
//	}
//	defer ethClient.Close()
//
//	gasTipCap, err := ethClient.SuggestGasTipCap()
//	if err != nil {
//		t.Fatalf("Failed to get suggested gas tip cap: %v", err)
//	}
//	t.Logf("Suggested gas tip cap: %s", gasTipCap.String())
//}

func TestEthClient_SendRawTransaction(t *testing.T) {
	ctx := context.Background()
	ethClient, err := DialEthClient(ctx, "http://127.0.0.1:8545")
	if err != nil {
		t.Fatalf("Failed to create EthClient: %v", err)
	}
	defer ethClient.Close()

	privateKey, err := crypto.HexToECDSA(globalconst.OWNER_PRV_KEY)
	if err != nil {
		t.Fatalf("Failed to create private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Fatal("Failed to get public key")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := ethClient.TxCountByAddress(fromAddress)
	if err != nil {
		t.Fatalf("Failed to get nonce: %v", err)
	}

	value := big.NewInt(1 * 1e18)
	gasLimit := uint64(21000)
	gasPrice, err := ethClient.SuggestGasPrice()
	if err != nil {
		t.Fatalf("Failed to get gas price: %v", err)
	}
	toAddress := common.HexToAddress(globalconst.TEMP_TO_ADDRESS)

	tx := types.NewTransaction(uint64(nonce), toAddress, value, gasLimit, gasPrice, nil)

	chainID := big.NewInt(globalconst.ChainId)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		t.Fatalf("Failed to sign transaction: %v", err)
	}

	rawTxBytes, err := signedTx.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to serialize transaction: %v", err)
	}
	rawTxHex := hexutil.Encode(rawTxBytes)

	printEthBalance(t, ctx, ethClient, fromAddress, "From (before)")
	printEthBalance(t, ctx, ethClient, toAddress, "To (before)")

	err = ethClient.SendRawTransaction(rawTxHex)
	if err != nil {
		t.Fatalf("Failed to send raw transaction: %v", err)
	}
	t.Logf("Raw transaction sent successfully: %s", signedTx.Hash().Hex())

	time.Sleep(5 * time.Second)

	printEthBalance(t, ctx, ethClient, fromAddress, "From (after)")
	printEthBalance(t, ctx, ethClient, toAddress, "To (after)")
}

func printEthBalance(t *testing.T, ctx context.Context, client EthClient, address common.Address, label string) {
	balance, err := client.BalanceAt(ctx, address, nil)
	if err != nil {
		t.Fatalf("Failed to get balance for %s address: %v", label, err)
	}
	t.Logf("%s address balance: %s ETH", label, new(big.Float).Quo(new(big.Float).SetInt(balance), new(big.Float).SetInt64(1e18)).Text('f', 6))
}

func TestEthClient_SendRawERC20Transaction(t *testing.T) {
	ctx := context.Background()
	ethClient, err := DialEthClient(ctx, "http://127.0.0.1:8545")
	if err != nil {
		t.Fatalf("Failed to create EthClient: %v", err)
	}
	defer ethClient.Close()

	erc20Client, err := NewTestErc20Client(ctx, "http://127.0.0.1:8545", globalconst.TEMP_TEST_ERC20_ADDRESS)
	if err != nil {
		t.Fatalf("Failed to create TestErc20Client: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(globalconst.OWNER_PRV_KEY)
	if err != nil {
		t.Fatalf("Failed to create private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Fatal("Failed to get public key")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := ethClient.TxCountByAddress(fromAddress)
	if err != nil {
		t.Fatalf("Failed to get nonce: %v", err)
	}

	toAddress := common.HexToAddress(globalconst.TEMP_TO_ADDRESS)
	amount := big.NewInt(9 * 1e6)

	gasPrice, err := ethClient.SuggestGasPrice()
	if err != nil {
		t.Fatalf("Failed to get gas price: %v", err)
	}

	printERC20Balance(t, ctx, erc20Client, fromAddress, "From (before)")
	printERC20Balance(t, ctx, erc20Client, toAddress, "To (before)")

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := crypto.Keccak256(transferFnSignature)
	methodID := hash[:4]
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	erc20Address := common.HexToAddress(globalconst.TEMP_TEST_ERC20_ADDRESS)
	tx := types.NewTransaction(uint64(nonce), erc20Address, big.NewInt(0), 300000, gasPrice, data)

	chainID := big.NewInt(globalconst.ChainId)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		t.Fatalf("Failed to sign transaction: %v", err)
	}

	rawTxBytes, err := signedTx.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to serialize transaction: %v", err)
	}
	rawTxHex := hexutil.Encode(rawTxBytes)

	err = ethClient.SendRawTransaction(rawTxHex)
	if err != nil {
		t.Fatalf("Failed to send raw transaction: %v", err)
	}
	t.Logf("ERC20 transfer transaction sent successfully: %s", signedTx.Hash().Hex())

	err = WaitForTransaction(ctx, ethClient, signedTx.Hash())
	if err != nil {
		t.Fatalf("Failed to wait for transaction: %v", err)
	}

	printERC20Balance(t, ctx, erc20Client, fromAddress, "From (after)")
	printERC20Balance(t, ctx, erc20Client, toAddress, "To (after)")
}
