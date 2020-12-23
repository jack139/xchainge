package client

import (
	"xchainge/types"

	"fmt"
	"io"
	"strconv"
	"time"
	crypto_rand "crypto/rand"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/nacl/box"
)



// 交易上链，数据加密
func (me *User) Deal(action, assetsId, data, refer string) error {
	// 交易所id
	exchangeId := *me.CryptoPair.PubKey

	sharedEncryptKey := new([32]byte)
	box.Precompute(sharedEncryptKey, &exchangeId, me.CryptoPair.PrivKey)

	var nonce [24]byte
	if _, err := io.ReadFull(crypto_rand.Reader, nonce[:]); err != nil {
		panic(err)
	}
	//fmt.Printf("data=>%v,nonce=>%v,sharedEncryptKey=>%v\n", data, nonce, *sharedEncryptKey)
	encrypted := box.SealAfterPrecomputation(nonce[:], []byte(data), &nonce, sharedEncryptKey)

	now := time.Now()
	tx := new(types.Transx)
	tx.SendTime = &now
	action0, _ := strconv.Atoi(action)

	deal := types.Deal{
		ID:         uuid.NewV4(),
		AssetsID:   []byte(assetsId),
		ExchangeID: exchangeId,
		Data:       encrypted,
		Refer:      []byte(refer),
		Action:     byte(action0),
	}

	tx.Payload = &deal

	tx.Sign(me.SignKey)
	tx.SignPubKey = me.SignKey.PubKey()
	// broadcast this tx
	bz, err := cdc.MarshalJSON(&tx)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ret, err := cli.BroadcastTxAsync(ctx, bz)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("deal => %+v\n", ret)
	return nil
}
