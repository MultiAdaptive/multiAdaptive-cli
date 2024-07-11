package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/MultiAdaptive/multiAdaptive-cli/bindings"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/kzg"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	kzgsdk "github.com/multiAdaptive/kzg-sdk"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
)

const (
	dataSize              = 5 * 1024 * 1024
	cmManagerAddress      = "0xa8ED91Eb2B65A681A742011798d7FB31C50FA724"
	nodeManagerAddress    = "0x97bE3172AEA87b60224e8d604aC4bAbe55F067EC"
	storageManagerAddress = "0x664250Fb3b1cd58f07683D957A34daf8A06130fe"
	chainID               = 11155111
	ethUrl                = "https://eth-sepolia.public.blastapi.io"
)

func main() {
	privateKeyHex := getEnv("PRIVATE_KEY", "")
	addressHex := getEnv("ADDRESS", "")
	validateEnvVars(privateKeyHex, addressHex)

	// Prompt user to select an action
	action := getAction()
	// Execute the selected action
	executeAction(action)
}

// Retrieve environment variable value or fallback if not set
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Validate that all environment variables are set
func validateEnvVars(vars ...string) {
	for _, v := range vars {
		if v == "" {
			log.Fatalf("Environment variable is not set")
		}
	}
}

// Display a prompt to select an action
func getAction() string {
	var action string
	actionPrompt := &survey.Select{
		Message: "What do you want to do?",
		Options: []string{
			"View Broadcast Node Information",
			"View Storage Node Information",
			"Register NodeGroup",
			"Register NameSpace",
			"General test",
			"Advanced test",
			"Register as a broadcast node",
			"Register as a storage node",
		},
		PageSize: 8,
	}
	err := survey.AskOne(actionPrompt, &action)
	if err != nil {
		log.Fatal(err.Error())
	}
	return action
}

// Execute the selected action based on the user input
func executeAction(action string) {
	switch action {
	case "View Broadcast Node Information":
		displayNodeInfo(viewBroadcastNodeInfo)
	case "View Storage Node Information":
		displayNodeInfo(viewStorageNodeInfo)
	case "Register NodeGroup":
		key := registerNodeGroup()
		log.Printf("nodeGroupKeys: %s", key.Hex())
	case "Register NameSpace":
		key := registerNameSpace()
		log.Printf("nameSpaceKey: %s", key.Hex())
	case "General test":
		generalTest()
	case "Advanced test":
		advancedTest()
	case "Register as a broadcast node":
		registerBroadcastNode()
	case "Register as a storage node":
		registerStorageNode()
	default:
		fmt.Println("Unknown action")
	}
}

// View information of broadcasting nodes
func viewBroadcastNodeInfo() ([]bindings.NodeManagerNodeInfo, error) {
	_, instance, err := getNodeManagerInstance()
	if err != nil {
		return nil, err
	}
	nodeList, err := instance.GetBroadcastingNodes(nil)
	if err != nil {
		return nil, err
	}
	return filterNodes(nodeList), nil
}

// View information of storage nodes
func viewStorageNodeInfo() ([]bindings.NodeManagerNodeInfo, error) {
	_, instance, err := getNodeManagerInstance()
	if err != nil {
		return nil, err
	}
	nodeList, err := instance.GetStorageNodes(nil)
	if err != nil {
		return nil, err
	}
	return filterNodes(nodeList), nil
}

// Register a new node group
func registerNodeGroup() common.Hash {
	// Get user input for node addresses and required signatures
	addresses := getCommaSeparatedInput("Enter a list of broadcast node addresses separated by commas:\n")
	addressList := parseAddresses(addresses)
	requiredAmountOfSignatures := getIntInput("Please enter the minimum number of signatures: \n")

	// Get instance of StorageManager contract
	client, instance, auth, address := getStorageManagerInstance()
	if client == nil || instance == nil || auth == nil || address == nil {
		return common.Hash{}
	}

	// Get node group key
	nodeGroupKey := getNodeGroupKey(instance, address, addressList, big.NewInt(requiredAmountOfSignatures))
	nodegroup, err := instance.NODEGROUP(&bind.CallOpts{
		Pending: false,
		From:    *address,
	}, nodeGroupKey)
	if err != nil || nodeGroupKey.Cmp(common.Hash{}) == 0 {
		log.Fatal(err)
	}

	// Check if the node group is already registered
	if len(nodegroup.Addrs) > 0 {
		log.Println("The nodeGroup has been registered and can be used directly.")
		return nodeGroupKey
	}

	// Register the node group
	tx := registerNodeGroupTransaction(instance, auth, big.NewInt(requiredAmountOfSignatures), addressList)
	waitForTransaction(tx, client)
	return nodeGroupKey
}

// Register a new namespace
func registerNameSpace() common.Hash {
	// Get user input for node addresses
	addresses := getCommaSeparatedInput("Enter a list of storage node addresses separated by commas:\n")
	addressList := parseAddresses(addresses)

	// Get user input for node addresses
	client, instance, auth, address := getStorageManagerInstance()
	if client == nil || instance == nil || auth == nil || address == nil {
		return common.Hash{}
	}

	// Get namespace key
	nameSpaceKey := getNameSpaceKey(instance, address, addressList)
	nameSpace, err := instance.NAMESPACE(&bind.CallOpts{
		Pending: false,
		From:    *address,
	}, nameSpaceKey)
	if err != nil || nameSpaceKey.Cmp(common.Hash{}) == 0 {
		log.Fatal(err)
	}

	// Check if the namespace is already registered
	if len(nameSpace.Addr) > 0 {
		log.Println("The nodeGroup has been registered and can be used directly.")
		return nameSpaceKey
	}

	// Register the namespace
	tx := registerNameSpaceTransaction(instance, auth, addressList)
	waitForTransaction(tx, client)
	return nameSpaceKey
}

// Perform a general test with simulated data
func generalTest() {
	fmt.Println("Running Parameter Test...")
	const nodeGroupKeyStr = "69506f24e7f886f0e3dd0a6f2137da279798d7bcfaef2071648748d6930594fe"
	nameSpaceKey := common.HexToHash("0x00")
	for {
		sendDA(nodeGroupKeyStr, nameSpaceKey)
		time.Sleep(5 * time.Minute)
	}
}

// Perform an advanced test with user inputs
func advancedTest() {
	var nodeGroupKeyStr string
	nodeGroupKeyInput := &survey.Input{
		Message: "Please enter NodeGroupKey:\n",
	}
	err := survey.AskOne(nodeGroupKeyInput, &nodeGroupKeyStr)
	if err != nil || nodeGroupKeyStr == "" {
		log.Fatal("NodeGroupKey cannot be empty")
	}
	var nameSpaceKeyStr string
	nameSpaceKeyInput := &survey.Input{
		Message: "Please enter NameSpaceKey (press Enter if empty):\n",
	}
	err = survey.AskOne(nameSpaceKeyInput, &nameSpaceKeyStr)
	if err != nil {
		log.Fatal(err.Error())
	}

	var timerInt int
	timerInput := &survey.Input{
		Message: "Please enter the sending interval (seconds):\n",
	}
	err = survey.AskOne(timerInput, &timerInt)
	if err != nil {
		log.Fatal("Please enter a valid time interval")
	}

	nameSpaceKey := common.HexToHash(nameSpaceKeyStr)
	for {
		sendDA(nodeGroupKeyStr, nameSpaceKey)
		time.Sleep(time.Duration(timerInt) * time.Second)
	}
}

func registerBroadcastNode() {
	url, name, location, stakedTokens, maxStorageSpace := getRegisterNodeInfo()

	privateKeyHex := os.Getenv("PRIVATE_KEY")
	private, address := privateKeyToAddress(privateKeyHex)
	client, instance, err := getNodeManagerInstance()
	if err != nil {
		log.Fatalf("cant create contract address err: %s", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(private, big.NewInt(chainID))
	if err != nil {
		log.Fatal(err)
	}
	tx, err := instance.RegisterBroadcastNode(auth, bindings.NodeManagerNodeInfo{
		Url:             url,
		Name:            name,
		StakedTokens:    big.NewInt(stakedTokens),
		Location:        location,
		MaxStorageSpace: big.NewInt(maxStorageSpace),
		Addr:            address,
	})

	if err != nil {
		log.Fatal(err.Error())
		return
	}
	waitForTransaction(tx, client)
}

func registerStorageNode() {
	url, name, location, stakedTokens, maxStorageSpace := getRegisterNodeInfo()

	privateKeyHex := os.Getenv("PRIVATE_KEY")
	private, address := privateKeyToAddress(privateKeyHex)
	client, instance, err := getNodeManagerInstance()
	if err != nil {
		log.Fatalf("cant create contract address err: %s", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(private, big.NewInt(chainID))
	if err != nil {
		log.Fatal(err)
	}
	tx, err := instance.RegisterStorageNode(auth, bindings.NodeManagerNodeInfo{
		Url:             url,
		Name:            name,
		StakedTokens:    big.NewInt(stakedTokens),
		Location:        location,
		MaxStorageSpace: big.NewInt(maxStorageSpace),
		Addr:            address,
	})

	if err != nil {
		log.Fatal(err.Error())
		return
	}
	waitForTransaction(tx, client)
}

// Register a node group transaction
func registerNodeGroupTransaction(instance *bindings.StorageManager, auth *bind.TransactOpts, requiredAmountOfSignatures *big.Int, addressList []common.Address) *types.Transaction {
	tx, err := instance.RegisterNodeGroup(auth, requiredAmountOfSignatures, addressList)
	if err != nil {
		log.Fatal(err)
	}
	return tx
}

// Register a namespace transaction
func registerNameSpaceTransaction(instance *bindings.StorageManager, auth *bind.TransactOpts, addressList []common.Address) *types.Transaction {
	tx, err := instance.RegisterNameSpace(auth, addressList)
	if err != nil {
		log.Fatal(err)
	}
	return tx
}

// Wait for the transaction to be mined
func waitForTransaction(tx *types.Transaction, client *ethclient.Client) {
	log.Printf("tx sent: %s", tx.Hash().Hex())
	_, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("tx confirmed")
}

// Send DA with simulated data
func sendDA(nodeGroupKeyStr string, nameSpaceKey [32]byte) {
	sdk, err := kzgsdk.InitMultiAdaptiveSdk("./srs")
	if err != nil {
		log.Fatalf("kzgsdk Error: %s", err)
	}

	data := simulatedData()

	cm, proof, err := sdk.GenerateDataCommitAndProof(data)
	if err != nil {
		log.Fatalf("kzgsdk Error: %s", err)
	}

	nodeGroupKey := common.HexToHash(nodeGroupKeyStr)
	addressHex := os.Getenv("ADDRESS")
	sender := common.HexToAddress(addressHex)

	index, err := getIndex(sender)
	if err != nil {
		log.Fatalf("getIndex Error: %s", err)
	}

	ti := time.Now()
	timeout := ti.Add(10 * time.Hour).Unix()

	signatures, err := GetSignature(nodeGroupKey, sender, index, uint64(len(data)), cm.Marshal(), data, proof.H.Marshal(), proof.ClaimedValue.Marshal(), uint64(timeout))
	if err != nil {
		log.Printf("GetSignature Error: %s", err)
	}

	SendCommitToL1(uint64(len(data)), nodeGroupKey, signatures, cm, nameSpaceKey, timeout)
}

// Generate simulated data
func simulatedData() []byte {
	data := make([]byte, dataSize)
	rand.Read(data)
	return data
}

func getIndex(sender common.Address) (uint64, error) {
	_, instance, err := getCommitmentManagerInstance()
	if err != nil {
		return 0, err
	}
	index, err := instance.Indices(nil, sender)
	if err != nil {
		return 0, err
	}
	return index.Uint64(), nil
}

func GetSignature(nodeGroupKey common.Hash, sender common.Address, index, length uint64, commitment, data, proof, claimedValue []byte, timeout uint64) (signatures [][]byte, err error) {
	client, err := ethclient.Dial(ethUrl)
	if err != nil {
		return nil, err
	}

	storageManager, err := bindings.NewStorageManager(common.HexToAddress(storageManagerAddress), client)
	if err != nil {
		return nil, err
	}

	nodeManager, err := bindings.NewNodeManager(common.HexToAddress(nodeManagerAddress), client)
	if err != nil {
		return nil, err
	}

	nodeGroup, err := storageManager.NODEGROUP(nil, nodeGroupKey)
	if err != nil {
		return nil, err
	}

	for _, add := range nodeGroup.Addrs {
		info, err := nodeManager.BroadcastingNodes(nil, add)
		if err != nil {
			signatures = append(signatures, nil)
			continue
		}
		sign, err := signature(info.Url, sender, index, length, commitment, data, nodeGroupKey, proof, claimedValue, timeout)
		if err != nil {
			signatures = append(signatures, nil)
			continue
		}
		signatures = append(signatures, sign)
	}

	return signatures, nil
}

func signature(url string, sender common.Address, index, length uint64, commitment, data []byte, nodeGroupKey [32]byte, proof, claimedValue []byte, timeout uint64) ([]byte, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	var result []byte
	err = client.Client().CallContext(ctx, &result, "mta_sendDAByParams", sender, index, length, commitment, data, nodeGroupKey, proof, claimedValue, timeout)
	return result, err
}

func SendCommitToL1(length uint64, dasKey [32]byte, sign [][]byte, commit kzg.Digest, nameSpaceId [32]byte, timeout int64) {
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	private, _ := privateKeyToAddress(privateKeyHex)
	client, instance, err := getCommitmentManagerInstance()
	if err != nil {
		log.Fatalf("cant create contract address err: %s", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(private, big.NewInt(chainID))
	if err != nil {
		log.Fatal(err)
	}

	commitData := bindings.PairingG1Point{
		X: new(big.Int).SetBytes(commit.X.Marshal()),
		Y: new(big.Int).SetBytes(commit.Y.Marshal()),
	}

	tx, err := instance.SubmitCommitment(auth, big.NewInt(int64(length)), big.NewInt(timeout), nameSpaceId, dasKey, sign, commitData)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	waitForTransaction(tx, client)
}

// Display node information by calling the respective function
func displayNodeInfo(fn func() ([]bindings.NodeManagerNodeInfo, error)) {
	list, err := fn()
	if err != nil {
		log.Printf(err.Error())
		return
	}
	for i, info := range list {
		log.Printf("%d Url:%s  Address:%s  Name:%s  Location:%s  StakedTokens:%s  MaxStorageSpace:%s", i, info.Url, info.Addr, info.Name, info.Location, info.StakedTokens, info.MaxStorageSpace)
	}
}

// Filter nodes based on registration and expiry times
func filterNodes(nodeList []bindings.NodeManagerNodeInfo) []bindings.NodeManagerNodeInfo {
	var filtered []bindings.NodeManagerNodeInfo
	for _, info := range nodeList {
		if info.StakedTokens.Cmp(big.NewInt(0)) != 0 {
			filtered = append(filtered, info)
		}
	}
	return filtered
}

// Get NodeManager instance
func getNodeManagerInstance() (*ethclient.Client, *bindings.NodeManager, error) {
	client, err := ethclient.Dial(ethUrl)
	if err != nil {
		return nil, nil, err
	}
	instance, err := bindings.NewNodeManager(common.HexToAddress(nodeManagerAddress), client)
	if err != nil {
		return nil, nil, err
	}
	return client, instance, nil
}

func getCommitmentManagerInstance() (*ethclient.Client, *bindings.CommitmentManager, error) {
	client, err := ethclient.Dial(ethUrl)
	if err != nil {
		return nil, nil, err
	}
	instance, err := bindings.NewCommitmentManager(common.HexToAddress(cmManagerAddress), client)
	if err != nil {
		return nil, nil, err
	}
	return client, instance, nil
}

// Get StorageManager instance
func getStorageManagerInstance() (*ethclient.Client, *bindings.StorageManager, *bind.TransactOpts, *common.Address) {
	client, err := ethclient.Dial(ethUrl)
	if err != nil {
		log.Fatal(err)
		return nil, nil, nil, nil
	}
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	private, address := privateKeyToAddress(privateKeyHex)
	auth, err := bind.NewKeyedTransactorWithChainID(private, big.NewInt(chainID))
	if err != nil {
		log.Fatal(err)
		return nil, nil, nil, nil
	}
	instance, err := bindings.NewStorageManager(common.HexToAddress(storageManagerAddress), client)
	if err != nil {
		log.Fatalf("cant create contract address err: %v", err)
		return nil, nil, nil, nil
	}
	return client, instance, auth, &address
}

func getRegisterNodeInfo() (url, name, location string, stakedTokens, maxStorageSpace int64) {
	url = getCommaSeparatedInput("url")
	name = getCommaSeparatedInput("name")
	location = getCommaSeparatedInput("location")
	stakedTokens = getIntInput("stakedTokens")
	maxStorageSpace = getIntInput("maxStorageSpace")
	return url, name, location, stakedTokens, maxStorageSpace
}

// Get comma-separated input from the user
func getCommaSeparatedInput(prompt string) string {
	var input string
	addressesInput := &survey.Input{
		Message: prompt,
	}
	err := survey.AskOne(addressesInput, &input)
	if err != nil {
		log.Fatal(err.Error())
	}
	return input
}

// Get integer input from the user
func getIntInput(prompt string) int64 {
	var input int64
	inputPrompt := &survey.Input{
		Message: prompt,
	}
	err := survey.AskOne(inputPrompt, &input)
	if err != nil {
		log.Fatal(err.Error())
	}
	return input
}

// Parse addresses from a comma-separated string
func parseAddresses(addresses string) []common.Address {
	addressList := strings.Split(addresses, ",")
	var commonAddresses []common.Address
	for _, addr := range addressList {
		trimmedAddr := strings.TrimSpace(addr)
		if !common.IsHexAddress(trimmedAddr) {
			log.Fatalf("Invalid address: %s", trimmedAddr)
		}
		commonAddresses = append(commonAddresses, common.HexToAddress(trimmedAddr))
	}
	return commonAddresses
}

func privateKeyToAddress(privateKeyHex string) (*ecdsa.PrivateKey, common.Address) {
	private, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Invalid private key: %v", err)
	}
	address := crypto.PubkeyToAddress(private.PublicKey)
	return private, address
}

func getNodeGroupKey(instance *bindings.StorageManager, address *common.Address, addressList []common.Address, requiredAmountOfSignatures *big.Int) common.Hash {
	nodeGroupKey, err := instance.GetNodeGroupKey(&bind.CallOpts{
		Pending: false,
		From:    *address,
	}, addressList, requiredAmountOfSignatures)
	if err != nil {
		return common.Hash{}
	}
	return nodeGroupKey
}

func getNameSpaceKey(instance *bindings.StorageManager, address *common.Address, addressList []common.Address) common.Hash {
	nameSpaceKey, err := instance.GetNameSpaceKey(&bind.CallOpts{
		Pending: false,
		From:    *address,
	}, addressList)
	if err != nil {
		return common.Hash{}
	}
	return nameSpaceKey
}
