NETWORK=upgrade-1
DAEMON=nibirud
HOME_DIR=~/.nibiru
CONFIG=~/.nibiru/config
TOKEN_DENOM=ugame

$DAEMON tx gov submit-proposal software-upgrade signal-module-upgrade --upgrade-height=25 --upgrade-info='{"binaries":{"linux/amd64":"https://github.com/cosmos-gaminghub/nibiru/releases/download/sm-upgrade/nibirud-sm-upgrade?checksum=sha256:78d44fe51c1c04a7b0ec7b77cd197a324290659548d80d0d8526094512e8e70b"}}' --from=eg --title='sm-upgrade' --description='add signal module' --chain-id=$NETWORK -y

$DAEMON tx gov deposit 1 10000000ugame --from=eg --chain-id=$NETWORK --sequence=2 -y

$DAEMON tx gov vote 1 yes --from=eg --chain-id=$NETWORK --sequence=3 -y
