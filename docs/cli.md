# How to use CLI of VATZ

## Init and start VATZ

Visit [Installation](./installation.md).


## Help
For more details, you can query helps by adding `--help` flag.
```
~$ ./vatz --help
Usage:
   [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        init
  plugin      Plugin commands
  start       start VATZ
  stop        stop VATZ
  version     VATZ Version

Flags:
      --debug   Enable debug mode on Log.
  -h, --help    help for this command
      --trace   Enable Trace mode on Log.

Use " [command] --help" for more information about a command.
```


## Init
To start VATZ you need to init first. Currently, there are 4 flags under `init`. 
```
Usage:
   init [flags]

Flags:
  -a, --all             Create config yaml with all default setting of official plugins.
  -h, --help            help for init
  -p, --home string     Home directory of VATZ (default "~/.vatz")
  -o, --output string   New config file to create (default "default.yaml")

Global Flags:
      --debug   Enable debug mode on Log.
      --trace   Enable Trace mode on Log.
```
With `--all` flag, you create config file including all default setting of official plugins like below. If you use the default settings, make sure to name the plugin when you install it to match the <pluginName> of the setting.
```
  default_verify_interval: 15
  default_execute_interval: 30
  default_plugin_name: "vatz-plugin"
  plugins:
    - plugin_name: "vatz_cpu_monitor"
      plugin_address: "localhost"
      plugin_port: 9001
      executable_methods:
        - method_name: "cpu_monitor"
    - plugin_name: "vatz_mem_monitor"
      plugin_address: "localhost"
      plugin_port: 9002
      executable_methods:
        - method_name: "mem_monitor"
    - plugin_name: "vatz_disk_monitor"
      plugin_address: "localhost"
      plugin_port: 9003
      executable_methods:
        - method_name: "disk_monitor"
    - plugin_name: "vatz_net_monitor"
      plugin_address: "localhost"
      plugin_port: 9004
      executable_methods:
        - method_name: "net_monitor"
    - plugin_name: "vatz_block_sync"
      plugin_address: "localhost"
      plugin_port: 10001
      executable_methods:
        - method_name: "node_block_sync"
    - plugin_name: "vatz_node_is_alived"
      plugin_address: "localhost"
      plugin_port: 10002
      executable_methods:
        - method_name: "node_is_alived"
    - plugin_name: "vatz_peer_count"
      plugin_address: "localhost"
      plugin_port: 10003
      executable_methods:
        - method_name: "node_peer_count"
    - plugin_name: "vatz_active_status"
      plugin_address: "localhost"
      plugin_port: 10004
      executable_methods:
        - method_name: "node_active_status"
    - plugin_name: "vatz_gov_alarm"
      plugin_address: "localhost"
      plugin_port: 10005
      executable_methods:
        - method_name: "node_governance_alarm"`
```
Also, You can set the home directory of VATZ (default "~/.vatz") with `--home` flag.
And you can create the new config file with your disired name with `--output` flag

## Start
You can start Vatz when you finish set up config default.yaml or your own config yaml file.
```
./vatz start --help                 
start VATZ

Usage:
   start [flags]

Flags:
      --config string       VATZ config file. (default "default.yaml")
  -h, --help                help for start
      --log string          log file export to.
      --prometheus string   prometheus port number. (default "18080")

Global Flags:
      --debug   Enable debug mode on Log.
      --trace   Enable Trace mode on Log.

```
You can simply start VATZ with command
```
~$ ./vatz start 
```
This will start VATZ with default config, that's been created from previous command `vatz init`. <br>
You can set exact config file with flag --config such as 
```
~$ ./vatz start --config /root/User/vatz-config.yaml
```
You need to use absolute path. 
You can also check your own log by its level, adding flag `--debug` or `--trace`.

## Stop
You can stop(kill gracefully) Vatz when you want to terminate vatz process. 
```shell
stop VATZ

Usage:
   stop [flags]

Flags:
  -h, --help   help for stop

Global Flags:
      --debug   Enable debug mode on Log.
      --trace   Enable Trace mode on Log.
```
You can simply stop VATZ with command
```
~$ ./vatz stop 
```
---
```
  ~$ ./vatz stop       
  2024-01-17T01:05:45-06:00 INF Sent termination signal to VATZ process, terminating ...
```
```
2024-01-17T01:05:33-06:00 INF Initialize Server module=main
2024-01-17T01:05:33-06:00 INF Start VATZ Server on Listening Port: :9090 module=main
2024-01-17T01:05:34-06:00 INF Client successfully connected to localhost:9001 (plugin:cpu_monitor). module=util
2024-01-17T01:05:34-06:00 INF start metric server: 127.0.0.1:18080 module=main
2024-01-17T01:05:34-06:00 INF start rpc server module=rpc
2024-01-17T01:05:34-06:00 INF start gRPC gateway server 127.0.0.1:19091 module=rpc
2024-01-17T01:05:34-06:00 INF start gRPC server 127.0.0.1:19090 module=rpc
2024-01-17T01:05:34-06:00 INF Client successfully connected to localhost:9001 (plugin:cpu_monitor). module=util
2024-01-17T01:05:42-06:00 INF Executor send request to cpu_monitor module=executor
2024-01-17T01:05:42-06:00 INF response: SUCCESS module=executor
2024-01-17T01:05:48-06:00 INF Received signal: interrupt module="cmd > start"
2024-01-17T01:05:48-06:00 INF Terminating VATZ...
```

## Plugin

  **VATZ** binary also supports several plugin commands. In this document, usage of plugin command will be described. Currently, there are 7 subcommands under the plugin. 

  ```
  ~$ ./vatz plugin --help
  Plugin commands
  
  Usage:
     plugin [command]
  
  Available Commands:
    disable     Disable plugin
    enable      Enable plugin
    install     Install new plugin
    list        List installed plugin
    start       Start installed plugin
    status      Get statuses of Plugin
    stop        Stop running plugin
    uninstall   Uninstall plugin from plugin registry
  
  Flags:
        --config string   VATZ config file. (default "default.yaml")
    -h, --help            help for plugin
  
  Global Flags:
        --debug   Enable debug mode on Log.
        --trace   Enable Trace mode on Log.
  
  Use " plugin [command] --help" for more information about a command.
  ```
  ### 1. disable
  You can temporarily disable an enabled plugin.
  VATZ does not execute plugin method when plugin is disabled. 
  ```
  ~$ ./vatz plugin disable -p <pluginName> 
  or
  ~$ ./vatz plugin disable --plugin <pluginName>
  ```
  ---
  #### 1.1. Example(`disable`)
  ```
  ~$ ./vatz plugin list  
  +-------------+------------+---------------------+-------------------------------------------------------------+---------+
  | NAME        | IS ENABLED | INSTALL DATE        | REPOSITORY                                                  | VERSION |
  +-------------+------------+---------------------+-------------------------------------------------------------+---------+
  | cpu_monitor | true       | 2023-12-29 00:12:08 | github.com/dsrvlabs/vatz-plugin-sysutil/plugins/cpu_monitor | latest  |
  +-------------+------------+---------------------+-------------------------------------------------------------+---------+
  ~$ ./vatz plugin disable -p cpu_monitor
  2024-01-16T23:11:35-06:00 INF Plugin cpu_monitor is disabled. module=db
  
  ~$ ./vatz plugin list                  
  +-------------+------------+---------------------+-------------------------------------------------------------+---------+
  | NAME        | IS ENABLED | INSTALL DATE        | REPOSITORY                                                  | VERSION |
  +-------------+------------+---------------------+-------------------------------------------------------------+---------+
  | cpu_monitor | false      | 2023-12-29 00:12:08 | github.com/dsrvlabs/vatz-plugin-sysutil/plugins/cpu_monitor | latest  |
  +-------------+------------+---------------------+-------------------------------------------------------------+---------+
  ```
  
  ### 2. enable
  You can temporarily enable a disabled plugin.
  ```
  ~$ ./vatz plugin enable -p <pluginName> 
  or
  ~$ ./vatz plugin enable --plugin <pluginName>
  ```
  ---
  #### 2.1. Example(`enable`)
  ```
  ~$ ./vatz plugin enable -p cpu_monitor
  2024-01-16T23:14:18-06:00 INF Plugin cpu_monitor is enabled. module=db
  
  ./vatz plugin list                 
  +-------------+------------+---------------------+-------------------------------------------------------------+---------+
  | NAME        | IS ENABLED | INSTALL DATE        | REPOSITORY                                                  | VERSION |
  +-------------+------------+---------------------+-------------------------------------------------------------+---------+
  | cpu_monitor | true       | 2023-12-29 00:12:08 | github.com/dsrvlabs/vatz-plugin-sysutil/plugins/cpu_monitor | latest  |
  +-------------+------------+---------------------+-------------------------------------------------------------+---------+

  ```

  ### 3. install
  
  First, you can install the plugin by CLI. For installing the plugin, you have to know the repository URL where the plugin is implemented.
  ```
  ~$ ./vatz plugin install <plugin's githubAddress> <pluginName>
  ```

  #### 3.1. Example(`install`)
  ```
  ~$ ./vatz plugin install github.com/dsrvlabs/vatz-plugin-cosmoshub/plugins/node_active_status cosmos-status
  or 
  ~$ ./vatz plugin install github.com/dsrvlabs/vatz-plugin-cosmoshub/plugins/node_active_status cosmos-status -v v1.0.0
  ```

  The last argument, `cosmos-status` is a simple name that is used for binary name on your machine. So, you could set the plugin name as desired.
  
  ### 4. list
  If you install a plugin, you can use the list subcommand to view the installed plugins. 
  In the previous example, we installed version v1.0.0, and the plugin list will display the exact version that was installed.

  ```
  ~$ ./vatz plugin list
  ```
  ```
  +---------------+------------+---------------------+----------------------------------------------------------------------+---------+
  | NAME          | IS ENABLED | INSTALL DATE        | REPOSITORY                                                           | VERSION |
  +---------------+------------+---------------------+----------------------------------------------------------------------+---------+
  | cosmos-status | true       | 2024-01-17 00:05:24 | github.com/dsrvlabs/vatz-plugin-cosmoshub/plugins/node_active_status | v1.0.0  |
  +---------------+------------+---------------------+----------------------------------------------------------------------+---------+
  ```

  ### 5. start
  There are 4 flags under `plugin start`

  ```
  ~$ ./vatz plugin start -h
  Start installed plugin

  Usage:
    plugin start [flags]

  Examples:
  vatz plugin start pluginName

  Flags:
    -a, --args string     Arguments
    -h, --help            help for start
    -l, --log string      Logfile
    -p, --plugin string   Installed plugin name
  ```

  You can start installed plugin like below.
  ```
  ~$ ./vatz plugin start --plugin <pluginName> or vatz plugin start --plugin <pluginName> --args <arguments>
  ```
  For certain plugins that necessitate arguments for the binary, these can be supplied using the --args flag. 
  The arguments should be enclosed within quotes as a single string. Details about the arguments are provided by the plugin itself. 
  For instance, in the case of the recently installed plugin, you can refer to https://github.com/dsrvlabs/vatz-plugin-cosmoshub for more information.
  
  ```
  ~$ ./vatz plugin start --plugin cosmos-status --args "--valoperAddr=5dsxaisdoifb2b194ajsllba7"
  ```

  ```
  2024-01-17T00:28:10-06:00 INF Start plugin cosmos-status --valoperAddr=5dsxaisdoifb2b194ajsllba7 module=plugin
  2024-01-17T00:28:10-06:00 INF Plugin cosmos-status is successfully started. module=plugin
  ```
  
  ### 6. status
  When VATZ is running, you can check status of plugins(OK or FAIL).<br>
  Note: You will get an error if VATZ isn't running. 
  ```
  ~$  ./vatz plugin status
  ```
  ---
  ```
  ./vatz plugin status
  ***** Plugin Status *****
  1: cosmos-status [OK]
  ```
  ### 7. stop
  You can stop running plugin.
  ```
  ~$ ./vatz plugin stop --plugin <pluginName> 
  or 
  ~$ ./vatz plugin stop -p <pluginName> 
  ```
  ---
  #### 7.1. Example(`stop`)
  ```
  ~$  ./vatz plugin stop --plugin cosmos-status      
  2024-01-17T00:35:23-06:00 INF Stop plugin cosmos-status module=plugin
  2024-01-17T00:35:25-06:00 INF Plugin cosmos-status is successfully stopped. module=plugin
  ~$  ./vatz plugin status
  ***** Plugin Status *****
  1: cosmos-status [FAIL]
  ```
  ### 8. uninstall
  You can uninstall the plugin.
  
  ```
  ~$ ./vatz plugin uninstall <pluginName>
  ```
  ---
  #### 8.1. Example(`uninstall`)
  ```
  ./vatz plugin uninstall cosmos-status
  2024-01-17T00:46:29-06:00 INF Plugin cosmos-status is successfully uninstalled from /Users/dongyookang/.vatz module=plugin
  ```
