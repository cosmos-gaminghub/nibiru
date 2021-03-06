# Install Nibiru fullnode

## Install go

:::tip Required
**Go 1.16.0+** is required for the Cosmos SDK.
:::

Firstly, install `golang` from [the official golang donwload page](https://golang.org/dl/).
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
$ go version
go version go1.16.4 linux/amd64
```


## Install Nibiru
With `golang`, you can compile and run `nibiru`.

```sh
git clone https://github.com/cosmos-gaminghub/nibiru.git
cd nibiru && git checkout -b tag0.3 tags/v0.3
make install
```

Try `nbrcli version` and `nbrd version` to verify that everything is fine.

