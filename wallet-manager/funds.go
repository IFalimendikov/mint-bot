package walletmanager

import (
	"crypto/ecdsa"
	"context"
    "fmt"
    "log"
    "os"
    "encoding/json"
	"strings"
	"math/big"

	_ "github.com/joho/godotenv/autoload"

	"github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/core/types"
)

func FundWallets(walletGroup string, ethAmount float64, client *ethclient.Client, FundingWallet string, privateKey *ecdsa.PrivateKey) {
	wallets := strings.ToLower(walletGroup)

	if _, err := os.Stat(fmt.Sprintf("wallet-manager/wallet-data/%s", wallets)); err == nil {
		walletList := []Wallet{}

		err := json.Unmarshal([]byte(fmt.Sprintf("wallet-manager/wallet-data/%s", wallets)), &walletList)
		if err != nil {
			log.Fatal(err)
		}

		for _, w := range walletList {
			// Declare Transaction Parameters
			nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(FundingWallet))
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

			// Convert wallet address string to Address type
			walletAddress := common.HexToAddress(w.Address)

			// Create and Sign Transaction
			transaction := types.NewTransaction(nonce, walletAddress, amount, gasLimit, gasPrice, nil)
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


func CheckBalance(walletGroup string, client *ethclient.Client) {
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

func WithdrawFunds(walletGroup string, client *ethclient.Client, fundingWallet string) {
	wallets := strings.ToLower(walletGroup)

	if _, err := os.Stat(fmt.Sprintf("wallet-manager/wallet-data/%s", wallets)); err == nil {
		walletList := []Wallet{}

		err := json.Unmarshal([]byte(fmt.Sprintf("wallet-manager/wallet-data/%s", wallets)), &walletList)
		if err != nil {
			log.Fatal(err)
		}

        for _, w := range walletList {
            // Declare Transaction Parameters
            nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(w.Address))
            if err != nil {
                log.Fatal(err)
            }

            balance, err := client.BalanceAt(context.Background(), common.HexToAddress(w.Address), nil)
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

			privateKey, err := crypto.HexToECDSA(w.PrivateKey)
			if err != nil {
                log.Fatal(err)
            }

            // Create and Sign Transaction
            transaction := types.NewTransaction(nonce, common.HexToAddress(fundingWallet), amount, gasLimit, gasPrice, nil)
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


type EtherscanResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func getABI(contractAddress string) (string, error) {
	url := fmt.Sprintf("https://api.etherscan.io/api?module=contract&action=getabi&address=%s&apikey=YourApiKeyToken", contractAddress)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var etherscanResponse EtherscanResponse
	err = json.Unmarshal(body, &etherscanResponse)
	if err != nil {
		return "", err
	}

	if etherscanResponse.Status != "1" {
		return "", fmt.Errorf("Error: %s", etherscanResponse.Message)
	}

	return etherscanResponse.Result, nil
}

func getNFTIDs(walletAddress string, nftContractAddress string) ([]*big.Int, error) {
    // Get NFT Contract
    abi, err := getABI(nftContractAddress)
    if err != nil {
        return nil, err
    }

    address := common.HexToAddress(walletAddress)
    balance, err := abi.BalanceOf(&bind.CallOpts{}, address)
    if err != nil {
        return nil, err
    }

    var ids []*big.Int
    for i := big.NewInt(0); i.Cmp(balance) < 0; i.Add(i, big.NewInt(1)) {
        id, err := abi.TokenOfOwnerByIndex(&bind.CallOpts{}, address, i)
        if err != nil {
            return nil, err
        }
        ids = append(ids, id)
    }

    return ids, nil
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

			// Get all NFT IDs owned by the wallet
			ids, err := getNFTIDs(w.address, nftContractAddress)
			if err != nil {
				log.Fatal(err)
			}

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
            abi, err := getABI(nftContractAddress)
            if err != nil {
                log.Fatal(err)
            }

			
		}
	}
}