# To get started with VATZ 

> This instruction is mainly based on Linux:Ubuntu

## 1. VATZ 
1. Clone vatz (public) repository with:
    ```
    $ git clone git@github.com:dsrvlabs/vatz.git
    ```
2. Compile VATZ
    ```
    $ cd vatz
    $ make
    ```
    you will see binary named `vatz`
3. Initialize VATZ
    ```
    $ ./vatz init
    ```
    You will see config file that named default.yaml once you initialize VATZ.  

## 2. VATZ Plugins
1. Clone vatz-plugin official repository <br> (**This repository can be changed.**)
    ```
    $ git clone git@github.com:dsrvlabs/vatz-plugin-common.git
    ```
2. Compile plugins
    ```
    $ cd plugins/cosmos-sdk-blocksync
    $ make
    ```
   You can see binary named `cosmos-sdk-blocksync`

---

# Usage
## 1. Modify default.yaml (VATZ)
There's yaml file as config, and default config must be updated by yourself to set up VATZ service and clients
sample:
> This is default yaml, please update secrets, etcs for your node use.
```
vatz_protocol_info:
  protocol_identifier: "Put Your Protocol here"
  port: 9090
  health_checker_schedule:
    - "0 1 * * *"
  notification_info:
    host_name: "Put your machine's host name"
    default_reminder_schedule:
      - "*/15 * * * *"
    dispatch_channels:
      - channel: "discord"
        secret: "Put your Discord Webhook"
      - channel: "pagerduty"
        secret: "Put your PagerDuty's Integration Key (Events API v2)"
      - channel: "telegram"
        secret: "Put Your Bot's Token"
        chat_id: "Put Your Chat's chat_id"
        reminder_schedule:
          - "*/5 * * * *"
  rpc_info:
    enabled: true
    address: "127.0.0.1"
    grpc_port: 19090
    http_port: 19091
    
plugins_infos:
  default_verify_interval: 15
  default_execute_interval: 30
  default_plugin_name: "vatz-plugin"
  plugins:
    - plugin_name: "samplePlugin1"
      plugin_address: "localhost"
      plugin_port: 9001
      verify_interval: 7
      execute_interval: 9
      executable_methods:
        - method_name: "sampleMethod1"                                                                                     
```
### VATZ protocol Infos
1. `port`: check your machine used port number. If you set it to a port number that is already in use, an error will occur.
2. `health_checker_schedule`: use cron expression to set time to check if vatz is alived. 
3. `dispatch_channels`: add infos to recieve arlerts through discord, pagerduty, and telegram.
4. `reminder_schedule`: use cron expression to set time to resend the alert if it is not confirmed.


### Plugin Infos
> **You can find further plugins here.**
> - [vatz-plugin-common](https://github.com/dsrvlabs/vatz-plugin-common)
> - [vatz-plugin-sysutill](https://github.com/dsrvlabs/vatz-plugin-sysutil)
> - [vztz-plugin-cosmoshub](https://github.com/dsrvlabs/vatz-plugin-cosmoshub)

1. plugin_name
   - refer to `vatz-plugin-common/plugins/cosmos-sdk-blocksync/main.go` of official repository
      - Search `pluginName`in vat-plugin and update the default.yaml with `pluginName`
2. plugin_port
   - set your plugin used port number to connect the plugin.


## 2. Execute each binary
1. Start VATZ
    ```
    $ ./vatz start
    ```
    If you see those following command you successfully started VATZ service.
    ```
    root@validator-node:~/vatz# go run main.go
    2022/04/22 05:06:54 Initialize Servers: VATZ Manager
    2022/04/22 05:06:54 Listening Port :9090
    2022/04/22 05:06:54 Node Manager Started
    ```

2. Start VATZ plugin
    ```
    $ ./cosmos-sdk-blocksync start #put your plugin binary file name
    ```
    If you see those following command you successfully started VATZ plugin.
    ```
    2022-11-25T11:27:01+09:00 INF Start main=statusCollector
    2022-11-25T11:27:01+09:00 INF Register module=grpc
    2022-11-25T11:27:01+09:00 INF Start 127.0.0.1 9091 module=sdk
    2022-11-25T11:27:01+09:00 INF Start module=grpc
    ```    
<br> <img width="1510" alt="image" src="https://user-images.githubusercontent.com/106724973/180962626-c4de859c-526d-4038-87b8-badebed44136.png">

3. The alert notification will be sent as below if there are problems in monitored node. The alert conditions vary for each plugin. <br> <img width="223" alt="image" src="https://user-images.githubusercontent.com/106724973/180962941-4db55b66-ceb6-4cd4-bd53-96d21bdf5596.png">
