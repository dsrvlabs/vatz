package config

// Default contents
const configDefaultContents = `
vatz_protocol_info:
  protocol_identifier: "vatz"
  port: 9090
  notification_info:
    discord_secret: "XXXXX"
    pager_duty_secret: "YYYYY"
  health_checker_schedule:
    - "* 1 * * *"
plugins_infos:
  default_verify_interval:  15
  default_execute_interval: 30
  default_plugin_name: "vatz-plugin"
  plugins:
    - plugin_name: "vatz-plugin-node-checker"
      plugin_address: "localhost"
      verify_interval: 7
      execute_interval: 9
      plugin_port: 9091
      executable_methods:
        - method_name: "isUp"
        - method_name: "getBlockHeight"
        - method_name: "getNumberOfPeers"
    - plugin_name: "vatz-plugin-machine-checker"
      plugin_address: "localhost"
      verify_interval: 8
      execute_interval: 10
      plugin_port: 9092
      executable_methods:
        - method_name: "getMemory"
        - method_name: "getDiscSize"
        - method_name: "getCPUInfo"
`

// "verify_interval", "execute_interval" and "plugin_name" on "plugins" are removed.
const configNoIntervalContents = `
vatz_protocol_info:
  protocol_identifier: "vatz"
  port: 9090
  notification_info:
    discord_secret: "hello"
    pager_duty_secret: "world"

plugins_infos:
  default_verify_interval:  15
  default_execute_interval: 30
  default_plugin_name: "vatz-plugin"
  plugins:
    - plugin_address: "localhost"
      plugin_port: 9091
      executable_methods:
        - method_name: "isUp"
        - method_name: "getBlockHeight"
        - method_name: "getNumberOfPeers"
`

// Intentionally ruin file contents.
const configInvalidYAMLContents = `
vatz_protocol_info
  protocol_identifier: "vatz"
  port: 9090
  "notification_info":
    discord_secret: "hello"
    pager_duty_secret: "world"

plugins_infos:
  default_verify_interval:  15
  default_execute_interval: 30
  default_plugin_name: "vatz-plugin"
  plugins:
    - plugin_address: "localhost"
      plugin_port: 9091
      executable_methods:
        - method_name: "isUp"
        - method_name: "getBlockHeight"
        - method_name: "getNumberOfPeers"
`
