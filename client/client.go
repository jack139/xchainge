package client

import (
	"xchainge/types"

	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"time"
	"context"
	crypto_rand "crypto/rand"
	"encoding/json"
	"encoding/base64"

	cfg "github.com/tendermint/tendermint/config"
	cmn "github.com/tendermint/tendermint/libs/os"
	rpcclient "github.com/tendermint/tendermint/rpc/client/http"
	uuid "github.com/satori/go.uuid"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"golang.org/x/crypto/nacl/box"
)

// KEYFILENAME 私钥文件名
const KEYFILENAME string = "exchange.key"

var (
	cli *rpcclient.HTTP
	cdc = types.AminoCdc
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

type User struct {
	SignKey    crypto.PrivKey `json:"sign_key"` // 节点私钥，用户签名
	CryptoPair cryptoPair     // 密钥协商使用
}


// 从文件装入key
func LoadOrGenUserKey() (*User, error) {
	if cmn.FileExists(KEYFILENAME) {
		uk, err := loadUserKey()
		if err != nil {
			return nil, err
		}
		return uk, nil
	}
	//fmt.Println("userkey file not exists")
	uk := new(User)
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

func loadUserKey() (*User, error) {
	//copy(privKey[:], bz)
	jsonBytes, err := ioutil.ReadFile(KEYFILENAME)
	if err != nil {
		return nil, err
	}
	uk := new(User)
	err = cdc.UnmarshalJSON(jsonBytes, uk)
	if err != nil {
		return nil, fmt.Errorf("Error reading UserKey from %v: %v", KEYFILENAME, err)
	}
	return uk, nil
}


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

// 授权操作 上链
// xcli auth 4 234 j9cIgmm17x0aLApf0i20UR7Pj34Ua/JwyWOuBGgYIFg= yyy
func (me *User) Auth(action, assetsId, toExchangeId, refer string) error {
	now := time.Now()

	// 检查 toExchangeId 合理性
	var toExchangeIdBytes [32]byte
	err := cdc.UnmarshalJSON([]byte("\""+toExchangeId+"\""), &toExchangeIdBytes) // 反序列化时需要双引号，因为是字符串
	if err != nil {
		return err
	}

	// 新建交易
	tx := new(types.Transx)
	tx.SendTime = &now

	auth := new(types.Auth)
	auth.ID = uuid.NewV4()
	auth.AssetsID = []byte(assetsId)
	auth.FromExchangeID = *me.CryptoPair.PubKey
	auth.ToExchangeID = toExchangeIdBytes
	auth.Refer = []byte(refer)
	action0, _ := strconv.Atoi(action)
	auth.Action = byte(action0)

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

// 链上查询  category取值： exchange, assets, refer
// xcli queryExchange _
// xcli queryExchange j9cIgmm17x0aLApf0i20UR7Pj34Ua/JwyWOuBGgYIFg=
func (me *User) Query(category, queryContent string) error {
	addr, _ := cdc.MarshalJSON(*me.CryptoPair.PubKey)

	//addr = addr[1 : len(addr)-1] // 移除两边的引号
	var buf bytes.Buffer
	buf.WriteString("/")
	buf.Write(addr)
	buf.WriteString("/query/")
	buf.WriteString(category)
	//获得拼接后的字符串
	path := buf.String()
	if category=="exchange" && queryContent!="_" {  // 用户公钥需要加双引号
		queryContent = "\"" + queryContent + "\""	
	}
	rsp, err := cli.ABCIQuery(ctx, path, []byte(queryContent))
	if err != nil {
		fmt.Println(err)
		return err
	}

	data := rsp.Response.Value
	//fmt.Printf("resp => %s\n", data)

	/*
		exchange 不解密
		assets 根据授权解密
		refer 不解密
	*/

	var txHistory []types.Transx
	var respList []types.RespQuery
	cdc.UnmarshalJSON(data, &txHistory)

	for _, tx := range txHistory {
		deal, ok := tx.Payload.(*types.Deal) // 交易
		if ok {
			//fmt.Printf("deal => %v\n", deal)

			// data 默认返回加密数据的 base64
			data := base64.StdEncoding.EncodeToString(deal.Data)
			
			// 如果查询 assets，则尝试解密 data
			if category=="assets" {
				var decryptKey, publicKey [32]byte

				if deal.ExchangeID==*me.CryptoPair.PubKey { // 是自己的交易, 进行解密
					publicKey = deal.ExchangeID

					// 解密 data 数据
					box.Precompute(&decryptKey, &publicKey, me.CryptoPair.PrivKey)
					var decryptNonce [24]byte
					copy(decryptNonce[:], deal.Data[:24])
					//fmt.Printf("data=>%v,decryptNonce=>%v,decryptKey=>%v\n", deal.Data[24:], decryptNonce, decryptKey)
					decrypted, ok := box.OpenAfterPrecomputation(nil, deal.Data[24:], &decryptNonce, &decryptKey)
					if !ok {
						return fmt.Errorf("decryption error")
					}
					data = string(decrypted)
				}
			}

			exchangeId, _ := cdc.MarshalJSON(deal.ExchangeID)
			respList = append(respList, types.RespQuery{
				Type          : "DEAL",
				ID            : deal.ID.String(),
				ExchangeId    : string(exchangeId[1 : len(exchangeId)-1]), // 去掉两边引号
				AssetsId      : string(deal.AssetsID),
				Data          : data,
				Refer         : string(deal.Refer),
				Action        : deal.Action,
			})
		} else {
			auth, ok := tx.Payload.(*types.Auth)	// 授权
			if ok {
				//fmt.Printf("auth => %v\n", auth)
				exchangeId, _ := cdc.MarshalJSON(auth.FromExchangeID)
				exchangeId2, _ := cdc.MarshalJSON(auth.ToExchangeID)
				respList = append(respList, types.RespQuery{
					Type           : "AUTH",
					ID             : auth.ID.String(),
					ExchangeId     : string(exchangeId[1 : len(exchangeId)-1]), // 去掉两边引号
					AuthExchangeId : string(exchangeId2[1 : len(exchangeId2)-1]),
					AssetsId       : string(auth.AssetsID),
					Refer          : string(auth.Refer),
					Action         : auth.Action,
				})
			}
		}
	}

	// 返回结果转为json
	respBytes, err := json.Marshal(respList)
	if err != nil {
		return err
	}

	fmt.Printf("json => %s\n", respBytes)

	return nil
}

