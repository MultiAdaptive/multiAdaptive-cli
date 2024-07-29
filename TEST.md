# multiAdaptive
Please use a new wallet.  
This activity supports Linux/Mac/Windows systems to participate in the test.

## Preparing to Test ETH
You can get test ETH from the following faucets:

- [Chainlink Faucet](https://faucets.chain.link/)
- [Sepolia Faucet](https://sepolia-faucet.pk910.de/)



## Linux/Mac

### Downloading the Client

To download the client, run the following command:

```sh
curl -L https://github.com/MultiAdaptive/multiAdaptive-cli/raw/master/scripts/download.sh | bash
```

### Participating in the Test  
To participate in the test, run the client with your private key:
```sh
./multiAdaptive-cli -privateKey="<your privateKey>"
```
Feel free to replace <your privateKey> with your actual private key.

## Windows

### Downloading the Client

To download the client, run the following command:

```sh
curl -L -o download.bat https://github.com/MultiAdaptive/multiAdaptive-cli/raw/master/scripts/download.bat && download.bat
```

### Participating in the Test
To participate in the test, run the client with your private key:
```sh
./multiAdaptive-cli.exe -privateKey="<your privateKey>"
```
Feel free to replace <your privateKey> with your actual private key.
