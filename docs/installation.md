# To get started with VATZ 

> This instruction is mainly based on Linux:Ubuntu
## 0. Prerequisites 
Install Go 
> You can skip this part if you already have installed go version > 1.18+ on your system. 

Please, check [latest version](https://github.com/dsrvlabs/vatz/releases) or check on [VATZ module](https://github.com/dsrvlabs/vatz/blob/main/go.mod) and follow these [instructions](https://go.dev/doc/install) to install. 


## 1. VATZ 
### 1. Clone VATZ repository

  ```
  $ git clone git@github.com:dsrvlabs/vatz.git
  ```

### 2. Compile VATZ

  ```
  $ cd vatz
  $ make
  ```
You will see binary named `vatz`

  
### 3. Initialize VATZ
  ```
  $ ./vatz init
  ```
  You will see config file named default.yaml once you initialize VATZ.

## 2. VATZ Plugin

### 1. Copy the plugin address

  We officially provide 3 plugins for System Utility and 5 plugins for CosmosHub node. 
  > **Official Plugin Repositories**
  > - [vatz-plugin-sysutill](https://github.com/dsrvlabs/vatz-plugin-sysutil) 
  >   - cpu_monitor
  >   - disk_monitor
  >   - mem_monitor
  > - [vztz-plugin-cosmoshub](https://github.com/dsrvlabs/vatz-plugin-cosmoshub)
  >   - node_active_status
  >   - node_block_sync
  >   - node_governance_alarm
  >   - node-is_alived
  >   - node_peer_count

 But you can easily develop custom plugins for the feature that you want with [VATZ SDK](https://github.com/dsrvlabs/vatz/tree/main/sdk). You could fine community provided plugins in [community_plugins.md](https://github.com/dsrvlabs/vatz/blob/main/docs/community_plugins.md). Feel free to add your custom plugins to community_plugins.md to share with others.


### 2. Install the plugin

  ```
   $ ./vatz plugin install <plugin_address> <name>
  ```
  Put git address of the plugin you want to install. You can set the plugin name as desired.

For example,
  ```
  $ ./vatz plugin install https://github.com/dsrvlabs/vatz-plugin-sysutil/plugins/cpu_monitor cpu_monitor
  ```

### 3. Check installation success

  ```
  $ ./vatz plugin list
  ```

  Verify that the plugin is installed successfully by checking whether the plugin name is added to the list as below.
  ```
  2023-05-26T11:14:28+09:00 INF Load Config default.yaml module=config
  2023-05-26T11:14:28+09:00 INF List plugins module=plugin
  2023-05-26T11:14:28+09:00 INF Create DB Instance module=db
  2023-05-26T11:14:28+09:00 INF List module=plugin
  2023-05-26T11:14:28+09:00 INF List module=db
  +-------------+------------+---------------------+---------------------------------------------------------------------+---------+
  | NAME        | IS ENABLED | INSTALL DATE        | REPOSITORY                                                          | VERSION |
  +-------------+------------+---------------------+---------------------------------------------------------------------+---------+
  | cpu_monitor | true       | 2023-05-26 11:14:22 | https://github.com/dsrvlabs/vatz-plugin-sysutil/plugins/cpu_monitor | latest  |
  +-------------+------------+---------------------+---------------------------------------------------------------------+---------+
  ```
---

# Usage
## 1. Modify default.yaml (VATZ)
There's yaml file as config, and the default config must be updated by yourself to set up VATZ service and clients
sample:
> This is default.yaml, please update secrets, and etc for your node use.
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

  monitoring_info:
    prometheus:
      enabled: true
      address: "127.0.0.1"
      port: 18080

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
### Update VATZ protocol Infos

1. `port`: Check your machine used port number. If you set it to a port number that is already in use, an error will occur.
2. `health_checker_schedule`: Use cron expression to set the time to check if vatz is alive. 
3. `notificaiton_info`:
   - `dispatch_channels`: Add infos to receive alerts through discord, pagerduty, and telegram. 
   - `reminder_schedule`: Use cron expression to set the time to resend the alert if it is not confirmed.
4. `rpc_info`: Check address and port to connect with plugins.
5. `monitoring_info`: Check address and port to use monitoring. Currently, VATZ supports monitoring only through Prometheous-Grafana.


## Update Monitoring info

### Update Plugin Infos

1. `plugin_name`: Refer to the plugin name you installed.
2. `plugin_port`: Set your plugin used port number to connect the plugin.
3. `verify_interval`: You can set the interval to check if the plugin is connected. The default interval is 15s.
4. `execute_interval`: You can set the interval to execute the plugin method. The default interval is 30s.


## 2. Start VATZ and VATZ Plugin

### 1. Start VATZ plugin

  ```
  $ ./vatz plugin start --plugin <name> --args <argumnets> --log <logfile>
  ```
  Put the plugin name you want to start. 

  For example,

  ```
  $ ./vatz plugin start -p cpu_monitor
  ```

  Or you could also start plugin with some arguments too. Check out the plugin repository to find available arguments.

  ```
  $ ./vatz plugin start -p cpu_monitor -a "-port=9094 -urgent=80"
  ```
  If you see the following commands, you successfully started VATZ plugin.

  ```
  2023-05-26T13:05:25+09:00 INF Load Config default.yaml module=config
  2023-05-26T13:05:25+09:00 INF Start plugin cpu_monitor -port=9094 -urgent=80 module=plugin
  2023-05-26T13:05:25+09:00 INF Plugin log redirect to /Users/user/.vatz/cpu_monitor.log module=plugin
  2023-05-26T13:05:25+09:00 INF Create DB Instance module=db
  2023-05-26T13:05:25+09:00 INF Start plugin cpu_monitor module=plugin
  2023-05-26T13:05:25+09:00 INF Get cpu_monitor module=db
  ```    


### 2. Start VATZ
  ```
  $ ./vatz start
  ```
  If you see following logs you successfully started VATZ service.
  ```
  % ./vatz start
  2023-05-26T13:13:23+09:00 INF start module=main
  2023-05-26T13:13:23+09:00 INF load config default.yaml module=main
  2023-05-26T13:13:23+09:00 INF logfile  module=main
  2023-05-26T13:13:23+09:00 INF Load Config default.yaml module=config
  2023-05-26T13:13:23+09:00 INF Initialize Servers: VATZ Manager module=main
  2023-05-26T13:13:23+09:00 INF VATZ Listening Port: :9090 module=main
  2023-05-26T13:13:23+09:00 INF start metric server: 127.0.0.1:18080 module=main
  2023-05-26T13:13:23+09:00 INF start rpc server module=rpc
  2023-05-26T13:13:23+09:00 INF start gRPC gateway server 127.0.0.1:19091 module=rpc
  2023-05-26T13:13:23+09:00 INF Create DB Instance module=db
  2023-05-26T13:13:23+09:00 INF start gRPC server 127.0.0.1:19090 module=rpc
  2023-05-26T13:13:23+09:00 INF Get cpu_monitor module=plugin
  2023-05-26T13:13:23+09:00 INF Get cpu_monitor module=db
  ```


### 3. Notification
The alert notification will be sent to configured channels if there are problems in the monitored node. The alert conditions vary for each plugin. 


