package scheduled

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	global_const "go-project/common"
	"go.uber.org/zap"
	"math/big"
	"time"

	"gorm.io/gorm"

	"go-project/chain/eth"
	"go-project/main/log"
)

type TestIncrementBlock struct {
	ctx         context.Context
	ethClient   eth.EthClient
	erc20Client eth.TestErc20Client
	db          *gorm.DB
	log         *log.ZapLogger
}

func NewTestIncrementBlock(ctx context.Context, client eth.EthClient, erc20Client eth.TestErc20Client, db *gorm.DB, log *log.ZapLogger) (*TestIncrementBlock, error) {
	return &TestIncrementBlock{
		ctx:         ctx,
		ethClient:   client,
		erc20Client: erc20Client,
		db:          db,
		log:         log,
	}, nil
}

func (s *TestIncrementBlock) Start() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			fmt.Println("incrementBlock done")
			return
		case <-ticker.C:
			err := s.incrementBlock()
			if err != nil {
				fmt.Printf("incrementBlock incrementBlock error: %v\n", err)
			}
			err = s.transferERC20()
			if err != nil {
				fmt.Printf("incrementBlock transferERC20 error: %v\n", err)
			}
		}
	}
}

func (s *TestIncrementBlock) incrementBlock() error {
	fromAddress := common.HexToAddress("0x976EA74026E726554dB657fA54763abd0C3a0aa9")
	toAddress := common.HexToAddress("0x14dC79964da2C08b23698B3D3cc7Ca32193d9955")
	privateKey, err := crypto.HexToECDSA("92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e")
	if err != nil {
		return fmt.Errorf("解析私钥失败: %w", err)
	}

	s.log.Info("incrementBlock", zap.Any("fromAddress", fromAddress), zap.Any("toAddress", toAddress))

	nonce, err := s.ethClient.TxCountByAddress(fromAddress)
	if err != nil {
		return fmt.Errorf("获取nonce失败: %w", err)
	}

	gasPrice, err := s.ethClient.SuggestGasPrice()
	if err != nil {
		return fmt.Errorf("获取gas价格失败: %w", err)
	}

	transferAmount := new(big.Int).Mul(big.NewInt(10000000000000000), big.NewInt(1))

	tx := types.NewTransaction(uint64(nonce), toAddress, transferAmount, 21000, gasPrice, nil)

	chainId := new(big.Int).SetUint64(global_const.ChainId)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
	if err != nil {
		return fmt.Errorf("签名交易失败: %w", err)
	}
	rawTxBytes, err := signedTx.MarshalBinary()
	if err != nil {
		return err
	}
	rawTxHex := hexutil.Encode(rawTxBytes)
	err = s.ethClient.SendRawTransaction(rawTxHex)
	if err != nil {
		return fmt.Errorf("发送交易失败: %w", err)
	}

	s.log.Info("转账成功", zap.String("From", fromAddress.Hex()), zap.String("To", toAddress.Hex()), zap.String("Amount", "0.01 ETH"), zap.String("TxHash", signedTx.Hash().Hex()))

	err = s.printBalances(fromAddress, toAddress)
	if err != nil {
		s.log.Error("打印余额失败", zap.Error(err))
	}

	return nil
}

func (s *TestIncrementBlock) printBalances(fromAddress, toAddress common.Address) error {
	fromBalance, err := s.ethClient.BalanceAt(s.ctx, fromAddress, nil)
	if err != nil {
		return fmt.Errorf("获取发送方余额失败: %w", err)
	}

	toBalance, err := s.ethClient.BalanceAt(s.ctx, toAddress, nil)
	if err != nil {
		return fmt.Errorf("获取接收方余额失败: %w", err)
	}

	s.log.Info("当前余额",
		zap.String("FromAddress", fromAddress.Hex()),
		zap.String("FromBalance", fromBalance.String()),
		zap.String("ToAddress", toAddress.Hex()),
		zap.String("ToBalance", toBalance.String()))

	return nil
}

// ... 现有代码 ...

func (s *TestIncrementBlock) transferERC20() error {
	privateKey, err := crypto.HexToECDSA(global_const.OWNER_PRV_KEY)
	if err != nil {
		return fmt.Errorf("解析私钥失败: %w", err)
	}
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	toAddress := common.HexToAddress("0x14dC79964da2C08b23698B3D3cc7Ca32193d9955")
	tokenAddress := common.HexToAddress(global_const.TEMP_TEST_ERC20_ADDRESS)

	amount := big.NewInt(1 * 1e6)

	ethBusiness := eth.NewEthBusinessService(s.ethClient, s.erc20Client, s.log)
	txHash, transferData, err := ethBusiness.TransferERC20(context.Background(), privateKey, fromAddress.Hex(), toAddress.Hex(), tokenAddress.Hex(), amount)
	if err != nil {
		return err
	}

	fmt.Printf("scanBlocks BlockByNumberV3 transactionsJson %s \n", txHash)
	fmt.Printf("scanBlocks BlockByNumberV3 transactionsJson %s \n", hexutil.Encode(transferData))

	printERC20Balance(s.ctx, s.erc20Client, fromAddress, "From")
	printERC20Balance(s.ctx, s.erc20Client, toAddress, "To")

	return nil
}
