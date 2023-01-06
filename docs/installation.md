# To get started with VATZ 

> This instruction is mainly based on Linux:Ubuntu

## 1. VATZ 
1. Clone VATZ repository
    ```
    $ git clone git@github.com:dsrvlabs/vatz.git
    ```
2. Compile VATZ
    ```
    $ cd vatz
    $ make
    ```
    You will see binary named `vatz`

  
3. Initialize VATZ
    ```
    $ ./vatz init
    ```
    You will see config file named default.yaml once you initialize VATZ.

## 2. VATZ Plugin

1. Copy the plugin address
  > **Plugin Repositories**
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

  We officially provide 3 plugins for System Utility and 5 plugins for CosmosHub node. But you can easily develop custom plugins for the feature that you want with [VATZ SDK](https://github.com/dsrvlabs/vatz/tree/main/sdk).

2. Install the plugin

  ```
   $ ./vatz plugin install <plugin_address> <name>
  ```
  Put git address of the plugin you want to install. You can set the plugin name as desired.

For example,
  ```
  $ ./vatz plugin install https://github.com/dsrvlabs/vatz-plugin-sysutil/plugins/cpu_monitor cpu_monitor
  ```

3. Check installation success

  ```
  $ ./vatz plugin list
  ```

  Verify that the plugin is installed successfully by checking whether the plugin name is added to the list.


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

1. `port`: check your machine used port number. If you set it to a port number that is already in use, an error will occur.
2. `health_checker_schedule`: Use cron expression to set the time to check if vatz is alive. 
3. `dispatch_channels`: Add infos to receive alerts through discord, pagerduty, and telegram. 
4. `reminder_schedule`: Use cron expression to set the time to resend the alert if it is not confirmed.


### Update Plugin Infos

1. `plugin_name`: Refer to the plugin name you installed.
2. `plugin_port`: Set your plugin used port number to connect the plugin.
3. `verify_interval`: You can set the interval to check if the plugin is connected. The default interval is 15s.
4. `execute_interval`: You can set the interval to execute the plugin method. The default interval is 30s.


## 2. Start VATZ and VATZ Plugin
1. Start VATZ
    ```
    $ ./vatz start
    ```
    If you see following logs you successfully started VATZ service.
    ```
    root@validator-node:~/vatz# go run main.go
    2022/04/22 05:06:54 Initialize Servers: VATZ Manager
    2022/04/22 05:06:54 Listening Port :9090
    2022/04/22 05:06:54 Node Manager Started
    ```

2. Start VATZ plugin
    ```
    $ ./vatz plugin start <name>
    ```
    Put the plugin name you want to start. You can search installed plugins with CLI command.

    For example, 

    ```
    $ ./vatz plugin start cpu_monitor
    ```

    If you see the following commands, you successfully started VATZ plugin.

    ```
    2022-09-14T08:17:33+02:00 INF Register module=grpc
    2022-09-14T08:17:33+02:00 INF Start 127.0.0.1 9094 module=sdk
    2022-09-14T08:17:33+02:00 INF Start module=grpc
    2022-09-14T08:17:48+02:00 INF Execute module=grpc
    ```    

3. The alert notification will be sent to configured channels if there are problems in the monitored node. The alert conditions vary for each plugin.
