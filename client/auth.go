package client

import (
	"xchainge/types"

	"fmt"
	"io"
	"time"
	crypto_rand "crypto/rand"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/nacl/box"
)


// 请求授权 上链
// xcli authRequest j9cIgmm17x0aLApf0i20UR7Pj34Ua/JwyWOuBGgYIFg= dcfe656c-6c65-45e7-9e94-f082a068a93d
func (me *User) AuthRequest(fromExchangeId, dealId string) error {
	now := time.Now()

	// 检查 toExchangeId 合理性
	var fromExchangeIdBytes [32]byte
	err := cdc.UnmarshalJSON([]byte("\""+fromExchangeId+"\""), &fromExchangeIdBytes) // 反序列化时需要双引号，因为是字符串
	if err != nil {
		return err
	}

	// 检查 dealID -->  UUID
	uuidDealId, err := uuid.FromString(dealId)
	if err != nil {
		return err
	}

	// 新建交易
	tx := new(types.Transx)
	tx.SendTime = &now

	auth := new(types.Auth)
	auth.ID = uuid.NewV4()
	auth.DealID = uuidDealId
	auth.FromExchangeID = fromExchangeIdBytes
	auth.ToExchangeID = *me.CryptoPair.PubKey
	auth.Action = 0x04

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
	fmt.Printf("auth request => %+v\n", ret)
	return nil
}


// 响应授权 上链
// xcli authRequest dcfe656c-6c65-45e7-9e94-f082a068a93d
func (me *User) AuthResponse(authId string) error {
	addr, _ := cdc.MarshalJSON(*me.CryptoPair.PubKey)

	now := time.Now()

	// 获取 authID 对应的 授权请求 块
	authTx, err := queryTx(addr, "_", authId)
	if err != nil {
		return err
	}
	if authTx==nil {
		return fmt.Errorf("AuthID not found")
	}
	auth, ok := (*authTx).Payload.(*types.Auth)	// 授权块
	if !ok {
		return fmt.Errorf("need a Auth Payload")
	}

	// 检查是否已响应过，在toExchangeID的列表里找
	toExchangeId, _ := cdc.MarshalJSON(auth.ToExchangeID)
	toAuthList, err := query(addr, "auth", string(toExchangeId[1 : len(toExchangeId)-1]))
	if err != nil {
		return err
	}

	isAuthorised := false
	for _, tx := range *toAuthList {
		authItem, ok := tx.Payload.(*types.Auth) // 交易
		if ok {
			if authItem.Action==0x05 && authItem.DealID==auth.DealID {
				isAuthorised = true
				break
			}
		}
	}

	if isAuthorised { // 已经授权过
		return fmt.Errorf("Authorized")
	}

	// 获取 authID 对应的 交易块
	dealTx, err := queryTx(addr, "_", auth.DealID.String())
	if err != nil {
		return err
	}

	if dealTx==nil {
		return fmt.Errorf("DealID not found")
	}

	deal, ok := (*dealTx).Payload.(*types.Deal)	// 交易块
	if !ok {
		return fmt.Errorf("need a Deal Payload")
	}

	// 解密
	var decryptKey, publicKey [32]byte

	publicKey = auth.FromExchangeID

	// 解密 data 数据
	box.Precompute(&decryptKey, &publicKey, me.CryptoPair.PrivKey)
	var decryptNonce [24]byte
	copy(decryptNonce[:], deal.Data[:24])
	//fmt.Printf("data=>%v,decryptNonce=>%v,decryptKey=>%v\n", deal.Data[24:], decryptNonce, decryptKey)
	decrypted, ok := box.OpenAfterPrecomputation(nil, deal.Data[24:], &decryptNonce, &decryptKey)
	if !ok {
		return fmt.Errorf("decryption error")
	}

	// 重新加密，使用toExchangeID
	publicKey = auth.ToExchangeID

	sharedEncryptKey := new([32]byte)
	box.Precompute(sharedEncryptKey, &publicKey, me.CryptoPair.PrivKey)

	var nonce [24]byte
	if _, err := io.ReadFull(crypto_rand.Reader, nonce[:]); err != nil {
		panic(err)
	}
	//fmt.Printf("data=>%v,nonce=>%v,sharedEncryptKey=>%v\n", decrypted, nonce, *sharedEncryptKey)
	encrypted := box.SealAfterPrecomputation(nonce[:], decrypted, &nonce, sharedEncryptKey)

	// 新建交易
	tx := new(types.Transx)
	tx.SendTime = &now

	authResp := new(types.Auth)
	authResp.ID = uuid.NewV4()
	authResp.DealID = auth.DealID
	authResp.FromExchangeID = auth.FromExchangeID
	authResp.ToExchangeID = auth.ToExchangeID
	authResp.Data = encrypted
	authResp.Action = 0x05

	tx.Payload = authResp

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
	fmt.Printf("auth respose => %+v\n", ret)
	return nil
}
