# RPC Provider Guideline
[![Discord](https://img.shields.io/discord/770382203782692945?label=Discord&logo=Discord)](https://discord.gg/MSXGzVsSYf)
[![Twitter Follow](https://img.shields.io/twitter/follow/0xfilswan)](https://twitter.com/0xfilswan)
[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg)](https://github.com/RichardLitt/standard-readme)

- Join us on our public [discord](https://discord.com/invite/KKGhy8ZqzK) channel for news, discussions, and status updates.
- Check out our [Blog](https://blog.filswan.com/) for the latest posts and announcements.

## Directory

- [Features](#Features)
- [Prerequisites](#Prerequisites)
- [Installation](#Installation)
- [Help](#Help)
- [License](#License)

## Features

The RPC Provider Guideline provides the following features:

* Automatically deploy a pocket node in a container.
* Provide commands to query basic settings of the node in the container.
* Provide an interface to monitor the maintenance of nodes in the container.

## Prerequisites
- Docker
### Install Docker
```shell
sudo apt install docker
```
Reference: [Official Installation Documentation](https://docs.docker.com/engine/install/)

## Installation
### Build and Install from Source Code:
```shell
git clone https://github.com/FogMeta/meta-provider.git
cd meta-provider
git checkout rpc-provider
./build_install_pokt.sh
```

### Configure `config-pokt.toml`
#### Edit the configuration file `~/.swan/provider/config-pokt.toml` :
- **swan_api_url:**  Swan API address. For Swan production, it is "https://go-swan-server.filswan.com"
- **swan_api_key:**  Your api key. Acquire from Filswan -> "My Profile"->"Developer Settings". You can also check the Guide.
- **swan_access_token:** Your access token. Acquire from Filswan -> "My Profile"->"Developer Settings". You can also check the Guide.
- **pokt_log_level:** The default is `INFO`, and it can also be set to `DEBUG`, `INFO`, `WARN`, `ERROR`, or `FATAL`.
- **pokt_api_url:** The default is `8081`, which is the pocket API port.
- **pokt_docker_image** The Docker image, for example, `filswan/pocket:RC-0.9.2`.
- **pokt_docker_name** The container name, which can be defined as desired, for example, `pokt-node-v0.9.2`.
- **pokt_path** The Pocket data storage path, for example, `/root/.pocket`.
- **pokt_scan_interval** 600 seconds. The time interval for scanning Pocket height status.
- **pokt_heartbeat_interval:** 180 seconds. The time interval for updating the status on the Swan platform.
- **pokt_server_api_url** Provider pocket service URL, for example, "http://127.0.0.1:8088/".
- **pokt_server_api_port** Provider pocket service port, for example, `8088`.
- **pokt_network_type** Pocket network type, which can be either `MAINNET` or `TESTNET`.

### Configure `chains.json`
- Configure `~/.swan/provider/chains.json` based on your requirements, for example:
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

### Download Snapshot
- Downloading the latest snapshot will greatly reduce the time required to synchronize with the blockchain. Use wget to download it, and then extract the archive. The directory `/root/.pocket` where you extract the archive should match the path specified in `config-pokt.toml` for `pokt_path`.
```
mkdir -p /root/.pocket/data
wget -qO- https://snapshot.nodes.pokt.network/latest.tar.gz | tar -xz -C /root/.pocket/data
chmod -R 777 /root/.pocket
```

### Run
- Run the `meta-provider` in the background using the start command parameter `passwd` as the `passphrase` for creating an account.
```
nohup ./meta-provider pocket start --passwd 123456 >> provider-pokt.log 2>&1 & 
```

### Deposit
- Deposit POKT tokens using a command or wallet. The minimum staking requirement is 15,000 POKT (or 15,000,000,000 uPOKT).
- You can use the test network [faucet](https://faucet.pokt.network) to fund your account if you're using the test network.

### Set Validator Node Address
- After your deposit has been received, set the validator node address using the pocket accounts set-validator command:
```
# Enter the container
docker exec -it  [CONTAINER_ID] /bin/sh

# Execute the command to set the validator node address
pocket accounts set-validator [YOUR_ACCOUNT_ADDRESS]

# Check the result
pocket accounts get-validator
```

### Stake
- Stake your POKT tokens using the stake command. There are custodial and non-custodial staking options available.
```
# Enter the container
docker exec -it  [CONTAINER_ID] /bin/sh

# Custodial staking
pocket nodes stake custodial <operatorAddress> <amount> <relayChainIDs> <serviceURI> <networkID> <fee> <isBefore8.0>

# Non-custodial staking
pocket nodes stake non-custodial <operatorPublicKey> <outputAddress> <amount> <RelayChainIDs> <serviceURI> <networkID> <fee> <isBefore8.0>
```

## Help

If you encounter any issues, you can contact the Meta Provider team via the [Discord](http://discord.com/invite/KKGhy8ZqzK) channel or create a new issue on GitHub.

## License

This software is licensed under [Apache](https://github.com/FogMeta/meta-provider/blob/main/LICENSE).
