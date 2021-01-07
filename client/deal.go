package client

import (
	"xchainge/types"

	"fmt"
	"io"
	"strconv"
	"time"
	"encoding/json"
	crypto_rand "crypto/rand"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/nacl/box"
)

func isASCII(s string) bool {
    for i := 0; i < len(s); i++ {
        if s[i] <= 32 || s[i] >= 127 {
            return false
        }
    }
    return true
}

// 交易上链，数据加密
func (me *User) Deal(action, assetsId, data, refer string) ([]byte, error) {

	if !isASCII(assetsId) {
		return nil, fmt.Errorf("assetsId must be visible ASCII")
	}

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
		return nil, err
	}

	ret, err := cli.BroadcastTxSync(ctx, bz)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Printf("deal => %+v\n", ret)

	// ret  *ctypes.ResultBroadcastTxCommit
	if ret.Code !=0 {
		fmt.Println(ret.Log)
		return nil, fmt.Errorf(ret.Log)
	}

	respMap := map[string]interface{}{"id" : deal.ID.String()}

	// 返回结果转为json
	respBytes, err := json.Marshal(respMap)
	if err != nil {
		return nil, err
	}

	return respBytes, nil
}
