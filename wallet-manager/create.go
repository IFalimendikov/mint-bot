package walletmanager

import (
    "crypto/ecdsa"
    "fmt"
    "log"
    "os"
    "encoding/json"

    "github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/ethereum/go-ethereum/crypto"
)

type Wallet struct {
    Address string `json:"address"`
    PrivateKey string `json:"privateKey"`
}

func CreateWallets() {
    wallets := make([]Wallet, 10)

    for i := 0; i < 10; i++ {
        privateKey, err := crypto.GenerateKey()
        if err != nil {
            log.Fatal(err)
        }

        privateKeyBytes := crypto.FromECDSA(privateKey)
        // fmt.Println("SAVE BUT DO NOT SHARE THIS (Private Key):", hexutil.Encode(privateKeyBytes))

        publicKey := privateKey.Public()
        publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
        if !ok {
            log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
        }

        // publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
        // fmt.Println("Public Key:", hexutil.Encode(publicKeyBytes)) 

        address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
        // fmt.Println("Address:", address)

        wallets[i] = Wallet{
            Address: address,
            PrivateKey: hexutil.Encode(privateKeyBytes),
        }
    }

    file, _ := json.MarshalIndent(wallets, "", " ")
    _ = os.WriteFile("wallet-manager/wallet-data/wallets.json", file, 0644)

    fmt.Println("Success! 10 new wallets have been created and written to wallets.json")
} 