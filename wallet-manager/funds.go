package walletmanager

import (
	"crypto/ecdsa"
    "fmt"
    "log"
    "os"
    "encoding/json"
	"strings"
	"math/big"

	_ "github.com/joho/godotenv/autoload"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

func init() {
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
	FundingWallet := crypto.PubkeyToAddress(*publicKeyECDSA)
}

func FundWallets(walletGroup string, ethAmount float64) {
	wallets := strings.ToLower(walletGroup)

	if _, err := os.Stat(fmt.Sprintf("wallet-manager/wallet-data/%s", wallets)); err == nil {
		walletList := []Wallet{}

		err := json.Unmarshal([]byte(fmt.Sprintf("wallet-manager/wallet-data/%s", wallets)), &walletList)
		if err != nil {
			log.Fatal(err)
		}

		for _, w := range walletList {
			// Declare Transaction Parameters
			nonce, err := client.PendingNonceAt(context.Background(), FundingWallet)
			if err != nil {
				log.Println(err)
			}

			amount := new(big.Int)
			amount.SetString(fmt.Sprintf("%.0f", ethAmount*params.Ether), 10)
			gasLimit := uint64(21000)
			gasPrice, err := client.SuggestGasPrice(context.Background())
			if err != nil {
				log.Println(err)
			}

			chainID, err := client.NetworkID(context.Background())
			if err != nil {
				log.Println(err)
			}

			// Create and Sign Transaction
			transaction := types.NewTransaction(nonce, w.Address, amount, gasLimit, gasPrice, nil)
			signedTx, err := types.SignTx(transaction, types.NewEIP155Signer(chainID), privateKey)
			if err != nil {
				log.Fatal(err)
			}

			err = client.SendTransaction(context.Background(), signedTx)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Transaction sent: %s", signedTx.Hash().Hex())
		}
	}
}

func CheckBalance(walletGroup string) {
	wallets := strings.ToLower(walletGroup)

	if _, err := os.Stat(fmt.Sprintf("wallet-manager/wallet-data/%s", wallets)); err == nil {
		walletList := []Wallet{}

		err := json.Unmarshal([]byte(fmt.Sprintf("wallet-manager/wallet-data/%s", wallets)), &walletList)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Wallet group %s\n", walletGroup)

		for _, w := range walletList {
			balance, err := client.BalanceAt(context.Background(), common.HexToAddress(w.Address), nil)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Wallet %s balance is: %s\n", w.Address, balance)
		}
	}
}

func WithdrawFunds(walletGroup string) {
	wallets := strings.ToLower(walletGroup)

	if _, err := os.Stat(fmt.Sprintf("wallet-manager/wallet-data/%s", wallets)); err == nil {
		walletList := []Wallet{}

		err := json.Unmarshal([]byte(fmt.Sprintf("wallet-manager/wallet-data/%s", wallets)), &walletList)
		if err != nil {
			log.Fatal(err)
		}

        for _, w := range walletList {
            // Declare Transaction Parameters
            nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(w.address))
            if err != nil {
                log.Fatal(err)
            }

            balance, err := client.BalanceAt(context.Background(), common.HexToAddress(w.address), nil)
            if err != nil {
                log.Fatal(err)
            }

            gasPrice, err := client.SuggestGasPrice(context.Background())
            if err != nil {
                log.Fatal(err)
            }

            gasLimit := uint64(21000) // this is a standard gas limit for ETH transfer
            totalGas := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))

            amount := new(big.Int).Sub(balance, totalGas)

            chainID, err := client.NetworkID(context.Background())

            // Create and Sign Transaction
            transaction := types.NewTransaction(nonce, FundingWallet, amount, gasLimit, gasPrice, nil)
            signedTx, err := types.SignTx(transaction, types.NewEIP155Signer(chainID), w.privateKey)
            if err != nil {
                log.Fatal(err)
            }

            err = client.SendTransaction(context.Background(), signedTx)
            if err != nil {
                log.Fatal(err)
            }

            fmt.Printf("Transaction sent: %s", signedTx.Hash().Hex())
		}
	}
}

func WithdrawNFTS(walletGroup string, nftContractAddress string) {
	wallets := strings.ToLower(walletGroup)

	if _, err := os.Stat(fmt.Sprintf("wallet-manager/wallet-data/%s", wallets)); err == nil {
		walletList := []Wallet{}

		err := json.Unmarshal([]byte(fmt.Sprintf("wallet-manager/wallet-data/%s", wallets)), &walletList)
		if err != nil {
			log.Fatal(err)
		}

        for _, w := range walletList {
            // Declare Transaction Parameters
            nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(w.address))
            if err != nil {
                log.Fatal(err)
            }

            // Declare NFT Transfer Parameters
            gasPrice, err := client.SuggestGasPrice(context.Background())
            if err != nil {
                log.Fatal(err)
            }

            gasLimit, err := client.EstimateGas(context.Background(), msg) 
            totalGas := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))

            // Get NFT Contract
            nftContract, err := abi.JSON(strings.NewReader(nftABI))
            if err != nil {
                log.Fatal(err)
            }

            // Transfer NFT
            tx := types.NewContractCreation
		}
	}
}