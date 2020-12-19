package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	//"strconv"
	"time"
	"context"
	"xchainge"
	"xchainge/types"

	cfg "github.com/tendermint/tendermint/config"
	cmn "github.com/tendermint/tendermint/libs/os"
	rpcclient "github.com/tendermint/tendermint/rpc/client/http"

	crypto_rand "crypto/rand"

	uuid "github.com/satori/go.uuid"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"golang.org/x/crypto/nacl/box"
)

// KEYFILENAME 私钥文件名
const KEYFILENAME string = "exchange.key"

var (
	cli *rpcclient.HTTP
	cdc = xchainge.AminoCdc
	ctx = context.Background()
)

func init() {
	addr := cfg.DefaultRPCConfig().ListenAddress
	cli, _ = rpcclient.New(addr, "/websocket")
}

type cryptoPair struct {
	PrivKey *[32]byte
	PubKey  *[32]byte
}

type user struct {
	SignKey    crypto.PrivKey `json:"sign_key"` // 节点私钥，用户签名
	CryptoPair cryptoPair     // 密钥协商使用
}


// 从文件装入key
func loadOrGenUserKey() (*user, error) {
	if cmn.FileExists(KEYFILENAME) {
		uk, err := loadUserKey()
		if err != nil {
			return nil, err
		}
		return uk, nil
	}
	//fmt.Println("userkey file not exists")
	uk := new(user)
	uk.SignKey = ed25519.GenPrivKey()
	pubKey, priKey, err := box.GenerateKey(crypto_rand.Reader)
	if err != nil {
		return nil, err
	}
	uk.CryptoPair = cryptoPair{PrivKey: priKey, PubKey: pubKey}
	jsonBytes, err := cdc.MarshalJSON(uk)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(KEYFILENAME, jsonBytes, 0644)
	if err != nil {
		return nil, err
	}
	return uk, nil
}

func loadUserKey() (*user, error) {
	//copy(privKey[:], bz)
	jsonBytes, err := ioutil.ReadFile(KEYFILENAME)
	if err != nil {
		return nil, err
	}
	uk := new(user)
	err = cdc.UnmarshalJSON(jsonBytes, uk)
	if err != nil {
		return nil, fmt.Errorf("Error reading UserKey from %v: %v", KEYFILENAME, err)
	}
	return uk, nil
}


// 交易上链，数据加密
func (me *user) deal(action, assetsId, data, refer string) {
	// 交易所id
	exchangeId := *me.CryptoPair.PubKey

	sharedEncryptKey := new([32]byte)
	box.Precompute(sharedEncryptKey, &exchangeId, me.CryptoPair.PrivKey)

	var nonce [24]byte
	if _, err := io.ReadFull(crypto_rand.Reader, nonce[:]); err != nil {
		panic(err)
	}
	//fmt.Printf("msg=>%v,nonce=>%v,sharedEncryptKey=>%v\n", msg, nonce, *sharedEncryptKey)
	encrypted := box.SealAfterPrecomputation(nonce[:], []byte(data), &nonce, sharedEncryptKey)

	now := time.Now()
	tx := new(types.Transx)
	tx.SendTime = &now

	deal := types.Deal{
		ID:         uuid.NewV4(),
		AssetsID:   []byte(assetsId),
		ExchangeID: exchangeId,
		Data:       encrypted,
		Refer:      []byte(refer),
		Action:     byte(action),
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

// 授权操作 上链
func (me *user) auth(action, assetsId, toExchangeId, refer string) error {
	now := time.Now()
	tx := new(types.Transx)
	tx.SendTime = &now

	auth := new(types.Auth)
	auth.ID = uuid.NewV4()
	auth.AssetsID = []byte(assetsId)
	auth.FromExchangeID = *me.CryptoPair.PubKey
	auth.ToExchangeID = [32]byte(toExchangeId)
	auth.Refer = []byte(refer)
	auth.Action = byte(action)

	tx.Payload = auth

	tx.Sign(me.SignKey)
	tx.SignPubKey = me.SignKey.PubKey()

	bz, err := cdc.MarshalJSON(&tx)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ret, err := cli.BroadcastTxSync(ctx, bz)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("auth => %+v\n", ret)
	return nil
}

// 链上查询
func (me *user) query(category, queryContent string) error {
	addr, _ := cdc.MarshalJSON(*me.CryptoPair.PubKey)

	//addr = addr[1 : len(addr)-1] // 移除两边的引号
	var buf bytes.Buffer
	buf.WriteString("/")
	buf.Write(addr)
	buf.WriteString("/query")
	buf.WriteString("/")
	buf.WriteString(category)  // TODO: 检查category合法值
	//获得拼接后的字符串
	path := buf.String()
	rsp, err := cli.ABCIQuery(ctx, path, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	data := rsp.Response.Value
	fmt.Println(data)

	return nil
}

