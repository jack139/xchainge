##  应用层API

###  一、 说明

​		应用层API与区块链节点一起部署，提供给客户端调用，进行基础的区块链功能操作。



### 二、 概念和定义

#### 1. 节点

​		节点是区块链上的一个业务处理和存储的单元，是一个具有独立处理区块链业务的服务程序。节点可以是一台物理服务器，也可以是多个节点公用一个物理服务器，通过不同端口提供各自节点的功能。

#### 2. 链用户

​		链用户是具有在提交区块链交易的用户，线下可定义为交易所。每个链用户通过一对密钥识别，同时使用此密钥进行数据的加密解密操作，因此链用户的密钥需要妥善保管。密钥类似如下形式：
```json
{
	"sign_key":{
		"type":"ed25519/privkey",
		"value":"UgM13IPx/BkwfQo8jceLq1CiXlT3lm4WLZ6K6TMR5bRueBGgTDVAv7ZLdBooTZWm2ixLaNitCW91NHW06h8VQw=="
	},
	"CryptoPair":{
		"PrivKey":"tgNfUoYkh9xKs1hVKs+5uXNetCxvDRRHBNmLMs5/NKk=",
		"PubKey":"qyBsXnVKKjvFNxHBRudc3tCp8t8ymqBSF1Ga8qlfqFs="
	}
}
```

#### 3. 交易区块
​		链上数据存储在区块中，区块目前分两类：（1）买卖交易区块；（2）数据授权区块。买卖交易用于存储买入卖出交易的交易信息和交易数据。交易区块中的部分数据是公开的，部分数据是加密的。链用户只能查看自己提交的区块的加密数据。如果要查看其他链用户的加密数据，需要使用向交易提交者进行请求授权。当交易提交者同意并授权后，请求方才能看到相应区块的加密数据。请求和授权过程也会记录在区块链上，用于追溯。

**交易区块内容：**

| 名称       | 类型   | 说明                                       |
| ---------- | ------ | ------------------------------------------ |
| ID         | uuid   | 交易ID，自动生成                           |
| ExchangeID | string | 交易所ID（即，链用户公钥）                 |
| AssetsID   | string | 资产ID，唯一标示交易资产，由客户端定义     |
| Data       | string | 加密交易数据（只有链用户ExchangeID可解密） |
| Refer      | string | 参考数据，可用于检索                       |
| Action     | byte   | 交易类型：1 买入， 2 卖出， 3 变更所有权   |

**授权区块内容：**

| 名称           | 类型   | 说明                                   |
| -------------- | ------ | -------------------------------------- |
| ID             | uuid   | 授权ID，自动生成                       |
| ExchangeID     | string | 数据原始提交者的交易所ID（链用户公钥） |
| AuthExchangeID | string | 请求授权的交易所ID（链用户公钥）       |
| Data           | string | 加密交易数据（AuthExchangeID可以解密） |
| Action         | byte   | 交易类型：4 请求授权， 5 响应授权      |

> 说明：
>
> 1. 上述字段中，AssetsID、Data、Refer均没有长度限制，但不建议放很大的数据块
> 2. 如果需要存储大型数据，请使用IPFS存储，然后在Data字段保存IPFS的文件哈希值
> 3. AssetsID必须是可显示字符（32<ASCII<127）



### 三、 API提供的区块链功能

| 序号 | 接口名称        | 接口功能                                     |
| :--: | :-------------- | -------------------------------------------- |
|  1   | deal            | 提交买卖交易                                 |
|  2   | auth_request    | 请求授权查询指定交易                         |
|  3   | auth_response   | 授权查看指定交易                             |
|  4   | query_deals     | 查询自己的所以交易                           |
|  5   | query_auths     | 查询授权请求和授权响应                       |
|  6   | query_by_assets | 按资产ID进行检索（可能包括其他链用户的区块） |
|  7   | query_by_refer  | 按参考值进行检索                             |
|  8   | query_block     | 按区块ID查询制定区块                         |
|  9   | ipfs_upload     | 数据上传到ipfs                               |
|  10  | ipfs_download   | 从ipfs下载数据                               |




### 四、接口定义

#### 1. 全局接口定义

输入参数

| 参数      | 类型   | 说明                          | 示例        |
| --------- | ------ | ----------------------------- | ----------- |
| appId     | string | 应用渠道编号                  |             |
| version   | string | 版本号                        |             |
| signType  | string | 签名算法，目前使用SHA256算法 | SHA256或SM2 |
| signData  | string | 签名数据，具体算法见下文      |             |
| encType   | string | 接口数据加密算法，目前不加密  | plain       |
| timestamp | int    | unix时间戳（秒）              |             |
| data      | json   | 接口数据，详见各接口定义      |             |

> 签名/验签算法：
>
> 1. 筛选，获取参数键值对，剔除signData、encData、extra三个参数。data参数按key升序排列进行json序列化。
> 2. 排序，按key升序排序。
> 3. 拼接，按排序好的顺序拼接请求参数
>
> ```key1=value1&key2=value2&...&key=appSecret```，key=appSecret固定拼接在参数串末尾，appSecret需替换成应用渠道所分配的appSecret。
>
> 4. 签名，使用制定的算法进行加签获取二进制字节，使用 16进制进行编码Hex.encode得到签名串，然后base64编码。
> 5. 验签，对收到的参数按1-4步骤签名，比对得到的签名串与提交的签名串是否一致。

签名示例：

```json
请求参数：
{
    "appid":"19E179E5DC29C05E65B90CDE57A1C7E5",
    "version": "1",
    "signType": "SHA256",
    "signData": "...",
    "encType": "plain",
    "timestamp":1591943910,
    "data": {
    	"user_id":"gt",
    	"face_id":"5ed21b1c262daabe314048f5"
    }
}

密钥：
appSecret="D91CEB11EE62219CD91CEB11EE62219C"

待加签串：
appid=19E179E5DC29C05E65B90CDE57A1C7E5&data={"user_id":"gt","face_id":"5ed21b1c262daabe314048f5"}&encType=plain&signType=SM2&timestamp=1591943910&version=1&key=D91CEB11EE62219CD91CEB11EE62219C

SHA256加签结果：
"10e13147546debbea157ec793170968c6c614f4eb13ccd9b7a9c193bf1b3bd78"

base64后结果：
"MTBlMTMxNDc1NDZkZWJiZWExNTdlYzc5MzE3MDk2OGM2YzYxNGY0ZWIxM2NjZDliN2E5YzE5M2JmMWIzYmQ3OA=="
```

返回结果

| 参数      | 类型    | 说明                                                         |
| --------- | ------- | ------------------------------------------------------------ |
| appId     | string  | 应用渠道编号                                                 |
| code      | string  | 返回状态代码                                                 |
| encType   | string  | 数据加密算法，目前不加密                                     |
| success   | boolean | 成功与否                                                     |
| timestamp | int     | unix时间戳                                                   |
| data      | json    | 成功时返回结果数据；出错时，data.msg返回错误说明。详见具体接口 |

> 成功时：code为0， success为True，data内容见各接口定义；
>
> 出错时：code返回错误代码，具体定义见各接口说明

返回示例

```json
{
    "appId": "19E179E5DC29C05E65B90CDE57A1C7E5", 
    "code": 0, 
    "encType": "plain",
    "success": true,
    "timestamp": 1591943910,
    "data": {
       "msg": "success", 
       ...
    }
}
```

全局出错代码

| 编码 | 说明                               |
| ---- | ---------------------------------- |
| 9800 | 无效签名                           |
| 9801 | 签名参数有错误                     |
| 9802 | 调用时间错误，unixtime超出接受范围 |



#### 2. 交易提交接口



#### 3. 请求授权接口



#### 4. 响应授权接口



#### 5. 查询接口



#### 6. IPFS接口