# RPC Provider Guideline API

* [版本](#版本)
* [高度](#高度)
* [账户余额](#账户余额)
* [状态信息](#状态信息)
* [设置验证节点](#设置验证节点)
* [查看验证节点](#查看验证节点)
* [Custodial抵押](#Custodial抵押)
* [Non-Custodial抵押](#Non-Custodial抵押)


## 版本

描述: 检查运行中 pocket 的当前版本

```shell
curl --url http://127.0.0.1:8088/poktsrv/version 
```

输出:

```shell
{
  "status": "success",
  "code": "",
  "data": {
    "version": "RC-0.9.2"
  }
}
```


## 高度

描述: 检查运行中 pocket 的当前区块高度

```shell
curl --url http://127.0.0.1:8088/poktsrv/height
```

输出:

```shell
{
  "status": "success",
  "code": "",
  "data": {
    "height": 99156
  }
}
```


## 账户余额

描述:查询指定账号的余额

```shell
curl --request POST --url http://127.0.0.1:8088/poktsrv/balance --header 'Content-Type: application/json' \
--data "{\"height\": 0,\"address\":\"ee60841d9afb70ba893c02965537bc0eec4ef1e4\"}"
```

参数：

- **height：** 要查询的区块的指定高度。默认为0，这将查询当前节点已知的最新块。
- **address：** 目标地址。


输出:

```shell
{
  "status": "success",
  "code": "",
  "data": {
    "height": 0,
    "address": "ee60841d9afb70ba893c02965537bc0eec4ef1e4",
    "balance": "39999930000"
  }
}
```


## 状态信息

描述:检查运行中 pocket 的状态信息

```shell
curl --url http://127.0.0.1:8088/poktsrv/status
```

输出:

```shell
{
  "status": "success",
  "code": "",
  "data": {
    "version": "RC-0.9.2",
    "height": 99167,
    "Synced": true,
    "address": "ee60841d9afb70ba893c02965537bc0eec4ef1e4",
    "publicKey": "7b1739685dcdc10fcc02bc21dd822ef3458fcf543cc89487af9fe512b573e74d",
    "balance": 39999910000,
    "staking": "20000000000",
    "award": "",
    "jailed": false,
    "jailedBlock": 0,
    "jailedUntil": "0001-01-01T00:00:00Z"
  }
}
```


## 设置验证节点

描述:设置运行中 pocket 的验证节点账户

```shell
curl --request POST --url http://127.0.0.1:8088/poktsrv/set-validator --header 'Content-Type: application/json' \
--data "{\"address\":\"ee60841d9afb70ba893c02965537bc0eec4ef1e4\",\"passwd\": \"123456\"}"
```

参数：

- **address：** 目标地址。
- **passwd：** address 对应的 Passphrase。

输出:

```shell
{
  "status": "success",
  "code": "",
  "data": {
    "result": "spawn sh -c pocket accounts set-validator ee60841d9afb70ba893c02965537bc0eec4ef1e4\r\n\2023/03/03 03:06:37 Initializing Pocket Datadir\r\n2023/03/03 03:06:37 datadir = /home/app/.pocket\r\nEnter the password:\r\n"
  }
}
```


## 查看验证节点

描述:查看当前验证节点账户信息

```shell
curl --url http://127.0.0.1:8088/poktsrv/validator
```

输出:

```shell
{
  "status": "success",
  "code": "",
  "data": "ee60841d9afb70ba893c02965537bc0eec4ef1e4"
}
```


## Custodial抵押

描述:设置节点抵押

```shell
curl --request POST --url http://127.0.0.1:8088/poktsrv/custodial --header 'Content-Type: application/json' \
--data "{\"address\":\"ee60841d9afb70ba893c02965537bc0eec4ef1e4\",\"amount\": \"20000000000\",\"relay_chain_ids\": \"0001,0021\",\"service_url\": \"http://pokt.storefrontiers.cn:80\",\"network_id\": \"testnet\",\"fee\": \"10000\",\"is_before\": \"false\",\"passwd\": \"123456\"}"
```

参数：

- **address：** 欲质押 uPOKT 的地址。
- **amount：** 要质押的 uPOKT 数量。必须高于 StakeMinimum 当前值，可以在此处找到。
- **relay_chain_ids：** 用逗号分隔的 RelayChain 网络标识符列表。可以在此处找到网络标识符列表。
- **service_url：** 应用程序用于与中继节点通信的服务 URI。
- **network_id：** Pocket 链标识符，可以是 "mainnet" 或 "testnet"。
- **fee：** 网络所需的 uPOKT 费用。
- **is_before：** 指示是否激活了非托管升级，可以是 "true" 或 "false"。
- **passwd：** fromAddr 账户对应的 Passphrase。


输出:

```shell
{
  "status": "success",
  "code": "",
  "data": {
    "result": "spawn sh -c pocket nodes stake custodial ee60841d9afb70ba893c02965537bc0eec4ef1e4 20000000000 0001,0021 http://pokt.storefrontiers.cn:80 testnet 10000 false\r\n 2023/03/03 03:15:32 Initializing Pocket Datadir\r\n2023/03/03 03:15:32 datadir = /home/app/.pocket\r\nEnter Passphrase: \r\nhttp://localhost:8081/v1/client/rawtx\r\n{\r\n    \"logs\": null,\r\n    \"txhash\": \"0A025220D33B84525E99AFD5BE7ECA95D6234AFB40CD21901700A7F706DE12E7\"\r\n}\r\n\r\n"
  }
}
```



## Non-Custodial抵押

描述:设置节点抵押

```shell
curl --request POST --url http://127.0.0.1:8088/poktsrv/noncustodial --header 'Content-Type: application/json' \
--data "{\"public_key\":\"f75e382d77893447b8c01d9a5787f5bf7f4446d8a02e2c6ed07fb02f08b8bb83\",\"output_addr\":\"f4daee9cdacdb76f658c571e6301723817bc588a\",\"amount\": \"15100000000\",\"relay_chain_ids\": \"0001,0021\",\"service_url\": \"http://pokt.storefrontiers.cn:80\",\"network_id\": \"testnet\",\"fee\": \"10000\",\"is_before\": \"false\",\"passwd\": \"123456\"}"
```

参数：

- **public_key：** OperatorAddress 是块和中继的唯一有效签名者，其对应的公钥。
- **output_addr：** outputAddress是奖励和托管资金的目的地。
- **amount：** 要质押的 uPOKT 数量。必须高于 StakeMinimum 当前值，可以在此处找到。
- **relay_chain_ids：** 用逗号分隔的 RelayChain 网络标识符列表。可以在此处找到网络标识符列表。
- **service_url：** 应用程序用于与中继节点通信的服务 URI。
- **network_id：** Pocket 链标识符，可以是 "mainnet" 或 "testnet"。
- **fee：** 网络所需的 uPOKT 费用。
- **is_before：** 指示是否激活了非托管升级，可以是 "true" 或 "false"。
- **passwd：** OperatorAddress 账户对应的 Passphrase。

输出:

```shell
{
  "status": "success",
  "code": "",
  "data": {
    "result": "spawn sh -c pocket nodes stake custodial ee60841d9afb70ba893c02965537bc0eec4ef1e4 20000000000 0001,0021 http://pokt.storefrontiers.cn:80 testnet 10000 false\r\n 2023/03/03 03:15:32 Initializing Pocket Datadir\r\n2023/03/03 03:15:32 datadir = /home/app/.pocket\r\nEnter Passphrase: \r\nhttp://localhost:8081/v1/client/rawtx\r\n{\r\n    \"logs\": null,\r\n    \"txhash\": \"0A025220D33B84525E99AFD5BE7ECA95D6234AFB40CD21901700A7F706DE12E7\"\r\n}\r\n\r\n"
  }
}
```
