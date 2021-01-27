## Roadmap



### 功能描述

1. 核心数据上链（交易链+ipfs），非核心数据不上链（保存在DB）
2. 交易链存储基本交易信息、账户信息，数据实体存储在ipfs
3. 数据实体为字节流，内容由应用确定，不做限制
4. 交易数据、授权数据全部上链
5. 以链用户（交易所）区分数据所有权，链用户上传的数据对自己是开放的
6. 交易数据使用链用户密钥加密，其他节点是否可见，需要获取数据所有者的授权
7. 交易节点提供功能：交易上传，交易查询，资产溯源，查询授权请求，授权查询
9. 资产链，在交易链上用链上链表实现
10. 交易区块的交易类型：资产买卖；查询授权


### 交易链请求

1. 交易内容
```json
{
	"Signature":"SXKaVqfAe5ypHpz1qM3tQTYa42F9JQoq4zwMYQLKN0E0s+nViVk2Z3b98mFXvTHnqRCFousPVCYdR7b+d21jCg==",
	"SendTime":"2020-12-18T05:36:00.281914675Z",
	"SignPubKey":{
			"type":"ed25519/pubkey",
			"value":"yaeWs6Y5a0djpvnShNwq+zZdeJmN9I+nrddWMMH+3Uo="
	},
	"Payload":{
		// 详见下
	}
}
```

2. 资产交易 payload
```json
{
	"type":"deal", // 交易
	"value":{
		"ID":"", // 资产交易ID
		"Assets":"K1kY3yTfSwW9lphr5RzjLw==", // 资产ID
		"Exchange":"P1ABCkAph4DQlnEMahGW6I2mfOOtfZKYyssOZ4L8MTc=", // 交易所ID（公钥）
		"Data":"", // 加密交易数据 （IPFS HASH）
		"Refer":"abc", // 用于索引（例如，可以存放第三方用户id）
		"Action":1, // 0x01 买入， 0x02 卖出， 0x03 变更所有权
	}
}
```

3. 查询授权 payload
```json
{
	"type":"auth", // 查询授权（授权其他交易所查看某资产），查询记录（只记录被授权方的查询动作）
	"value":{
		"ID":"", // 授权操作ID
		"DealID":"", // 交易ID
		"FromExchange":"P1ABCkAph4DQlnEMahGW6I2mfOOtfZKYyssOZ4L8MTc=", // 授权交易所ID（公钥）
		"ToExchange":"P1ABCkAph4DQlnEMahGW6I2mfOOtfZKYyssOZ4L8MTc=", // 被授权交易所ID（公钥）
		"Data":"", // FromExchange加密数据，被授权者ToExchangeID可以解密
		"Action":4, // 0x04 请求授权， 0x05 响应授权
	}
}
```



### 区块例子

```json
{
	"header":{
		"version":{
			"block":"11",
			"app":"1"
		},
		"chain_id":"test-chain-FEeTGF",
		"height":"8",
		"time":"2020-12-24T05:24:01.760181367Z",
		"last_block_id":{
			"hash":"3D326CA03E1D0E6D9C80FB6B788AD1A72BB12E10B9DE617B13AF311E5258ABA8",
			"parts":{
				"total":1,
				"hash":"0A675856141DB019BB245E87896DEA5A3C7BAE8CC2C3C1A7666DF96236B802CB"
			}
		},
		"last_commit_hash":"63AA80A1CE7261E001383A3754DAF09C52DA25881EFDD6E3E1F1541A937C5AAE",
		"data_hash":"D14D445C897F2E7A2518FB1EEE69A969F178F6D819ADEBB8D91B5B528CCC01C7",
		"validators_hash":"82F872A1F21F7C05578D5397DA499A2D656E61D0DD7F9EFE6531F433F72306EB",
		"next_validators_hash":"82F872A1F21F7C05578D5397DA499A2D656E61D0DD7F9EFE6531F433F72306EB",
		"consensus_hash":"048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F",
		"app_hash":"0000000000000000",
		"last_results_hash":"6E340B9CFFB37A989CA544E6BB780A2C78901D3FB33738768511A30617AFA01D",
		"evidence_hash":"E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855",
		"proposer_address":"89581502243C3D3401C38EEF4C1A1145AB514B11"
	},
	"data":{
		"txs":[
			"eyJTaWduYXR1cmUiOiIwbDFvTFhoayt4YXF2UC96MXJMZ1U5Mzh3d01wZzExYlMwcko3V1lrNmNlNzg1ZFBiQWN0dE9sRzZ5UkJXQ1RlU09kZnh2TlBScDN3a0JFbVVOaGhCdz09IiwiU2VuZFRpbWUiOiIyMDIwLTEyLTI0VDA1OjU3OjQ0LjUzNzg1NDA3OFoiLCJTaWduUHViS2V5Ijp7InR5cGUiOiJlZDI1NTE5L3B1YmtleSIsInZhbHVlIjoiVXJobzU5a3UrOEF1Mit6aS9lTDVNKzB5SjhUQVdYMDJkY005K2pTZDE5ND0ifSwiUGF5bG9hZCI6eyJ0eXBlIjoiYXV0aCIsInZhbHVlIjp7IklEIjoiSnU2bXFla2pSRmE1dG1Ob0RnR0R2UT09IiwiRGVhbElEIjoiWHVFSkRhckxUM3lkU2Y0YkY2ZUVZZz09IiwiRnJvbUV4Y2hhbmdlSUQiOiJxeUJzWG5WS0tqdkZOeEhCUnVkYzN0Q3A4dDh5bXFCU0YxR2E4cWxmcUZzPSIsIlRvRXhjaGFuZ2VJRCI6Imo5Y0lnbW0xN3gwYUxBcGYwaTIwVVI3UGozNFVhL0p3eVdPdUJHZ1lJRmc9IiwiRGF0YSI6bnVsbCwiQWN0aW9uIjo0fX19"
		]
	},
	"evidence":{
		"evidence":[]
	},
	"last_commit":{
		"height":"7",
		"round":0,
		"block_id":{
			"hash":"3D326CA03E1D0E6D9C80FB6B788AD1A72BB12E10B9DE617B13AF311E5258ABA8",
			"parts":{
				"total":1,
				"hash":"0A675856141DB019BB245E87896DEA5A3C7BAE8CC2C3C1A7666DF96236B802CB"
			}
		},
		"signatures":[
			{
				"block_id_flag":2,
				"validator_address":"89581502243C3D3401C38EEF4C1A1145AB514B11",
				"timestamp":"2020-12-24T05:24:01.760181367Z",
				"signature":"ZfgneVPY/pOEjygwmEQnMIu4iQT8QgRf/AdjHptkbpqT57dCMFa4V+7bxAKIzoCUcJgtLtrg1bJtJdXQNk7gCA=="
			}
		]
	}
}
```



### leveldb 逻辑分表

| 前缀       | key             | value     |
| ---------- | --------------- | --------- |
| blockLink: | 区块高度        | 区块高度  |
| assetsLink:  | 资产id        | 区块高度  |


### IPFS

启动
```
nohup ipfs daemon --enable-namesys-pubsub > /tmp/ipfs.log 2>&1 &
```

### 技术栈

1. 区块链 Tendermint 0.34.0
2. 分布式存储 IPFS 0.7.0
3. 节点数据库 LevelDB 1.20
4. 开发语言 Go 1.15.6



### TODO

- [x] 交易上链
- [x] 链上查询
- [x] 链上支持多链表
- [x] 用户认证，使用ed25519签名
- [x] 数据使用curve25519加密
- [x] 交易查询授权和响应
- [x] ipfs支持
- [x] 应用层api
- [ ] credit产生机制，产生方案记录在链上（app states），通过区块修改
- [ ] 增加区块类型：系统区块，credit交易区块
- [x] 调整验签机制
- [x] 业务处理api(注册、合同、验收)
