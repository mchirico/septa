
## Creating a Service

```bash

useradd -m -d /docker/routefirebase routefirebase
sudo -i -u routefirebase

touch /etc/systemd/system/routefirebase.service

# Copy files


```

## Edit /etc/systemd/system/routefirebase.service

```bash

[Unit]
Description=Route Firebase Service  
After=network.target

[Service]
Type=simple
WorkingDirectory=/docker/routefirebase
User=routefirebase
ExecStart=/docker/routefirebase/septa-project/wrapperScript.sh
Restart=always
RestartSec=10
StandardOutput=syslog
StandardError=syslog

[Install]
WantedBy=multi-user.target

```


## Install the service

```bash
systemctl daemon-reload;sudo systemctl enable routefirebase.service
```

## Test User

```bash
su -s /bin/bash routefirebase

```


## Start/Stop/Status

```bash
systemctl start routefirebase
systemctl stop routefirebase
systemctl status routefirebase
systemctl restart routefirebase

```

## Check Service

```bash
systemctl status routefirebase
● routefirebase.service - Route Firebase Service
   Loaded: loaded (/etc/systemd/system/routefirebase.service; enabled)
   Active: active (running) since Sun 2018-04-01 00:45:31 UTC; 2min 11s ago
 Main PID: 22563 (wrapperScript.s)
   CGroup: /system.slice/routefirebase.service
           ├─22563 /bin/bash /docker/routefirebase/septa-project/wrapperScript.sh
           └─22564 /docker/routefirebase/septa-project/bin/routefirebase -token=/docker/routefirebase/septa-project/tokens/token.json -qui...

Apr 01 00:45:31 mce8 systemd[1]: Started Route Firebase Service.

```

