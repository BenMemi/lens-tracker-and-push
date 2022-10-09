package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	//eth "github.com/ethereum/go-ethereum"

	"main/database"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	dotenv "github.com/profclems/go-dotenv" //Import dotenv library to deal with env variables before CICD is needed
	"github.com/shopspring/decimal"

	//Import GORM (go ORM) to interact with the database
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	ABI = `[{"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"name","type":"string"},{"indexed":false,"internalType":"string","name":"symbol","type":"string"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"BaseInitialized","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"collectModule","type":"address"},{"indexed":true,"internalType":"bool","name":"whitelisted","type":"bool"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"CollectModuleWhitelisted","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"profileId","type":"uint256"},{"indexed":true,"internalType":"uint256","name":"pubId","type":"uint256"},{"indexed":true,"internalType":"address","name":"collectNFT","type":"address"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"CollectNFTDeployed","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"profileId","type":"uint256"},{"indexed":true,"internalType":"uint256","name":"pubId","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"CollectNFTInitialized","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"profileId","type":"uint256"},{"indexed":true,"internalType":"uint256","name":"pubId","type":"uint256"},{"indexed":true,"internalType":"uint256","name":"collectNFTId","type":"uint256"},{"indexed":false,"internalType":"address","name":"from","type":"address"},{"indexed":false,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"CollectNFTTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"collector","type":"address"},{"indexed":true,"internalType":"uint256","name":"profileId","type":"uint256"},{"indexed":true,"internalType":"uint256","name":"pubId","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"rootProfileId","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"rootPubId","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"collectModuleData","type":"bytes"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"Collected","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"profileId","type":"uint256"},{"indexed":true,"internalType":"uint256","name":"pubId","type":"uint256"},{"indexed":false,"internalType":"string","name":"contentURI","type":"string"},{"indexed":false,"internalType":"uint256","name":"profileIdPointed","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"pubIdPointed","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"referenceModuleData","type":"bytes"},{"indexed":false,"internalType":"address","name":"collectModule","type":"address"},{"indexed":false,"internalType":"bytes","name":"collectModuleReturnData","type":"bytes"},{"indexed":false,"internalType":"address","name":"referenceModule","type":"address"},{"indexed":false,"internalType":"bytes","name":"referenceModuleReturnData","type":"bytes"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"CommentCreated","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"wallet","type":"address"},{"indexed":true,"internalType":"uint256","name":"profileId","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"DefaultProfileSet","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"profileId","type":"uint256"},{"indexed":true,"internalType":"address","name":"dispatcher","type":"address"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"DispatcherSet","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"caller","type":"address"},{"indexed":true,"internalType":"address","name":"oldEmergencyAdmin","type":"address"},{"indexed":true,"internalType":"address","name":"newEmergencyAdmin","type":"address"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"EmergencyAdminSet","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"moduleGlobals","type":"address"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"FeeModuleBaseConstructed","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"profileId","type":"uint256"},{"indexed":false,"internalType":"address","name":"followModule","type":"address"},{"indexed":false,"internalType":"bytes","name":"followModuleReturnData","type":"bytes"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"FollowModuleSet","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"followModule","type":"address"},{"indexed":true,"internalType":"bool","name":"whitelisted","type":"bool"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"FollowModuleWhitelisted","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"delegate","type":"address"},{"indexed":true,"internalType":"uint256","name":"newPower","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"FollowNFTDelegatedPowerChanged","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"profileId","type":"uint256"},{"indexed":true,"internalType":"address","name":"followNFT","type":"address"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"FollowNFTDeployed","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"profileId","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"FollowNFTInitialized","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"profileId","type":"uint256"},{"indexed":true,"internalType":"uint256","name":"followNFTId","type":"uint256"},{"indexed":false,"internalType":"address","name":"from","type":"address"},{"indexed":false,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"FollowNFTTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"profileId","type":"uint256"},{"indexed":false,"internalType":"string","name":"followNFTURI","type":"string"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"FollowNFTURISet","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"follower","type":"address"},{"indexed":false,"internalType":"uint256[]","name":"profileIds","type":"uint256[]"},{"indexed":false,"internalType":"bytes[]","name":"followModuleDatas","type":"bytes[]"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"Followed","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"uint256","name":"profileId","type":"uint256"},{"indexed":false,"internalType":"address[]","name":"addresses","type":"address[]"},{"indexed":false,"internalType":"bool[]","name":"approved","type":"bool[]"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"FollowsApproved","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"uint256[]","name":"profileIds","type":"uint256[]"},{"indexed":false,"internalType":"bool[]","name":"enabled","type":"bool[]"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"FollowsToggled","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"caller","type":"address"},{"indexed":true,"internalType":"address","name":"prevGovernance","type":"address"},{"indexed":true,"internalType":"address","name":"newGovernance","type":"address"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"GovernanceSet","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"profileId","type":"uint256"},{"indexed":true,"internalType":"uint256","name":"pubId","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"profileIdPointed","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"pubIdPointed","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"referenceModuleData","type":"bytes"},{"indexed":false,"internalType":"address","name":"referenceModule","type":"address"},{"indexed":false,"internalType":"bytes","name":"referenceModuleReturnData","type":"bytes"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"MirrorCreated","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"hub","type":"address"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"ModuleBaseConstructed","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"currency","type":"address"},{"indexed":true,"internalType":"bool","name":"prevWhitelisted","type":"bool"},{"indexed":true,"internalType":"bool","name":"whitelisted","type":"bool"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"ModuleGlobalsCurrencyWhitelisted","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"prevGovernance","type":"address"},{"indexed":true,"internalType":"address","name":"newGovernance","type":"address"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"ModuleGlobalsGovernanceSet","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint16","name":"prevTreasuryFee","type":"uint16"},{"indexed":true,"internalType":"uint16","name":"newTreasuryFee","type":"uint16"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"ModuleGlobalsTreasuryFeeSet","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"prevTreasury","type":"address"},{"indexed":true,"internalType":"address","name":"newTreasury","type":"address"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"ModuleGlobalsTreasurySet","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"profileId","type":"uint256"},{"indexed":true,"internalType":"uint256","name":"pubId","type":"uint256"},{"indexed":false,"internalType":"string","name":"contentURI","type":"string"},{"indexed":false,"internalType":"address","name":"collectModule","type":"address"},{"indexed":false,"internalType":"bytes","name":"collectModuleReturnData","type":"bytes"},{"indexed":false,"internalType":"address","name":"referenceModule","type":"address"},{"indexed":false,"internalType":"bytes","name":"referenceModuleReturnData","type":"bytes"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"PostCreated","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"profileId","type":"uint256"},{"indexed":true,"internalType":"address","name":"creator","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"string","name":"handle","type":"string"},{"indexed":false,"internalType":"string","name":"imageURI","type":"string"},{"indexed":false,"internalType":"address","name":"followModule","type":"address"},{"indexed":false,"internalType":"bytes","name":"followModuleReturnData","type":"bytes"},{"indexed":false,"internalType":"string","name":"followNFTURI","type":"string"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"ProfileCreated","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"profileCreator","type":"address"},{"indexed":true,"internalType":"bool","name":"whitelisted","type":"bool"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"ProfileCreatorWhitelisted","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"profileId","type":"uint256"},{"indexed":false,"internalType":"string","name":"imageURI","type":"string"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"ProfileImageURISet","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"referenceModule","type":"address"},{"indexed":true,"internalType":"bool","name":"whitelisted","type":"bool"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"ReferenceModuleWhitelisted","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"caller","type":"address"},{"indexed":true,"internalType":"enum DataTypes.ProtocolState","name":"prevState","type":"uint8"},{"indexed":true,"internalType":"enum DataTypes.ProtocolState","name":"newState","type":"uint8"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"StateSet","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"number","type":"uint256"}],"name":"bigint","type":"event"}]`
)

type CommentCreated struct {
	profileId                 decimal.Decimal
	pubId                     decimal.Decimal
	contentURI                string
	profileIdPointed          decimal.Decimal
	pubIdPointed              decimal.Decimal
	collectModule             string
	collectModuleReturnData   string
	referenceModule           string
	referenceModuleReturnData string
	timestamp                 decimal.Decimal
}

type FollowNFTDeployed struct {
	profileId decimal.Decimal
	followNFT string
	timestamp decimal.Decimal
}

func main() {

	dsn := ""
	err := dotenv.LoadConfig()
	if err != nil {
		//panic if we cannot load the .env
		fmt.Println("error loading .env file")
	}

	dsn = dotenv.GetString("DATABASE_URL")
	rpc := dotenv.GetString("RPC_URL")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//panic if we cannot connect to the database
	if err != nil {
		panic("failed to connect database")
	} else {
		//or else we are good to go
		fmt.Println("Connected to database")
		fmt.Println(db)
	}
	db.AutoMigrate(&database.User{})
	db.AutoMigrate(&database.FollowMessage{})
	db.AutoMigrate(&database.CommentMessage{})

	client, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal("Failed to connect to the websocket of the Node (RPC) ", err)
	} else {
		fmt.Println("successfully connected to the RPC endpoint!")
	}

	contractABI, err := abi.JSON(strings.NewReader(ABI))
	if err != nil {
		log.Fatal("could not convert JSON ABI string to ABI object")
	}

	contractAddress := common.HexToAddress("0x60Ae865ee4C725cd04353b5AAb364553f56ceF82")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{common.HexToHash("0x44403e38baed5e40df7f64ff8708b076c75a0dfda8380e75df5c36f11a476743")}},
	}

	query2 := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{common.HexToHash("0x7b4d1aa33773161799847429e4fbf29f56dbf1a3fe815f5070231cbfba402c37")}},
	}

	logs1 := make(chan types.Log)
	logs2 := make(chan types.Log)

	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs1)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("successfully subscribed to the contract events!")
	}

	sub2, err := client.SubscribeFilterLogs(context.Background(), query2, logs2)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("successfully subscribed to the contract events!")
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs1:
			fmt.Println("The Topic 0 of this event is;: ", vLog.Topics[0].Hex())
			Topics := vLog.Topics
			ProfileIDData, err := hexutil.Decode(Topics[1].Hex())
			if err != nil {
				log.Fatal(err)
			}

			ProfileIdInterface, err := contractABI.Unpack("bigint", ProfileIDData)

			if err != nil {
				log.Fatal(err)
			}

			ProfileIDBI := ProfileIdInterface[0].(*big.Int)
			profileID := decimal.NewFromBigInt(ProfileIDBI, 0)
			followNFT := common.HexToAddress((Topics[2].Hex())).Hex()

			Timestamp, err := contractABI.Unpack("bigint", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}
			TimestampBI := Timestamp[0].(*big.Int)
			TimestampDecimal := decimal.NewFromBigInt(TimestampBI, 0)

			myvar := database.FollowMessage{
				MessageID: uuid.New(),
				Sent:      false,
				ProfileId: profileID,
				FollowNFT: followNFT,
				Timestamp: TimestampDecimal,
			}

			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&myvar)

			fmt.Println(common.HexToAddress((Topics[2].Hex())).Hex())
		case err := <-sub2.Err():
			log.Fatal(err)
		case vLog2 := <-logs2:
			fmt.Println("The Topic 0 of this event is;: ", vLog2.Topics[0].Hex())
			//Topics2 := vLog.Topics
			inrerfaces, err := contractABI.Unpack("CommentCreated", vLog2.Data)

			fmt.Println("I am up to here")
			if err != nil {
				log.Fatal(err)
			}

			profileIdData, err := hexutil.Decode(vLog2.Topics[1].Hex())
			if err != nil {
				log.Fatal(err)
			}
			profileIdInterface, err := contractABI.Unpack("bigint", profileIdData)
			if err != nil {
				log.Fatal(err)
			}
			profileId := profileIdInterface[0].(*big.Int)
			prodileIdDecimal := decimal.NewFromBigInt(profileId, 0)

			pubIdData, err := hexutil.Decode(vLog2.Topics[2].Hex())
			if err != nil {
				log.Fatal(err)
			}
			pubIdInterface, err := contractABI.Unpack("bigint", pubIdData)
			if err != nil {
				log.Fatal(err)
			}
			pubId := pubIdInterface[0].(*big.Int)
			pubIdDecimal := decimal.NewFromBigInt(pubId, 0)

			contentURI := inrerfaces[0].(string)
			profileIdPointed := inrerfaces[1].(*big.Int)
			pubIdPointed := inrerfaces[2].(*big.Int)

			myvar := database.CommentMessage{
				MessageID:        uuid.New(),
				Sent:             false,
				ProfileId:        prodileIdDecimal,
				PubId:            pubIdDecimal,
				ContentURI:       contentURI,
				ProfileIdPointed: decimal.NewFromBigInt(profileIdPointed, 0),
				PubIdPointed:     decimal.NewFromBigInt(pubIdPointed, 0),
			}

			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&myvar)

		}
	}
}

// Main
// wss://polygon-mainnet.g.alchemy.com/v2/3ZR9MXbyYN4nBj4IWZJX9XNqxVpYUK2M
// Test
// wss://polygon-mumbai.g.alchemy.com/v2/-xLct1D6mffFUeh-NhHKvIQ1Qe6sNBqe