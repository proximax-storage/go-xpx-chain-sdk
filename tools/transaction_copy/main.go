package main

import (
	"context"
	"fmt"
	"time"
	"flag"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

const (
	// Sirius api rest server
	baseUrl = "http://localhost:3000"
	mainurl = "https://betelgeuse.xpxsirius.io"
	// Valid private key
	privateKey = "819F72066B17FFD71B8B4142C5AEAE4B997B0882ABDF2C263B02869382BD93A0"
	
	// Mainnet mosaic id
	xpx  = 0x402B2F579FAEBC59
	so   = 0x42F11DF73D00296A
	sm   = 0x7E332C41B04E4E05
	sc   = 0x373215BB26B0DB2C
	metx = 0x7FAE09FA0288ED69

	// Testnet mosaic id
	xpx2  = 0x0DC67FBE1CAD29E3
	sm2   = 0x6C5D687508AC9D75
	so2   = 0x26514E2A1EF33824
	sc2   = 0x6EE955268A1C33D9
	metx2 = 0x00

)

func main() {
	// Blockheight to check
	height := flag.Uint("height", 0, "height to start copy")

	for {
		fmt.Println("job start", height)
		task(height)
		*height++
		time.Sleep(15 * time.Second)
	}
}

func task(height *uint) {
	mainnetConfig, err := sdk.NewConfig(context.Background(), []string{mainurl})
	if err != nil {
		fmt.Printf("NewConfig returned error: %s", err)
		return
	}

	// Use the default http client
	mainnetClient := sdk.NewClient(nil, mainnetConfig)

	testnetConfig, err := sdk.NewConfig(context.Background(), []string{baseUrl})
	if err != nil {
		fmt.Printf("NewConfig returned error: %s", err)
		return
	}

	// Use the default http client
	testnetClient := sdk.NewClient(nil, testnetConfig)

	options := &sdk.TransactionsPageOptions{Height: *height, FirstLevel: false}

	// Get transaction informations for transactionIds or transactionHashes
	transactionPage, err := mainnetClient.Transaction.GetTransactionsByGroup(context.Background(), sdk.Confirmed, options)
	if err != nil {
		fmt.Printf("Transaction.GetTransactions returned error: %s", err)
		return
	}

	// Loop through all transactions in specific height
	for _, t := range transactionPage.Transactions {

		switch t.GetAbstractTransaction().Type {
		case sdk.Transfer:
			fmt.Println("transfer switch")
			txn := t.(*sdk.TransferTransaction)

			// Create signer clone
			signer := t.GetAbstractTransaction().Signer
			clone, err := testnetClient.NewAccountFromPrivateKey(signer.PublicKey)
			if err != nil {
				fmt.Printf("Clone failed to create: %s", err)
				return
			}
			signerInfo, err := mainnetClient.Account.GetAccountInfo(context.Background(), signer.Address)
			if err != nil {
				fmt.Printf("SignerInfo fail to get: %s", err)
				return
			}
			cloneInfo, err := testnetClient.Account.GetAccountInfo(context.Background(), clone.Address)
			if err != nil {
				fmt.Printf("CloneInfo fail to get: %s", err)
				return
			}
			if cloneInfo == nil {
				restored := restoreMosaics(signerInfo.Mosaics, txn)
				balance := pairMosaics(restored)
				success := createCloneAccount(testnetClient, signerInfo, balance)
				// success := assignSigner(testnetClient, newAcc, balance, genesis)
				if !success {
					return
				}
			}
			
			// Create recipient clone
			newRecipient, err := testnetClient.NewAccount()
			if err != nil {
				fmt.Printf("Failed to create recipient clone: %s", err)
				return
			}
			txn.Recipient = newRecipient.Address
			
			// Convert mosaic id
			toTransferedMosaic := pairMosaics(txn.Mosaics)
			txn.Mosaics = toTransferedMosaic
			
			// Update txn info
			t.GetAbstractTransaction().Deadline = sdk.NewDeadline(time.Hour * 1)
			t.GetAbstractTransaction().NetworkType = testnetConfig.NetworkType

			// Sign transaction
			signedTransaction, err := clone.Sign(transactionPage.Transactions[0])
			if err != nil {
				fmt.Printf("Sign returned error: %s", err)
				return
			}
			
			// Announce transaction
			_, err = testnetClient.Transaction.Announce(context.Background(), signedTransaction)
			if err != nil {
				fmt.Printf("Transaction.Announce returned error: %s", err)
				return
			}

		case sdk.AggregateCompleted, sdk.AggregateBonded:
			fmt.Println("Aggregate switch")
			
			cosigners := []*sdk.Account{}
			aggTx := transactionPage.Transactions[0].(*sdk.AggregateTransaction)

			// Create signer clone
			signer := t.GetAbstractTransaction().Signer
			clone, err := testnetClient.NewAccountFromPrivateKey(signer.PublicKey)
			if err != nil {
				fmt.Printf("Clone failed to create: %s", err)
				return
			}
			signerInfo, err := mainnetClient.Account.GetAccountInfo(context.Background(), signer.Address)
			if err != nil {
				fmt.Printf("SignerInfo fail to get: %s", err)
				return
			}
			cloneInfo, err := testnetClient.Account.GetAccountInfo(context.Background(), clone.Address)
			if err != nil {
				fmt.Printf("CloneInfo fail to get: %s", err)
				return
			}
			if cloneInfo == nil {
				balance := pairMosaics(signerInfo.Mosaics)
				success := createCloneAccount(testnetClient, signerInfo, balance)
				if !success {
					return
				}
			}

			// Region add cosignature of agg txn
			for _, cosigner := range aggTx.Cosignatures {
				signer := cosigner.Signer
				clone, err:= testnetClient.NewAccountFromPrivateKey(signer.PublicKey)
				if err != nil {
					fmt.Printf("Clone failed to create: %s", err)
					return
				}
				signerInfo, err := mainnetClient.Account.GetAccountInfo(context.Background(), signer.Address)
				if err != nil {
					fmt.Printf("SignerInfo fail to get: %s", err)
					return
				}
				cloneInfo, err := testnetClient.Account.GetAccountInfo(context.Background(), clone.Address)
				if err != nil {
					fmt.Printf("CloneInfo fail to get: %s", err)
					return
				}
				if cloneInfo == nil {
					balance := pairMosaics(signerInfo.Mosaics)
					success := createCloneAccount(testnetClient, signerInfo, balance)
					if !success {
						return
					}
				}
				cosigners = append(cosigners, clone)
			}

			for _, t := range aggTx.InnerTransactions {
				clone, err := testnetClient.NewAccountFromPrivateKey(t.GetAbstractTransaction().Signer.PublicKey)
				if err != nil {
					fmt.Printf("CloneInfo fail to get: %s", err)
					return
				}
				
				if t.GetAbstractTransaction().Type == sdk.Transfer {
					txn := t.(*sdk.TransferTransaction)
					
					// Create recipient clone
					newRecipient, _ := testnetClient.NewAccount()
					txn.Recipient = newRecipient.Address
					
					// Convert mosaic id
					toTransferedMosaic := pairMosaics(txn.Mosaics)
					txn.Mosaics = toTransferedMosaic
					fmt.Println(toTransferedMosaic)
					fmt.Println(txn.Mosaics)
				}

				t.GetAbstractTransaction().ToAggregate(clone.PublicAccount)
				t.GetAbstractTransaction().Deadline = sdk.NewDeadline(time.Hour * 1)
			}

			t.GetAbstractTransaction().Deadline = sdk.NewDeadline(time.Hour * 1)
			t.GetAbstractTransaction().NetworkType = testnetConfig.NetworkType
			t.GetAbstractTransaction().Type = sdk.AggregateCompleted

			// Sign transaction
			signedTransaction, err := clone.SignWithCosignatures(aggTx, cosigners)
			if err != nil {
				fmt.Printf("Transaction.Announce returned error: %s", err)
				return
			}

			// Announce transaction
			_, err = testnetClient.Transaction.Announce(context.Background(), signedTransaction)
			if err != nil {
				fmt.Printf("Transaction.Announce returned error: %s", err)
				return
			}

		default:
			fmt.Println("default")
		}

	}
}

func createCloneAccount(client *sdk.Client, signerInfo *sdk.AccountInfo, balance []*sdk.Mosaic) bool {
	genesis, err := client.NewAccountFromPrivateKey(privateKey)
	if err != nil {
		fmt.Printf("Failed to create genesis: %s", err)
		return false
	}
	clone, err := client.NewAccountFromPrivateKey(signerInfo.PublicKey)
	if err != nil {
		fmt.Printf("Failed to create clone: %s", err)
		return false
	}
	transaction, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour*1),
		sdk.NewAddress(clone.Address.Address, client.NetworkType()),
		balance,
		sdk.NewPlainMessage(""),
	)
	if err != nil {
		fmt.Printf("NewTransferTransaction returned error: %s", err)
		return false
	}

	signedTransaction, err := genesis.Sign(transaction)
	if err != nil {
		fmt.Printf("Sign returned error: %s", err)
		return false
	}
	_, err = client.Transaction.Announce(context.Background(), signedTransaction)
	if err != nil {
		fmt.Printf("Transaction.Announce returned error: %s", err)
		return false
	}
	time.Sleep(10 * time.Second)
	return true
}

func pairMosaics(mosaics []*sdk.Mosaic) []*sdk.Mosaic {
	// Mainnet asset
	xpx, _ := sdk.NewMosaicId(xpx)
	so, _ := sdk.NewMosaicId(so)
	sm, _ := sdk.NewMosaicId(sm)
	sc, _ := sdk.NewMosaicId(sc)
	metx, _ := sdk.NewMosaicId(metx)

	// Testnet asset
	xpx2, _ := sdk.NewMosaicId(xpx2)
	so2, _ := sdk.NewMosaicId(so2)
	sm2, _ := sdk.NewMosaicId(sm2)
	sc2, _ := sdk.NewMosaicId(sc2)
	// metx2, _ := sdk.NewMosaicId(metx2)

	var updatedMosaics []*sdk.Mosaic
	for _, m := range mosaics {

		mid := m.AssetId.String()

		switch mid {
		case xpx.String(), sdk.XpxNamespaceId.String():
			m.AssetId = xpx2
			updatedMosaics = append(updatedMosaics, m)
		case so.String():
			m.AssetId = so2
		case sm.String():
			m.AssetId = sm2
		case sc.String():
			m.AssetId = sc2
		case metx.String():
			fmt.Println("skip metx") //todo
			
		default:
			continue
		}
	}
	return updatedMosaics
}

func restoreMosaics(mosaics []*sdk.Mosaic, txn *sdk.TransferTransaction) []*sdk.Mosaic {
	fmt.Println("restoring mosaic")
	for _, transfered := range txn.Mosaics {
		for _, balance := range mosaics {
			if transfered.AssetId.String() == balance.AssetId.String() {
				balance.Amount = balance.Amount + transfered.Amount
			}
		}
	}
	return mosaics
}