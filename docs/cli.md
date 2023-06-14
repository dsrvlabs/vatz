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
  init        Init
  plugin      Plugin commands
  start       start VATZ
  version     VATZ Version

Flags:
  -h, --help   help for this command

Use " [command] --help" for more information about a command.
```


## Init
To start VATZ you need to init first. Currently, there are 4 flags under `init`. 
```
~$ ./vatz init --help
init

Usage:
   init [flags]

Flags:
  -a, --all             Create config yaml with all default setting of official plugins.
  -h, --help            help for init
  -p, --home string     Home directory of VATZ (default "~/.vatz")
  -o, --output string   New config file to create (default "default.yaml")
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


## Plugin

  **VATZ** binary also supports several plugin commands. In this document, usage of plugin command will be described. Currently, there are 7 subcommands under the plugin. 

  ```
  ~$ ./vatz plugin --help
  Plugin commands

  Usage:
    plugin [command]

  Available Commands:
    enable      Enabled or Disable plugin
    install     Install new plugin
    list        List installed plugin
    start       Start installed plugin
    status      Get statuses of Plugin
    stop        Stop running plugin
    uninstall   Uninstall plugin from plugin registry

  Flags:
    -h, --help   help for plugin

  Use " plugin [command] --help" for more information about a command.
  ```
  ### 1. enable
  You can temparaily enable or disable the plugin.
  ```
  ~$ ./vatz plugin enable <pluginName> <true/false>
  ```
  ### 2. install
  
  First, you can install the plugin by CLI. For installing the plugin, you have to know the repository URL where the plugin is implemented.
  For example, 
  ```
  ~$ ./vatz plugin install github.com/dsrvlabs/vatz-plugin-cosmoshub/plugins/node_active_status cosmos-status
  ```

  The last argument, `cosmos-status` is a simple name that is used for binary name on your machine. So, you could set the plugin name as desired.
  ### 3. list
  If you install some plugin, you can query installed plugins by `list` subcommand.

  ```
  ~$ ./vatz plugin list
  ```
  ```
  2023-05-26T14:40:10+09:00 INF Load Config default.yaml module=config
  2023-05-26T14:40:10+09:00 INF List plugins module=plugin
  2023-05-26T14:40:10+09:00 INF Create DB Instance module=db
  2023-05-26T14:40:10+09:00 INF List module=plugin
  2023-05-26T14:40:10+09:00 INF List module=db
  +---------------+------------+---------------------+----------------------------------------------------------------------+---------+
  | NAME          | IS ENABLED | INSTALL DATE        | REPOSITORY                                                           | VERSION |
  +---------------+------------+---------------------+----------------------------------------------------------------------+---------+
  | cosmos-status | true       | 2023-05-26 14:40:00 | github.com/dsrvlabs/vatz-plugin-cosmoshub/plugins/node_active_status | latest  |
  +---------------+------------+---------------------+----------------------------------------------------------------------+---------+
  ```
  ### 4. start
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
  ~$ ./vatz plugin start --plugin <pluginName>
  ```
  For some plugins that require arguments for binary, you can provide it with `--args` flag. It should be surrounded by quotes like as one string. For example,

  ```
  ~$ ./vatz plugin start --plugin cosmos-status --args "--valoperAddr=<YOUR VALOPER ADDRESS>" 
  ```

  ```
  2022-11-21T02:57:14Z INF Start plugin cosmos-status%!(EXTRA string=--valoperAddr=<HIDDEN>) module=plugin
  2022-11-21T02:57:14Z INF Start plugin cosmos-status module=plugin
  2022-11-21T02:57:14Z INF newReader module=db
  2022-11-21T02:57:14Z INF Create DB Instance module=db
  ```

### 5. status
When VATZ is running, you can check status of plugins(OK or FAIL).

  ```
  ~$  ./vatz plugin status
  ```

### 6. stop
You can stop running plugin.
  ```
  ~$ ./vatz plugin stop --plugin <pluginName>
  ```
  ```
  2023-05-26T15:22:30+09:00 INF Load Config default.yaml module=config
  ***** Plugin status *****
  1: cpu_monitor [FAIL]
  ```
### 7. uninstall
You can uninstall the plugin.

  ```
  ~$ ./vatz plugin uninstall cosmos-status
  ```
  ```
  2023-05-26T15:17:12+09:00 INF Load Config default.yaml module=config
  2023-05-26T15:17:12+09:00 INF Uninstall a plugin cosmos-status from /Users/user/.vatz module=plugin
  2023-05-26T15:17:12+09:00 INF Create DB Instance module=db
  2023-05-26T15:17:12+09:00 INF List module=plugin
  2023-05-26T15:17:12+09:00 INF Find Process cosmos-status module=plugin
  2023-05-26T15:17:16+09:00 INF Get cosmos-status module=plugin
  2023-05-26T15:17:16+09:00 INF Get cosmos-status module=db
  2023-05-26T15:17:16+09:00 INF DeletePlugin module=db
  ```
