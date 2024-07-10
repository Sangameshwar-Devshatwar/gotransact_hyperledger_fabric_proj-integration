/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fabric

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func Initfabric() *gateway.Contract {
	log.Println("============ application-golang starts ============")

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environment variable: %v", err)
	}

	walletPath := "wallet"
	// remove any existing wallet from prior runs
	os.RemoveAll(walletPath)
	wallet, err := gateway.NewFileSystemWallet(walletPath)
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := filepath.Join(
		// "..",
		// "..",
		// "..",
		"/home/trellis",
		"hyperledger",
		"fabric-samples",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	channelName := "mychannel"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	log.Println("--> Connecting to channel", channelName)
	network, err := gw.GetNetwork(channelName)
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	chaincodeName := "petmgmt"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	log.Println("--> Using chaincode", chaincodeName)
	contract := network.GetContract(chaincodeName)
	return contract
	// log.Println("--> Submit Transaction: InitLedger, function creates the initial set of owners & doctors on the ledger")
	// result, err := contract.SubmitTransaction("InitLedger")
	// if err != nil {
	// 	log.Fatalf("Failed to Submit transaction: %v", err)
	// }
	// log.Println(result)
	// fmt.Println(contract.)

	// log.Println("--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")
	// result, err := contract.EvaluateTransaction("GetAllAssets")
	// if err != nil {
	// 	log.Fatalf("Failed to evaluate transaction: %v", err)
	// }
	// log.Println(string(result))

	// log.Println("--> Submit Transaction: CreateAsset, creates new pet with ID, name, barcode, ownerid")
	// // {ID: "asset1", Color: "blue", Size: 5, Owner: "Tomoko", AppraisedValue: 300},
	// result, err = contract.SubmitTransaction("CreateAsset", "asset7", "Pink", "10", "sourabh", "700")
	// if err != nil {
	// 	log.Fatalf("Failed to Submit transaction: %v", err)
	// }
	// log.Println(string(result))

	// |
	// |
	// |

	// log.Println("--> Evaluate Transaction: QueryPet, function returns  pet with a given petID")
	// result, err = contract.EvaluateTransaction("QueryPet", "pet1")
	// if err != nil {
	// 	log.Fatalf("Failed to evaluate transaction: %v\n", err)
	// }
	// log.Println(string(result))

	// log.Println("--> Evaluate Transaction: AssetExists, function returns 'true' if an asset with given assetID exist")
	// result, err = contract.EvaluateTransaction("AssetExists", "asset1")
	// if err != nil {
	//  log.Fatalf("Failed to evaluate transaction: %v\n", err)
	// }
	// log.Println(string(result))

	// log.Println("--> Submit Transaction: TransferAsset asset1, transfer to new owner of Tom")
	// _, err = contract.SubmitTransaction("TransferAsset", "asset1", "Tom")
	// if err != nil {
	//  log.Fatalf("Failed to Submit transaction: %v", err)
	// }

	// log.Println("--> Evaluate Transaction: ReadAsset, function returns 'asset1' attributes")
	// result, err = contract.EvaluateTransaction("ReadAsset", "asset1")
	// if err != nil {
	//  log.Fatalf("Failed to evaluate transaction: %v", err)
	// }
	// log.Println(string(result))
	// |
	// |
	// |
	// log.Println("--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")
	// result, err = contract.EvaluateTransaction("GetAllAssets")
	// if err != nil {
	// 	log.Fatalf("Failed to evaluate transaction: %v", err)
	// }
	// log.Println(string(result))
	// log.Println("============ application-golang ends ============")
}

func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
	credPath := filepath.Join(
		// "..",
		// "..",
		// "..",
		"/home/trellis",
		"hyperledger",
		"fabric-samples",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := os.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := os.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := os.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}
