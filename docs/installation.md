# To get started with VATZ 

> This instruction is mainly based on Linux:Ubuntu
## 0. Prerequisites 
Install Go 
> You can skip this part if you already have installed go version > 1.18+ on your system. 

Please, check [latest version](https://github.com/dsrvlabs/vatz/releases) or check on [VATZ module](https://github.com/dsrvlabs/vatz/blob/main/go.mod) and follow these [instructions](https://go.dev/doc/install) to install. 


## 1. VATZ 
### 1.1. Clone VATZ repository

  ```
  $ git clone git@github.com:dsrvlabs/vatz.git
  ```

### 1.2. Compile VATZ

  ```
  $ cd vatz
  $ make
  ```
You will see binary named `vatz`

  
### 1.3. Initialize VATZ
```
$ ./vatz init
```
You will see config file named default.yaml once you initialize VATZ.

```
$ ./vatz init --all
```
You can also use `--all` flag to add all default setting of [official plugins](https://github.com/dsrvlabs/vatz/blob/main/docs/installation.md#2-vatz-plugin) to config file. For more details, please check [cli.md](https://github.com/dsrvlabs/vatz/blob/main/docs/cli.md) file or use `--help` flag. 

## 2. VATZ Plugin

### 2.1. Copy the plugin address

  We officially provide 4 plugins for System Utility and 5 plugins for CosmosHub node. 
  > **Official Plugin Repositories**
  > - [vatz-plugin-sysutil](https://github.com/dsrvlabs/vatz-plugin-sysutil) 
  >   - vatz_cpu_monitor
  >   - vatz_mem_monitor
  >   - vatz_disk_monitor
  >   - vatz_net_monitor
  > - [vatz-cosmos-hub](https://github.com/dsrvlabs/vatz-plugin-cosmoshub)
  >   - vatz_block_sync
  >   - vatz_node_is_alived
  >   - vatz_peer_count
  >   - vatz_active_status
  >   - vatz_gov_alarm

 But you can easily develop custom plugins for the feature that you want with [VATZ SDK](https://github.com/dsrvlabs/vatz/tree/main/sdk). You could fine community provided plugins in [community_plugins.md](https://github.com/dsrvlabs/vatz/blob/main/docs/community_plugins.md). Feel free to add your custom plugins to community_plugins.md to share with others.


### 2.2. Install the plugin

  ```
   $ ./vatz plugin install <plugin_address> <name>
  ```
  Put git address of the plugin you want to install. You can set the plugin name as desired.

For example,
  ```
  $ ./vatz plugin install https://github.com/dsrvlabs/vatz-plugin-sysutil/plugins/cpu_monitor cpu_monitor
  ```

### 2.3. Check installation success

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

# 4. Set VATZ configs for your own use

## 4.1. Modify default.yaml (VATZ)
There's a YAML file for configuration, and you'll need to update the default settings yourself to properly set up the VATZ service and clients.
sample:
> This is default.yaml, please update secrets, and etc for your node use.
```
vatz_protocol_info:
  home_path: "~/.vatz"
  protocol_identifier: "Put Your Protocol here"
  port: 9090
  health_checker_schedule:
    - "0 1 * * *"
  notification_info:
    host_name: "Put your machine's host name"
    default_reminder_schedule:
      - "*/30 * * * *"
    dispatch_channels:
      - channel: "discord"
        secret: "Put your Discord Webhook"
      - channel: "pagerduty"
        secret: "Put your PagerDuty's Integration Key (Events API v2)"
      - channel: "slack"
        secret: "Put Your Slack Webhook url"
        subscriptions:
          - "Please, put Plugin Name that you specifically subscribe to send notification."
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
    gcp:
      gcp_cloud_logging_info:
        enabled: true
        cloud_logging_credential_info:
          project_id: "Please, Set your GCP Project id"
          credentials_type: "Check the Credential Type: ADC: Application, SAC: Default Credentials, Service Account Credentials, APIKey: API Key, OAuth: OAuth2"
          credentials: "Put your credential Info"
          checker_schedule:
            - "* * * * *"
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
      executable_methods:
        - method_name: "sampleMethod1"
    - plugin_name: "samplePlugin2"
      plugin_address: "localhost"
      verify_interval: 7
      execute_interval: 9
      plugin_port: 10002
      executable_methods:
        - method_name: "sampleMethod2"
                                                                               
```

### 4.1.1. update VATZ protocol Infos

1. `port`: Check your machine used port number. If you set it to a port number that is already in use, an error will occur.
2. `health_checker_schedule`: Use cron expression to set the time to check if vatz is alive. 
3. `notificaition_info`:
   - `dispatch_channels`: Add infos to receive alerts through discord, pagerduty, and telegram. 
   - `reminder_schedule`: Use cron expression to set the time to resend the alert if it is not confirmed.
4. `rpc_info`: Check address and port to connect with plugins.
5. `monitoring_info`: Check address and port to use monitoring. Currently, VATZ supports monitoring only through Prometheus-Grafana.


### 4.1.2. Update Monitoring info

#### 1. GCP Integration
> GCP: VATZ can be integrated with Google Cloud Platform (GCP) to send specific data, such as liveness information and plugin statuses, to Google Cloud Logging.
- **``gcp``**
  - **`gcp_cloud_logging_info`**: Manages the integration of VATZ with GCP Cloud Logging.
    - **`enabled`**: Set to `true` or `false` to enable or disable this feature.
    - **`project_id`**: Specify your GCP `project_id` to enable the logging service.
    - **`credentials_type`**: Define the type of credentials used. Supported options:
      - 1. `ADC`: Application Default Credentials.
      - 2. `SAC`: Service Account Credentials.
      - 3. `APIKey`: API Key.
      - 4. `OAuth`: OAuth2.
    - **`credentials`**: Provide the path to your credentials or include the key directly. If a URL is provided, the credentials will be downloaded to memory, not stored on disk.
    - **`checker_schedule`**: Set the interval for periodic checks, in seconds.

#### 2. Prometheus Integration
> VATZ can collect its own metrics using Prometheus which helps you to manage VATZ by Grafana.
- **``prometheus``**  
  - **`enabled`**: Set to `true` or `false` to enable or disable this feature.
  - **`address`**: Specify the address where Prometheus will run.
  - **`port`**: Define the port number for the Prometheus service.


### 4.1.3. Update Plugin Infos

1. `plugin_name`: Refer to the plugin name you installed.
2. `plugin_port`: Set your plugin used port number to connect the plugin.
3. `verify_interval`: You can set the interval to check if the plugin is connected. The default interval is 15s.
4. `execute_interval`: You can set the interval to execute the plugin method. The default interval is 30s.


## 5. Start VATZ service and VATZ Plugin service

### 5.1. Start VATZ plugin

  ```
  $ ./vatz plugin start --plugin <name> --args <arguments> --log <logfile>
  ```
  Put the plugin name you want to start.

  ```
  $ ./vatz plugin start -p cpu_monitor
  ```

  Alternatively, you can start the plugin with some arguments with `-arg` or `-a` flag. For example, you could change port number as below.
  Some plugins might require arguments to start. Check the plugin repository for available arguments or requirements.

  ```
  $ ./vatz plugin start -p cpu_monitor -a "-port=9094 -urgent=80"
  ```
  If you see the following commands, you successfully started VATZ plugin.
  ```
  2024-01-17T01:01:32-06:00 INF Start plugin cpu_monitor  module=plugin
  2024-01-17T01:01:32-06:00 INF Plugin cpu_monitor is successfully started. module=plugin
  ```
### 5.2. Start VATZ
  
  ```
  ~$ ./vatz start
  ```
  If you see following logs you successfully started VATZ service.

  ```
  ~$ ./vatz start
  2024-01-17T01:01:38-06:00 INF Initialize Server module=main
  2024-01-17T01:01:38-06:00 INF Start VATZ Server on Listening Port: :9090 module=main
  2024-01-17T01:01:38-06:00 INF Client successfully connected to localhost:9001 (plugin:cpu_monitor). module=util
  2024-01-17T01:01:38-06:00 INF start metric server: 127.0.0.1:18080 module=main
  2024-01-17T01:01:38-06:00 INF start rpc server module=rpc
  2024-01-17T01:01:38-06:00 INF start gRPC gateway server 127.0.0.1:19091 module=rpc
  2024-01-17T01:01:38-06:00 INF start gRPC server 127.0.0.1:19090 module=rpc
  2024-01-17T01:01:38-06:00 INF Client successfully connected to localhost:9001 (plugin:cpu_monitor). module=util
  2024-01-17T01:01:46-06:00 INF Executor send request to cpu_monitor module=executor
  2024-01-17T01:01:46-06:00 INF response: SUCCESS module=executor
  ```
  You can set your own path for config file such as 
  ```
  ~$ ./vatz start --config <path_to_your_own_file_at/config_file.yaml>
  ```

### 5.3. Notification
Alert notifications are sent to the configured channels when issues are detected in the monitored node. The specific alert conditions depend on each plugin.

> Dispatcher channels can be configured to operate in two modes
> - **Non-subscription mode**
> - **Subscription-based mode**

Below are examples of both modes, along with the list of installed plugins:
```
 vatz plugin list                                                                                                                                                             
+-------------------+------------+---------------------+-----------------------------------------------------------------+---------+
|     NAME          | IS ENABLED | INSTALL DATE        | REPOSITORY                                                      | VERSION |
+-------------------+------------+---------------------+-----------------------------------------------------------------+---------+
| cosmos_watcher    | true       | 2024-08-28 18:50:38 | github.com/dsrvlabs/vatz-plugin-watchers/plugins/watcher_cosmos | latest  |
| shentu_watcher    | true       | 2024-08-28 18:50:44 | github.com/dsrvlabs/vatz-plugin-watchers/plugins/watcher_cosmos | latest  |
| osmosis_watcher   | true       | 2024-08-28 18:50:44 | github.com/dsrvlabs/vatz-plugin-watchers/plugins/watcher_cosmos | latest  |
+-------------------+------------+---------------------+-----------------------------------------------------------------+---------+
```


1. **Non-subscription mode**: In this mode, notifications are automatically delivered from any declared plugin that meets its own alert conditions. All registered dispatcher channels will send alert messages from all plugins.
```
    dispatch_channels:
      - channel: "discord"
        secret: "https://discord.com/api/webhooks/221449818154281/7PfPpuWv4uK4wkPp-uWT1nJalAesD0YgSZA2j2EL7YvAN1ah32"
        subscriptions:
      - channel: "pagerduty"
        secret: "y_NbAkKc66ryYTWUX4YEu801s"
```
  In this mode, there is no need to set any specific subscriptions. If any plugin’s alert conditions are met, the alert messages will be sent to both Discord and PagerDuty.

2. **Subscription-based mode**: This mode allows you to subscribe to specific plugins, ensuring that you receive notifications only from those you select. In this case, only alert messages from your subscribed plugins will be delivered through the dispatcher channel.
```
    dispatch_channels:
      - channel: "discord"
        secret: "https://discord.com/api/webhooks/221449818154281/7PfPpuWv4uK4wkPp-uWT1nJalAesD0YgSZA2j2EL7YvAN1ah32"
        subscriptions:
      - channel: "pagerduty"
        secret: "y_NbAkKc66ryYTWUX4YEu801s"
        subscriptions:
         - "cosmos_watcher"
         - "shentu_watcher"
```
  In the example above, only alert messages from the `cosmos_watcher` and `shentu_watcher` plugins will be sent to PagerDuty. Alerts from the `osmosis_watcher` plugin will be ignored for Pagerduty channel.
    
