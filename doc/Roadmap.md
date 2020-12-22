## Roadmap

### 功能描述

1. 核心数据上链（交易链+ipfs），非核心数据不上链（保存在DB）
2. 交易链存储基本交易信息、账户信息，数据实体存储在ipfs
3. 数据实体为字节流，内容由应用确定，不做限制
4. 交易数据、授权数据全部上链
5. 以交易节点（交易所）区分交易数据所有权，交易节点上传的数据对本节点是开放的
6. 交易数据由交易节点加密，其他节点是否可见，需要获取上传节点的授权
7. 交易节点提供功能：交易上传，交易查询，资产溯源，交易查询授权
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



### leveldb 逻辑分表

| 前缀       | key             | value     |
| ---------- | --------------- | --------- |
| blockLink: | 区块高度        | 区块高度  |
| assetsLink:  | 资产id        | 区块高度  |



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
- [ ] 交易查询授权和响应
- [ ] ipfs支持
- [ ] 应用层api
