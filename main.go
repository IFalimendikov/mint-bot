package main

import (
	"os"
    "fmt"
    "strconv"

    "github.com/IFalimendikov/mint-bot/wallet-manager"
)

func main() {
	
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
            walletmanager.FundWallets(walletGroup, ethAmount)
        } else {
            fmt.Println("Wallet group and ETH amount are required for funding! Usage: go run main.go fundWallets [walletGroup] [amount]")
        }
    case "checkBalance":
        if arg2 != "" {
            walletGroup := arg2
            walletmanager.CheckBalance(walletGroup)
        } else {
            fmt.Println("Wallet group are required for wallet balance checking! Usage: go run main.go checkBalance [walletGroup]")
        }
    default:
        fmt.Println("Invalid argument. Please use 'createWallets' or 'fundWallets'.")
    }
	
}