# RPC Provider Guideline 命令

* [启动节点](#启动节点)
* [查看版本](#查看版本)
* [查看验证节点](#查看验证节点)
* [查看余额](#查看余额)
* [状态信息](#状态信息)
* [Custodial抵押](#Custodial抵押)
* [Non-Custodial抵押](#Non-Custodial抵押)


## 启动节点

描述: 在容器中部署运行pocket节点：
- 拉取 `pokt_docker_image` 指定的image镜像到本地;
- 创建 `pokt_docker_name` 指定的容器，并根据命令参数passwd，创建pocket初始账号;
- 启动 `pokt_docker_name` 指定的容器；
- 等待容器中 pocket node 正常运行，获取pocket版本信息及区块高度。

```shell
meta-provider pocket start --passwd "123456"
```

参数：

- **passwd：** 设置初始账户的 `Passphrase`。


## 查看版本

描述: 检查运行中 pocket 的当前版本

```shell
meta-provider pocket version
```

输出:

```shell
Pocket Version  : RC-0.9.2
```


## 查看验证节点

描述: 检查运行中 pocket 的当前验证节点账户地址

```shell
meta-provider pocket validator
```

输出:

```shell
Validator Address       : ee60841d9afb70ba893c02965537bc0eec4ef1e4
```


## 查看余额

描述: 查看指定账户的余额

```shell
meta-provider pocket balance --addr ee60841d9afb70ba893c02965537bc0eec4ef1e4
```

参数：

- **addr：** 查询账户地址。

输出:

```shell
Address : ee60841d9afb70ba893c02965537bc0eec4ef1e4
Balance : 39999970000
```


## 状态信息

描述: 检查运行中 pocket 节点的状态信息

```shell
meta-provider pocket status
```

输出:

```shell
Version         : RC-0.9.2
Height          : 99131
Synced          : true
Address         : ee60841d9afb70ba893c02965537bc0eec4ef1e4
PublicKey       : 7b1739685dcdc10fcc02bc21dd822ef3458fcf543cc89487af9fe512b573e74d
Balance         : 39999970000
Staking         : 20000000000
Jailed          : false
JailedBlock     : 0
JailedUntil     : 0001-01-01 00:00:00 +0000 UTC
```


## Custodial抵押

描述:设置节点抵押

```shell
meta-provider pocket custodial --fromAddr="ee60841d9afb70ba893c02965537bc0eec4ef1e4" --amount="20000000000" --relayChainIDs="0001,0021" --serviceURI="http://pokt.storefrontiers.cn:80" --networkID="testnet" --fee="10000" --isBefore="false" --passwd="123456"

```

参数：

- **fromAddr：** 欲质押 uPOKT 的地址。
- **amount：** 要质押的 uPOKT 数量。必须高于 StakeMinimum 当前值，可以在此处找到。
- **relayChainIDs：** 用逗号分隔的 RelayChain 网络标识符列表。可以在此处找到网络标识符列表。
- **serviceURI：** 应用程序用于与中继节点通信的服务 URI。
- **networkID：** Pocket 链标识符，可以是 "mainnet" 或 "testnet"。
- **fee：** 网络所需的 uPOKT 费用。
- **isBefore：** 指示是否激活了非托管升级，可以是 "true" 或 "false"。
- **passwd：** fromAddr 账户对应的 Passphrase。

输出:

```shell
{
    Result: spawn sh -c pocket nodes stake custodial ee60841d9afb70ba893c02965537bc0eec4ef1e4 20000000000 0001,0021 http://pokt.storefrontiers.cn:80 testnet 10000 false
    2023/03/02 21:50:02 Initializing Pocket Datadir
    2023/03/02 21:50:02 datadir = /home/app/.pocket
    Enter Passphrase: 
    http://localhost:8081/v1/client/rawtx
    {
        "logs": null,
        "txhash": "487F8E6FEFCDB1B8324572B411DC1E4239CEAA915958FB06BA6E6655978ADF43"
    }
}
```


## Non-Custodial抵押

描述:设置节点抵押

```shell
meta-provider pocket non-custodial --operatorPublicKey="f75e382d77893447b8c01d9a5787f5bf7f4446d8a02e2c6ed07fb02f08b8bb83" --outputAddress="f4daee9cdacdb76f658c571e6301723817bc588a" --amount="20000000000" --relayChainIDs="0001,0021" --serviceURI="http://pokt.storefrontiers.cn:80" --networkID="testnet" --fee="10000" --isBefore="false" --passwd="123456"
```

参数：

- **operatorPublicKey：** OperatorAddress 是块和中继的唯一有效签名者，其对应的公钥。
- **outputAddress：** outputAddress是奖励和托管资金的目的地。
- **amount：** 要质押的 uPOKT 数量。必须高于 StakeMinimum 当前值，可以在此处找到。
- **relayChainIDs：** 用逗号分隔的 RelayChain 网络标识符列表。可以在此处找到网络标识符列表。
- **serviceURI：** 应用程序用于与中继节点通信的服务 URI。
- **networkID：** Pocket 链标识符，可以是 "mainnet" 或 "testnet"。
- **fee：** 网络所需的 uPOKT 费用。
- **isBefore：** 指示是否激活了非托管升级，可以是 "true" 或 "false"。
- **passwd：** OperatorAddress 账户对应的 Passphrase。

输出:

```shell
{
    result: "spawn sh -c pocket nodes stake non-custodial f75e382d77893447b8c01d9a5787f5bf7f4446d8a02e2c6ed07fb02f08b8bb83 f4daee9cdacdb76f658c571e6301723817bc588a 15100000000 0001,0021 http://pokt.storefrontiers.cn:80 testnet 10000 false"
    2023/03/12 23:59:29 Initializing Pocket Datadir
    2023/03/12 23:59:29 datadir = /home/app/.pocket
    Enter Passphrase: 
    http://localhost:8081/v1/client/rawtx
    {  
      "logs": null,
      "txhash": "2C175EB657C665CAAAAE0303F3DF13CEFCEDE1FE5349D227B449CF3B42C13515"
    }
}
```

