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
Now let's try to use Cosmovisor in your localnet. In this simulation, you will upddate nibirud version from v to v.

[TODO]
