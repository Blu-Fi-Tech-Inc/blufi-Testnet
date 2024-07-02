package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/blu-fi-tech-inc/boriqua_project/core"
	"github.com/blu-fi-tech-inc/boriqua_project/crypto"
	"github.com/blu-fi-tech-inc/boriqua_project/network"
	"github.com/blu-fi-tech-inc/boriqua_project/types"
	"github.com/blu-fi-tech-inc/boriqua_project/utils"
)

func main() {
	validatorPrivKey := crypto.GeneratePrivateKey()

	localNode := makeServer("LOCAL_NODE", &validatorPrivKey, ":3000", []string{":4000"}, ":9000")
	go localNode.Start()

	remoteNode := makeServer("REMOTE_NODE", nil, ":4000", []string{":5000"}, "")
	go remoteNode.Start()

	remoteNodeB := makeServer("REMOTE_NODE_B", nil, ":5000", nil, "")
	go remoteNodeB.Start()

	go func() {
		time.Sleep(11 * time.Second)
		lateNode := makeServer("LATE_NODE", nil, ":6000", []string{":4000"}, "")
		go lateNode.Start()
	}()

	// Keep the main function running
	select {}
}

// sendTransaction sends a transaction from one node to another
func sendTransaction(privKey crypto.PrivateKey) error {
	toPrivKey := crypto.GeneratePrivateKey()
	tx := core.NewTransaction(nil)
	tx.To = toPrivKey.PublicKey()
	tx.Value = 666

	if err := tx.Sign(privKey); err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "http://localhost:9000/tx", buf)
	if err != nil {
		return err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// makeServer creates and initializes a server with the given options
func makeServer(id string, pk *crypto.PrivateKey, addr string, seedNodes []string, apiListenAddr string) *network.Server {
	opts := network.ServerOpts{
		APIListenAddr: apiListenAddr,
		SeedNodes:     seedNodes,
		ListenAddr:    addr,
		PrivateKey:    pk,
		ID:            id,
	}

	s, err := network.NewServer(opts)
	if err != nil {
		log.Fatalf("Failed to create server %s: %v", id, err)
	}

	return s
}

// createCollectionTx creates a collection transaction and sends it to the network
func createCollectionTx(privKey crypto.PrivateKey) types.Hash {
	tx := core.NewTransaction(nil)
	tx.TxInner = core.CollectionTx{
		Fee:      200,
		MetaData: []byte("chicken and egg collection!"),
	}
	if err := tx.Sign(privKey); err != nil {
		log.Fatalf("Failed to sign collection transaction: %v", err)
	}

	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		log.Fatalf("Failed to encode collection transaction: %v", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:9000/tx", buf)
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	return tx.Hash(core.TxHasher{})
}

// nftMinter mints an NFT and sends the transaction to the network
func nftMinter(privKey crypto.PrivateKey, collection types.Hash) {
	metaData := map[string]interface{}{
		"power":  8,
		"health": 100,
		"color":  "green",
		"rare":   "yes",
	}

	metaBuf := new(bytes.Buffer)
	if err := json.NewEncoder(metaBuf).Encode(metaData); err != nil {
		log.Fatalf("Failed to encode metadata: %v", err)
	}

	tx := core.NewTransaction(nil)
	tx.TxInner = core.MintTx{
		Fee:             200,
		NFT:             utils.RandomHash(),
		MetaData:        metaBuf.Bytes(),
		Collection:      collection,
		CollectionOwner: privKey.PublicKey(),
	}
	if err := tx.Sign(privKey); err != nil {
		log.Fatalf("Failed to sign mint transaction: %v", err)
	}

	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		log.Fatalf("Failed to encode mint transaction: %v", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:9000/tx", buf)
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}
}
