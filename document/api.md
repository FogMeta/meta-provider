# RPC Provider Guideline API

* [Version](#version)
* [Height](#height)
* [Account Balance](#account-balance)
* [Status Information](#status-information)
* [Set Validator Node](#set-validator-node)
* [View Validator Node](#view-validator-node)
* [Custodial Staking](#custodial-staking)
* [Non-Custodial Staking](#non-custodial-staking)

## Version

Description: Check the current version of running Pocket.

```shell
curl --url http://127.0.0.1:8088/poktsrv/version 
```

Output:

```shell
{
  "status": "success",
  "code": "",
  "data": {
    "version": "RC-0.9.2"
  }
}
```


## Height

Description: Check the current block height of running Pocket.

```shell
curl --url http://127.0.0.1:8088/poktsrv/height
```

Output:

```shell
{
  "status": "success",
  "code": "",
  "data": {
    "height": 99156
  }
}
```


## Account Balance

Description: Query the balance of a specified account.

```shell
curl --request POST --url http://127.0.0.1:8088/poktsrv/balance --header 'Content-Type: application/json' \
--data "{\"height\": 0,\"address\":\"ee60841d9afb70ba893c02965537bc0eec4ef1e4\"}"
```

Parameters:

- **height:** The specified height of the block to query. Defaults to 0, which will query the latest known block by the node.
- **address:** The target address.

Output:

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


## Status Information

Description: Check the status information of running Pocket.

```shell
curl --url http://127.0.0.1:8088/poktsrv/status
```

Output:

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


## Set Validator Node

Description: Set up a validator node account for running Pocket

```shell
curl --request POST --url http://127.0.0.1:8088/poktsrv/set-validator --header 'Content-Type: application/json' \
--data "{\"address\":\"ee60841d9afb70ba893c02965537bc0eec4ef1e4\",\"passwd\": \"123456\"}"
```

Parameters:

- **address:** The target address.
- **passwd:** Passphrase corresponding to the address.

Output:

```shell
{
  "status": "success",
  "code": "",
  "data": {
    "result": "spawn sh -c pocket accounts set-validator ee60841d9afb70ba893c02965537bc0eec4ef1e4\r\n\2023/03/03 03:06:37 Initializing Pocket Datadir\r\n2023/03/03 03:06:37 datadir = /home/app/.pocket\r\nEnter the password:\r\n"
  }
}
```


## View Validator Node

Description: View current validator node account information

```shell
curl --url http://127.0.0.1:8088/poktsrv/validator
```

Output:

```shell
{
  "status": "success",
  "code": "",
  "data": "ee60841d9afb70ba893c02965537bc0eec4ef1e4"
}
```


## Custodial Staking

Description: Set up node staking

```shell
curl --request POST --url http://127.0.0.1:8088/poktsrv/custodial --header 'Content-Type: application/json' \
--data "{\"address\":\"ee60841d9afb70ba893c02965537bc0eec4ef1e4\",\"amount\": \"20000000000\",\"relay_chain_ids\": \"0001,0021\",\"service_url\": \"http://pokt.storefrontiers.cn:80\",\"network_id\": \"testnet\",\"fee\": \"10000\",\"is_before\": \"false\",\"passwd\": \"123456\"}"
```

Parameters:

- **address:** The address to stake uPOKT from.
- **amount:** The amount of uPOKT to stake. Must be above the current StakeMinimum value.
- **relay_chain_ids:** Comma-separated list of RelayChain network identifiers.
- **service_url:** Service URI that an application uses to communicate with a relay node.
- **network_id:** Pocket chain identifier which can either be "mainnet" or "testnet".
- **fee:** The uPOKT fee required by the network.
- **is_before:** Indicates whether non-custodial upgrade has been activated and can be either "true" or "false".
- **passwd:** Passphrase of the fromAddr account.

Output:

```shell
{
  "status": "success",
  "code": "",
  "data": {
    "result": "spawn sh -c pocket nodes stake custodial ee60841d9afb70ba893c02965537bc0eec4ef1e4 20000000000 0001,0021 http://pokt.storefrontiers.cn:80 testnet 10000 false\r\n 2023/03/03 03:15:32 Initializing Pocket Datadir\r\n2023/03/03 03:15:32 datadir = /home/app/.pocket\r\nEnter Passphrase: \r\nhttp://localhost:8081/v1/client/rawtx\r\n{\r\n    \"logs\": null,\r\n    \"txhash\": \"0A025220D33B84525E99AFD5BE7ECA95D6234AFB40CD21901700A7F706DE12E7\"\r\n}\r\n\r\n"
  }
}
```



## Non-Custodial Staking

Description: Set up node staking

```shell
curl --request POST --url http://127.0.0.1:8088/poktsrv/noncustodial --header 'Content-Type: application/json' \
--data "{\"public_key\":\"f75e382d77893447b8c01d9a5787f5bf7f4446d8a02e2c6ed07fb02f08b8bb83\",\"output_addr\":\"f4daee9cdacdb76f658c571e6301723817bc588a\",\"amount\": \"15100000000\",\"relay_chain_ids\": \"0001,0021\",\"service_url\": \"http://pokt.storefrontiers.cn:80\",\"network_id\": \"testnet\",\"fee\": \"10000\",\"is_before\": \"false\",\"passwd\": \"123456\"}"
```

Parameters:

- **public_key:** The public key corresponding to the OperatorAddress which is the unique valid signer for blocks and relays.
- **output_addr:** The destination for rewards and staked funds.
- **amount:** The amount of uPOKT to stake. Must be above the current StakeMinimum value.
- **relay_chain_ids:** Comma-separated list of RelayChain network identifiers. 
- **service_url:** Service URI that an application uses to communicate with a relay node.
- **network_id:** Pocket chain identifier which can either be "mainnet" or "testnet".
- **fee:** The uPOKT fee required by the network.
- **is_before:** Indicates whether non-custodial upgrade has been activated and can be either "true" or "false".
- **passwd:** Passphrase of the OperatorAddress account.

Output:

```shell
{
  "status": "success",
  "code": "",
  "data": {
    "result": "spawn sh -c pocket nodes stake custodial ee60841d9afb70ba893c02965537bc0eec4ef1e4 20000000000 0001,0021 http://pokt.storefrontiers.cn:80 testnet 10000 false\r\n 2023/03/03 03:15:32 Initializing Pocket Datadir\r\n2023/03/03 03:15:32 datadir = /home/app/.pocket\r\nEnter Passphrase: \r\nhttp://localhost:8081/v1/client/rawtx\r\n{\r\n    \"logs\": null,\r\n    \"txhash\": \"0A025220D33B84525E99AFD5BE7ECA95D6234AFB40CD21901700A7F706DE12E7\"\r\n}\r\n\r\n"
  }
}
```
