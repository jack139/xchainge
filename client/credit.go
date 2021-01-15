package client

import (
	"xchainge/types"

	"fmt"
	"strconv"
	"time"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
)

// CR上链，数据加密
func (me *User) Credit(action, data, num string) ([]byte, error) {

	// 用户id
	userId := *me.CryptoPair.PubKey

	now := time.Now()
	tx := new(types.Transx)
	tx.SendTime = &now
	action0, _ := strconv.Atoi(action)
	num0, _ := strconv.Atoi(num)

	// 不是占位块，都只上链一块
	if action0!=1 { 
		num0 = 1
	}

	// 准备数据
	credit := types.Credit{
		//ID:     uuid.NewV4(),
		UserID: userId,
		Data:   []byte(data),
		Action: byte(action0),
	}
	tx.Payload = &credit

	for i:=0;i<num0;i++ {
		// 每次更新ID，ID是不一样的
		credit.ID = uuid.NewV4()

		tx.Sign(me.SignKey)
		tx.SignPubKey = me.SignKey.PubKey()
		// broadcast this tx
		bz, err := cdc.MarshalJSON(&tx)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		// 上链
		ret, err := cli.BroadcastTxSync(ctx, bz)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		fmt.Printf("credit => %+v\n", ret)

		// ret  *ctypes.ResultBroadcastTxCommit
		if ret.Code !=0 {
			fmt.Println(ret.Log)
			return nil, fmt.Errorf(ret.Log)
		}

	}

	//  最后一次的ID
	respMap := map[string]interface{}{"id" : credit.ID.String()} 

	// 返回结果转为json
	respBytes, err := json.Marshal(respMap)
	if err != nil {
		return nil, err
	}

	return respBytes, nil
}
