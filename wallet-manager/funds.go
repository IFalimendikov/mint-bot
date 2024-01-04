package walletmanager

import (
	"crypto/ecdsa"
	"context"
    "fmt"
    "log"
    "os"
	"io"
    "encoding/json"
	"strings"
	"math/big"
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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



func getABI(contractAddress string) (abi.ABI, error) {
	url := fmt.Sprintf("https://api.etherscan.io/api?module=contract&action=getabi&address=%s&apikey=YourApiKeyToken", contractAddress)

	resp, err := http.Get(url)
	if err != nil {
		return abi.ABI{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return abi.ABI{}, err
	}

	var etherscanResponse EtherscanResponse
	err = json.Unmarshal(body, &etherscanResponse)
	if err != nil {
		return abi.ABI{}, err
	}

	if etherscanResponse.Status != "1" {
		return abi.ABI{}, fmt.Errorf("Error: %s", etherscanResponse.Message)
	}

	parsedAbi, err := abi.JSON(strings.NewReader(etherscanResponse.Result))
	if err != nil {
		return abi.ABI{}, err
	}

	return parsedAbi, nil
}

func getIDs(contractAddress string, client *ethclient.Client) ([]uint256, error) {
	abi, err := getABI(contractAddress)
	if err != nil {
		return nil, err
	}

	address := common.HexToAddress(contractAddress)
	contract := bind.NewBoundContract(address, abi, client, client, client)

	totalSupply, err := contract.TotalSupply(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	ids := make([]uint256, totalSupply)
	for i := uint256(0); i < totalSupply; i++ {
		id, err := contract.TokenOfOwnerByIndex(&bind.CallOpts{}, common.HexToAddress(fundingWallet), i)
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}

	return ids, nil
}




func transferNFT(walletGroup string, client *ethclient.Client, contractAddress, fundingWallet string) error {
	wallets := strings.ToLower(walletGroup)

	if _, err := os.Stat(fmt.Sprintf("wallet-manager/wallet-data/%s", wallets)); err == nil {
		walletList := []Wallet{}

		err := json.Unmarshal([]byte(fmt.Sprintf("wallet-manager/wallet-data/%s", wallets)), &walletList)
		if err != nil {
			return err
		}

        abi, err := getABI(contractAddress)
        if err != nil {
            return err
        }

		chainID, err := client.NetworkID(context.Background())

		for _, w := range walletList {
			// Set Private Key
			privateKey, err := crypto.HexToECDSA(w.PrivateKey)
			if err != nil {
                log.Fatal(err)
            }
			// Create a transactor for each wallet
			auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
			if err != nil {
				return err
			}
			
			// Create a new instance of the contract for each wallet
			contract := bind.NewBoundContract(common.HexToAddress(contractAddress), abi, client, client, client)
		  
			// Call the Transfer function
			tx, err := contract.Transact(auth, "setState", )
			if err != nil {
				log.Fatalf("Failed to call state-changing method: %v", err)
			}
			fmt.Printf("Transaction sent: %s\n", tx.Hash().Hex())
		  }
	}
}
