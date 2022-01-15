# Use Cosmovisor

::: tip Tips
Check [the Cosmovisor official document](https://github.com/cosmos/cosmos-sdk/tree/master/cosmovisor) 
:::

The Cosmos team provides a tool named **Cosmovisor** that allows your node to perform some automatic operations when needed. This is particularly useful when dealing with on-chain upgrades, because Cosmovisor can help you by taking care of downloading the updated binary and restarting the node for you.  

If you want to learn how to setup Cosmovisor inside your full or validator node, please follow the guide below. 

## Setup
### Downloading Cosmovisor
To install the latest version of `cosmovisor`, run the following command:
```sh
go install github.com/cosmos/cosmos-sdk/cosmovisor/cmd/cosmovisor@latest
```

To install a previous version, you can specify the version. 

```sh
go install github.com/cosmos/cosmos-sdk/cosmovisor/cmd/cosmovisor@v1.0.0
```

To check your cosmovisor version, you can use this command and make sure it matches the version you've installed:
```sh
strings $(which cosmovisor) | egrep -e "mod\s+github.com/cosmos/cosmos-sdk/cosmovisor"
```

### Setting up environmental variables
Cosmovisor relies on the following environmental variables to work properly:

| Env | Value | Description |
| --------|---------|---|
| `DAEMON_NAME` (required)| nibirud | The name of the binary itself
| `DAEMON_HOME` (required)| $HOME/.nibiru | The location where upgrade binaries should be kept
| `DAEMON_RESTART_AFTER_UPGRADE` (optional) | true | If set to true, it will restart the process with the new binary after a successful upgrade.
| `DAEMON_ALLOW_DOWNLOAD_BINARIES` (optional) | false | If set to true,  it will enable auto-downloading of new binaries (for security reasons, this is intended for full nodes rather than validators)
| `UNSAFE_SKIP_BACKUP` (optional) | false | If set to true, it upgrades directly without performing a backup.

**IMPORTANT**: If you don't have much free disk space, please set `UNSAFE_SKIP_BACKUP=true` to avoid your node failing the upgrade due to insufficient disk space when creating the backup.

#### Updating the service file
If you are running your node using a service( ref: [service](./service.md) ), you need to update your service file to use `cosmovisor` instead of `nibirud`. To do this you can simply run the following command:

```shell
sudo tee /etc/systemd/system/nibirud.service > /dev/null <<EOF  

[Unit]
Description=Nibiru Full Node
After=network-online.target

[Service]
User=root
ExecStart=$(which cosmovisor) start
Restart=always
RestartSec=3
LimitNOFILE=65535

Environment="DAEMON_HOME=$HOME/.nibiru"
Environment="DAEMON_NAME=nibirud"
Environment="DAEMON_RESTART_AFTER_UPGRADE=true"
Environment="DAEMON_ALLOW_DOWNLOAD_BINARIES=false"
Environment="UNSAFE_SKIP_BACKUP=false"

[Install]
WantedBy=multi-user.target
EOF
```


Once you have edited your system file, you need to reload it using the following command:

```shell
sudo systemctl daemon-reload
```

Finally, you can restart is as follows: 

```shell
sudo systemctl restart nibirud
```

## Practice
Now let's try to use Cosmovisor in your localnet. In this simulation, you will upddate nibirud version from `v0.9` to `sm-upgrade`.


### Preparation before network launch

```sh
git checkout -b v0.9 tags/v0.9
make install

# build binary will be used for cosmovisor
make build
```

```sh
NETWORK=upgrade-1
DAEMON=nibirud
HOME_DIR=~/.nibiru
CONFIG=~/.nibiru/config
TOKEN_DENOM=ugame

rm -rf $HOME_DIR

$DAEMON init $NETWORK --chain-id $NETWORK

$DAEMON config chain-id $NETWORK
$DAEMON config keyring-backend test

$DAEMON keys add eg --keyring-backend test
$DAEMON add-genesis-account $($DAEMON keys show eg -a --keyring-backend test) 100000000000000$TOKEN_DENOM

sed -i "s/\"stake\"/\"$TOKEN_DENOM\"/g" $CONFIG/genesis.json
jq '.app_state.gov.voting_params.voting_period = "60s"' $CONFIG/genesis.json > tmp.json && mv tmp.json $CONFIG/genesis.json

$DAEMON gentx eg 50000000000000$TOKEN_DENOM --commission-max-change-rate=0.1 --commission-max-rate=1 --commission-rate=0.1 --moniker=eg-validator --keyring-backend test --chain-id $NETWORK

$DAEMON collect-gentxs

$DAEMON validate-genesis

```

Cosmovisor settings
```sh
export DAEMON_NAME=nibirud
export DAEMON_HOME=$HOME/.nibiru
export DAEMON_RESTART_AFTER_UPGRADE=true
export DAEMON_ALLOW_DOWNLOAD_BINARIES=true

# make sure v0.9 nibirud is in genesis/bin
mkdir -p $DAEMON_HOME/cosmovisor/genesis/bin
cp ./build/nibirud $DAEMON_HOME/cosmovisor/genesis/bin
```

Set target binary to upgrade
```sh
git checkout -b sm-upgrade tags/sm-upgrade
make build

# make sure sm-upgrade nibirud is in upgrades/signal-module-upgrade/bin
mkdir -p $DAEMON_HOME/cosmovisor/upgrades/signal-module-upgrade/bin
cp ./build/nibirud $DAEMON_HOME/cosmovisor/upgrades/signal-module-upgrade/bin
```

Now we are ready. Start the localnet using cosmovisor with the command `cosmovisor start`. In this simulation, there is 60 seconds before voting ends.

In another shell, you have to do three things below.

1. submit `softwareUpgrade` proposal
2. deposit fund
3. vote for the proposal

```sh
NETWORK=upgrade-1
DAEMON=nibirud
HOME_DIR=~/.nibiru
CONFIG=~/.nibiru/config
TOKEN_DENOM=ugame

$DAEMON tx gov submit-proposal software-upgrade sm-upgrade --upgrade-height=25 --upgrade-info='{"binaries":{"linux/amd64":"https://github.com/cosmos-gaminghub/nibiru/releases/download/sm-upgrade/nibirud-sm-upgrade?checksum=sha256:78d44fe51c1c04a7b0ec7b77cd197a324290659548d80d0d8526094512e8e70b"}}' --from=eg --title='sm-upgrade' --description='add signal module' --chain-id=$NETWORK

$DAEMON tx gov deposit 1 10000000ugame --from=eg --chain-id=$NETWORK

$DAEMON tx gov vote 1 yes --from=eg --chain-id=$NETWORK
```

:::tip
`--upgrade-info` is used for the auto downloading feature(`DAEMON_ALLOW_DOWNLOAD_BINARIES=true`)of cosmovisor@v0.1.0. 

:::
:::warning
In v1.0.0, this auto downloading feature is **NOT** available. So try just for test purpose only.
:::
:::tip
you can check the checksum with the command `sha256sum <file name>`.
:::


If everything is ok, you can see no downtime upgrade from `v0.9` to `sm-upgrade`.
