# RPC Provider Guideline
[![Discord](https://img.shields.io/discord/770382203782692945?label=Discord&logo=Discord)](https://discord.gg/MSXGzVsSYf)
[![Twitter Follow](https://img.shields.io/twitter/follow/0xfilswan)](https://twitter.com/0xfilswan)
[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg)](https://github.com/RichardLitt/standard-readme)

- 加入FilSwan的[Slack](https://filswan.slack.com)频道，了解新闻、讨论和状态更新。
- 查看FilSwan的[Medium](https://filswan.medium.com)，获取最新动态和公告。

## 目录

- [特性](#特性)
- [前提条件](#前提条件)
- [安装部署](#安装部署)
- [许可证](#许可证)

## 特性

RPC Provider Guideline 提供以下功能：

* 在容器中自动部署pocket节点。
* 提供对容器中节点基本设置查询命令。
* 提供对容器中节点维护监控的接口。

## 前提条件
- Docker
### 安装 Docker
```shell
sudo apt install docker
```
参考: [官方安装文档](https://docs.docker.com/engine/install/)

## 安装部署
### 从源代码构建安装: 
```shell
git clone https://github.com/FogMeta/meta-provider.git
cd meta-provider
git checkout release-2.1.0
./build_install_pock.sh
```

### 配置 `config-pokt.toml`
#### 编辑配置文件 **~/.swan/provider/config-pokt.toml** :
- **swan_api_url:**  Swan API 地址。 对于 Swan production， 地址为 `https://go-swan-server.filswan.com`。
- **swan_api_key:**  API key。可以通过 [Swan Platform](https://console.filswan.com/#/dashboard) -> "个人信息"->"开发人员设置" 获得， 也可以访问操作指南。
- **swan_access_token:** 访问令牌。可以通过 [Swan Platform](https://console.filswan.com/#/dashboard) -> "个人信息"->"开发人员设置"， 可以访问操作指南查看。
- **pokt_log_level:** 默认`INFO`，可选 DEBUG INFO WARN ERROR FATAL。
- **pokt_api_url:** 默认 `8081`，pocket API 端口。
- **pokt_docker_image** Docker 镜像，例如 `filswan/pocket:RC-0.9.2`。
- **pokt_docker_name** 容器名称，可自行定义，例如 `pokt-node-v0.9.2`。
- **pokt_path** pocket 数据存储路径，例如 `/root/.pocket`。
- **pokt_scan_interval** 600秒或10分钟。扫描Pocket高度状态的时间间隔。
- **pokt_heartbeat_interval:** 180秒或3分钟。在Swan平台更新状态的时间间隔。
- **pokt_server_api_url** provider pocket 服务Url，例如 `http://127.0.0.1:8088/`。
- **pokt_server_api_port** provider pocket 服务Port，例如 `8088`。
- **pokt_network_type** pocket网络类型，可以是 MAINNET 和 TESTNET 其中之一。

### 配置 `chains.json`
- 根据自身需求，配置 **~/.swan/provider/chains.json** ，例如：
```
[
    {
      "id": "0001",
      "url": "http://localhost:8081/",
      "basic_auth": {
        "username": "",
        "password": ""
      }
    },
    {
      "id": "0021",
      "url": "https://eth-rpc.gateway.pokt.network/",
      "basic_auth": {
          "username": "",
          "password": ""
      }
    }
]
```

### 下载快照
- 从最新快照下载将极大地缩短主网同步区块链所需的时间。使用wget进行下载，并在下载后解压缩存档。解压目录 `/root/.pocket` 需要与 `config-pokt.toml` 中 `pokt_path` 指定的路径保持一致。
```
mkdir -p /root/.pocket/data
wget -qO- https://snapshot.nodes.pokt.network/latest.tar.gz | tar -xz -C /root/.pocket/data
chmod -R 777 /root/.pocket
```

### 运行
- 后台运行 `meta-provider`, 其中 `start` 命令参数 `passwd` 初始创建账号的 `Passphrase`
```
nohup ./meta-provider pocket start --passwd 123456 >> provider-pokt.log 2>&1 & 
```

### 充值
- 使用命令或钱包，充值高于最低抵押值的 POCK，最低抵押值为15,000 POKT（或15,000,000,000 uPOKT）。
- 如果正在使用测试网络，可以使用[测试网络水龙头](https://faucet.pokt.network)为账户提供资金。

### 设置验证节点
- 充值到账后，通过命令设置验证节点地址：
```
# 进入容器
docker exec -it  [CONTAINER_ID] /bin/sh

# 执行设置命令
pocket accounts set-validator [YOUR_ACCOUNT_ADDRESS]

# 查看执行结果
pocket accounts get-validator
```

### 抵押
- 通过命令抵押 POCK
```
# 进入容器
docker exec -it  [CONTAINER_ID] /bin/sh

# custodial 抵押
pocket nodes stake custodial <operatorAddress> <amount> <relayChainIDs> <serviceURI> <networkID> <fee> <isBefore8.0>

# non-custodial 抵押
pocket nodes stake non-custodial <operatorPublicKey> <outputAddress> <amount> <RelayChainIDs> <serviceURI> <networkID> <fee> <isBefore8.0>
```

## 帮助

如有任何使用问题，请在 [Discord 频道](http://discord.com/invite/KKGhy8ZqzK) 联系 Meta Provider 团队或在Github上创建新的问题.

## 许可证

[Apache](https://github.com/FogMeta/meta-provider/blob/main/LICENSE)
