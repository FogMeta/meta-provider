# RPC Provider Guideline Commands

* [Start Node](#start-node)
* [View Version](#view-version)
* [View Validator Node](#view-validator-node)
* [View Balance](#view-balance)
* [Status Information](#status-information)
* [Custodial Staking](#custodial-staking)
* [Non-Custodial Staking](#non-custodial-staking)


## Start Node

Description: Deploy and run the Pocket node in a container:

- Pull the image specified by `pokt_docker_image` to the local machine;
- Create a container specified by `pokt_docker_name`, and create a Pocket initial account based on the command parameters `passwd`;
- Start the container specified by `pokt_docker_name`;
- Wait for the Pocket node to run normally in the container, and obtain the Pocket version information and block height.


```shell
meta-provider pocket start --passwd "123456"
```

Parameters:

- **passwd:** Set the `Passphrase` for the initial account.


## View Version

Description: Check the current version of running Pocket.

```shell
meta-provider pocket version
```

Output:

```shell
Pocket Version  : RC-0.9.2
```


## View Validator Node

Description: Check the current validator node account address of running Pocket.

```shell
meta-provider pocket validator
```

Output:

```shell
Validator Address       : ee60841d9afb70ba893c02965537bc0eec4ef1e4
```


## View Balance

Description: Check the balance of a specified account.

```shell
meta-provider pocket balance --addr ee60841d9afb70ba893c02965537bc0eec4ef1e4
```

Parameters:

- **addr:** Address of the account to query.

Output:

```shell
Address : ee60841d9afb70ba893c02965537bc0eec4ef1e4
Balance : 39999970000
```


## Status Information

Description: Check the status information of the running Pocket node.

```shell
meta-provider pocket status
```

Output:

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


## Custodial Staking

Description: Set up node staking with custodial mode.

```shell
meta-provider pocket custodial --fromAddr="ee60841d9afb70ba893c02965537bc0eec4ef1e4" --amount="20000000000" --relayChainIDs="0001,0021" --serviceURI="http://pokt.storefrontiers.cn:80" --networkID="testnet" --fee="10000" --isBefore="false" --passwd="123456"

```

Parameters:

- **fromAddr:** Address of the account to stake from.
- **amount:** Number of uPOKT tokens to stake.
- **relayChainIDs:** List of RelayChain network identifiers separated by commas.
- **serviceURI:** Service URI used by applications to communicate with Relay Nodes.
- **networkID:** Pocket chain identifier, can be "mainnet" or "testnet".
- **fee:** Amount of uPOKT fee required by the network.
- **isBefore:** Indicates if non-custodial upgrade is activated, can be "true" or "false".
- **passwd:** Passphrase for the fromAddr account.


Output:

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


## Non-Custodial Staking

Description: Set up node staking with non-custodial mode.

```shell
meta-provider pocket non-custodial --operatorPublicKey="f75e382d77893447b8c01d9a5787f5bf7f4446d8a02e2c6ed07fb02f08b8bb83" --outputAddress="f4daee9cdacdb76f658c571e6301723817bc588a" --amount="20000000000" --relayChainIDs="0001,0021" --serviceURI="http://pokt.storefrontiers.cn:80" --networkID="testnet" --fee="10000" --isBefore="false" --passwd="123456"
```

Parameters:

- **operatorPublicKey:** The unique valid signer for blocks and relays, corresponding to its public key.
- **outputAddress:** The destination of the rewards and staked funds.
- **amount:** Number of uPOKT tokens to stake. Must be higher than the current StakeMinimum value.
- **relayChainIDs:** List of RelayChain network identifiers separated by commas. 
- **serviceURI:** Service URI used by applications to communicate with Relay Nodes.
- **networkID:** Pocket chain identifier, can be "mainnet" or "testnet".
- **fee:** Amount of uPOKT fee required by the network.
- **isBefore:** Indicates if non-custodial upgrade is activated, can be "true" or "false".
- **passwd:** Passphrase for the OperatorAddress account.

Output:

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

