# Install Nibiru fullnode

## Install go

:::tip Required
**Go 1.17.0+** is required for the Cosmos SDK.
:::

Firstly, install `golang` from [the official golang donwload page](https://golang.org/dl/).

```sh
wget https://go.dev/dl/go1.17.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.17.6.linux-amd64.tar.gz
```

Be sure to set your `$GOPATH`, `$GOBIN`, and `$PATH` environment variables, for example:

```sh
echo '#golang' >> ~/.bashrc
echo 'export PATH="$PATH:/usr/local/go/bin"' >> ~/.bashrc
echo 'export GOPATH="$HOME/go"' >> ~/.bashrc
echo 'export PATH="$PATH:$GOPATH/bin"' >> ~/.bashrc
echo 'export GOBIN=$GOPATH/bin' >> ~/.bashrc

source ~/.bashrc
```

Verify that `golang` has been installed successfully.

```sh
go version
go version go1.17.6 linux/amd64
```


## Install Nibiru
With `golang`, you can compile and run `nibiru`.

```sh
git clone https://github.com/cosmos-gaminghub/nibiru.git
cd nibiru && git checkout -b 0.9 tags/0.9
make install
```

Try `nibirud version` to verify that everything is fine.

::: tip Tips
If you are using ubuntu, make sure to install build tools with the command `apt install build-essential` before `make install`.
:::

::: tip Tips
If you want to use a ledger device, make sure to set `LEDGER_ENABLED`. Ex: `LEDGER_ENABLED=true make install`.
:::
