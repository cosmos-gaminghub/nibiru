# Node Configuration
## Configure the service

To allow your nibiru node to run in the background as a service you need to execute the following command.

```sh
tee /etc/systemd/system/nibirud.service > /dev/null <<EOF
[Unit]
Description=Nibiru Full Node
After=network-online.target
[Service]
User=root
ExecStart=/root/go/bin/nibirud start
Restart=always
RestartSec=3
LimitNOFILE=65535
[Install]
WantedBy=multi-user.target
EOF
```

::: warning
If you are logged as a user which is not `root`, make sure to edit the User value accordingly
:::

Once you have successfully created the service, you need to enable it. After this, you can run Nibiru node.

```
systemctl daemon-reload
systemctl enable nibirud
systemctl start nibirud
```


### Service operations
```
systemctl status nibirud
```

### Check the node status

```
journalctl -f -u nibirud
```
