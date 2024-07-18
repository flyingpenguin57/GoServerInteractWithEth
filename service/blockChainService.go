package blockChain

import (
	"bubble/util"
	"context"
	"fmt"
	"log"
	"strings"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
)

var ethClient *ethclient.Client
var blockList = queue.NewQueue[types.Header](10)
var TransferInfoList = queue.NewQueue[types.Header](100)

func InitEthClient() {
	// 连接到以太坊节点
	_ethClient, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/1d37b9d398af4b81baca54ea5f164f17")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	ethClient = _ethClient
	if ethClient == nil {
		log.Println("client init fail")
	} else {
		log.Println("client init success")
	}
}

func GetLatestBlock(c *gin.Context) *queue.Queue[types.Header] {
	return blockList
}

func QueryLatestBlockFromChain() {

	log.Println("start get data from chain")

	header, err := ethClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to retrieve latest block header: %v", err)
	}

	if blockList.IsEmpty() {
		blockList.Enqueue(*header)
	} else {
		if header.Number.Cmp(blockList.GetHeader().Number) == 1 {
			blockList.Enqueue(*header)
		}
	}
	log.Println("end get data from chain")
}

func QueryTransferInfoFromBlockChain() {

	// 将字节切片转换为字符串
	contractABI := abiString

	contractAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")

	// 解析ABI
	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	// 创建查询过滤器
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	// 监听事件
	logs := make(chan types.Log)
	sub, err := ethClient.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("Failed to subscribe to logs: %v", err)
	}
	fmt.Println("Listening for Transfer events...")

	// 处理事件
	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("Error: %v", err)
		case vLog := <-logs:
			// 解析事件
			event := new(TransferEvent)
			err := parsedABI.UnpackIntoInterface(event, "Transfer", vLog.Data)
			if err != nil {
				log.Fatalf("Failed to unpack log: %v", err)
			}

			// 获取 indexed 参数
			event.From = common.HexToAddress(vLog.Topics[1].Hex())
			event.To = common.HexToAddress(vLog.Topics[2].Hex())

			fmt.Printf("Transfer event detected: from %s to %s value %s\n", event.From.Hex(), event.To.Hex(), event.Value.String())
		}
	}

}
