package service

import (
	"context"
	"github.com/sasaxie/monitor/api"
	"github.com/sasaxie/monitor/common/base58"
	"github.com/sasaxie/monitor/common/hexutil"
	"github.com/sasaxie/monitor/core"
	"google.golang.org/grpc"
	"log"
)

type GrpcClient struct {
	Address        string
	Conn           *grpc.ClientConn
	WalletClient   api.WalletClient
	DatabaseClient api.DatabaseClient
}

func NewGrpcClient(address string) *GrpcClient {
	client := new(GrpcClient)
	client.Address = address
	return client
}

func (g *GrpcClient) Start() {
	var err error
	g.Conn, err = grpc.Dial(g.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}

	g.WalletClient = api.NewWalletClient(g.Conn)
	g.DatabaseClient = api.NewDatabaseClient(g.Conn)
}

func (g *GrpcClient) ListWitnesses() *api.WitnessList {
	witnessList, err := g.WalletClient.ListWitnesses(context.Background(),
		new(api.EmptyMessage))

	if err != nil {
		log.Fatalf("get witnesses error: %v\n", err)
	}

	return witnessList
}

func (g *GrpcClient) ListNodes() *api.NodeList {
	nodeList, err := g.WalletClient.ListNodes(context.Background(),
		new(api.EmptyMessage))

	if err != nil {
		log.Fatalf("get nodes error: %v\n", err)
	}

	return nodeList
}

func (g *GrpcClient) GetAccount(address string) *core.Account {
	account := new(core.Account)

	account.Address = base58.DecodeCheck(address)

	result, err := g.WalletClient.GetAccount(context.Background(), account)

	if err != nil {
		log.Fatalf("get account error: %v\n", err)
	}

	return result
}

func (g *GrpcClient) GetNowBlock() *core.Block {
	result, err := g.WalletClient.GetNowBlock(context.Background(), new(api.EmptyMessage))

	if err != nil {
		log.Printf("get now block error: %v\n", err)
		return nil
	}

	return result
}

func (g *GrpcClient) GetAssetIssueByAccount(address string) *api.AssetIssueList {
	account := new(core.Account)

	account.Address = base58.DecodeCheck(address)

	result, err := g.WalletClient.GetAssetIssueByAccount(context.Background(),
		account)

	if err != nil {
		log.Fatalf("get asset issue by account error: %v", err)
	}

	return result
}

func (g *GrpcClient) GetNextMaintenanceTime() *api.NumberMessage {

	result, err := g.WalletClient.GetNextMaintenanceTime(context.Background(),
		new(api.EmptyMessage))

	if err != nil {
		log.Fatalf("get next maintenance time error: %v", err)
	}

	return result
}

func (g *GrpcClient) TotalTransaction() *api.NumberMessage {

	result, err := g.WalletClient.TotalTransaction(context.Background(),
		new(api.EmptyMessage))

	if err != nil {
		log.Fatalf("total transaction error: %v", err)
	}

	return result
}

func (g *GrpcClient) GetAccountNet(address string) *api.AccountNetMessage {
	account := new(core.Account)

	account.Address = base58.DecodeCheck(address)

	result, err := g.WalletClient.GetAccountNet(context.Background(), account)

	if err != nil {
		log.Fatalf("get account net error: %v", err)
	}

	return result
}

func (g *GrpcClient) GetAssetIssueByName(name string) *core.AssetIssueContract {

	assetName := new(api.BytesMessage)
	assetName.Value = []byte(name)

	result, err := g.WalletClient.GetAssetIssueByName(context.Background(), assetName)

	if err != nil {
		log.Fatalf("get asset issue by name error: %v", err)
	}

	return result
}

func (g *GrpcClient) GetBlockByNum(num int64) *core.Block {
	numMessage := new(api.NumberMessage)
	numMessage.Num = num

	result, err := g.WalletClient.GetBlockByNum(context.Background(), numMessage)

	if err != nil {
		log.Fatalf("get block by num error: %v", err)
	}

	return result
}

func (g *GrpcClient) GetBlockById(id string) *core.Block {
	blockId := new(api.BytesMessage)
	var err error

	blockId.Value, err = hexutil.Decode(id)

	if err != nil {
		log.Fatalf("get block by id error: %v", err)
	}

	result, err := g.WalletClient.GetBlockById(context.Background(), blockId)

	if err != nil {
		log.Fatalf("get block by id error: %v", err)
	}

	return result
}

func (g *GrpcClient) GetAssetIssueList() *api.AssetIssueList {

	result, err := g.WalletClient.GetAssetIssueList(context.Background(), new(api.EmptyMessage))

	if err != nil {
		log.Fatalf("get asset issue list error: %v", err)
	}

	return result
}

func (g *GrpcClient) GetBlockByLimitNext(start, end int64) *api.BlockList {
	blockLimit := new(api.BlockLimit)
	blockLimit.StartNum = start
	blockLimit.EndNum = end

	result, err := g.WalletClient.GetBlockByLimitNext(context.Background(), blockLimit)

	if err != nil {
		log.Fatalf("get block by limit next error: %v", err)
	}

	return result
}

func (g *GrpcClient) GetTransactionById(id string) *core.Transaction {
	transactionId := new(api.BytesMessage)
	var err error

	transactionId.Value, err = hexutil.Decode(id)

	if err != nil {
		log.Fatalf("get transaction by id error: %v", err)
	}

	result, err := g.WalletClient.GetTransactionById(context.Background(), transactionId)

	if err != nil {
		log.Fatalf("get transaction by limit next error: %v", err)
	}

	return result
}

func (g *GrpcClient) GetBlockByLatestNum(num int64) *api.BlockList {
	numMessage := new(api.NumberMessage)
	numMessage.Num = num

	result, err := g.WalletClient.GetBlockByLatestNum(context.Background(), numMessage)

	if err != nil {
		log.Fatalf("get block by latest num error: %v", err)
	}

	return result
}

func (g *GrpcClient) GetLastSolidityBlockNum() int64 {
	dynamicProperties, err := g.DatabaseClient.GetDynamicProperties(context.
		Background(), new(api.EmptyMessage))

	if err != nil {
		log.Printf("get last solidity block num error: %v", err)
		return 0
	}

	return dynamicProperties.LastSolidityBlockNum
}
