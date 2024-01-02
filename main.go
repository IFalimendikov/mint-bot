package main

import (
	"os"
    "fmt"
    "log"
    "strconv"
    

    "github.com/IFalimendikov/mint-bot/wallet-manager"
)

func main() {

    // Connect to an Ethereum node
	client, err := ethclient.Dial(os.Getenv("ETH_NODE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(os.Getenv("FUNDING_WALLET_PRIVATE_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	// Declare Sender Address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Public Key Error")
	}
	fundingWallet := crypto.PubkeyToAddress(*publicKeyECDSA)
	
    // The first argument
    // is always program name
    // So os.Args[1] is the first dynamic argument
    arg1 := os.Args[1]
    arg2 := os.Args[2]
    arg3 := os.Args[3]

    // use arg1 to decide which packages to call
    switch arg1 {
    case "createWallets":
        // option1 code executes here.
        walletmanager.CreateWallets()
    case "fundWallets":
        // option2 code executes here.
        if arg2 != "" && arg3 != "" {
            walletGroup := arg2
            ethAmount, _ := strconv.ParseFloat(arg3, 64)
            walletmanager.FundWallets(walletGroup, ethAmount, client, fundingWallet, privateKey)
        } else {
            fmt.Println("Wallet group and ETH amount are required for funding! Usage: go run main.go fundWallets [walletGroup] [amount]")
        }
    case "checkBalance":
        if arg2 != "" {
            walletGroup := arg2
            walletmanager.CheckBalance(walletGroup, client)
        } else {
            fmt.Println("Wallet group is required for wallet balance checking! Usage: go run main.go checkBalance [walletGroup]")
        }
    case "withdawFunds":
        if arg2 != "" {
            walletGroup := arg2
            walletmanager.WithdrawFunds(walletGroup, client, fundingWallet)
        } else {
            fmt.Println("Wallet group is required for funds withdrawal! Usage: go run main.go withdrawFunds [walletGroup]")
        }
    case "withdrawNFTS":
        if arg2 != "" && arg3 != "" {
            walletGroup := arg2
            nftAddress := arg3
            walletmanager.WithdrawNFTS(walletGroup, nftAddress, client, fundingWallet)
        } else {
            fmt.Println("Wallet group and NFT address is required for NFT withdrawal! Usage: go run main.go withdrawNFTS [walletGroup] [nftAddress]")
        }
    default:
        fmt.Println("Invalid argument. Please use 'createWallets' or 'fundWallets'.")
    }
}