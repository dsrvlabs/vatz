# VATZ Project Design (v1)

![Vatz Project Design](https://user-images.githubusercontent.com/6308023/210964352-d7c59b6e-9219-4d96-bf78-ccf2a21d1be6.png)

> **VATZ** is mainly designed to check the node status in real time and get the alert notification of all blockchain protocols, including metrics that doesn't supported by the protocol itself. Features for helping node operators such as automation that enable node manage orchestration and controlling VATZ over CLI commands are going to be available in near future.

### **VATZ** Project consists of 3 major components for followings: <br>
(Will be upgraded or added for the future)
1. [VATZ proto](https://github.com/dsrvlabs/vatz-proto) 
2. [VATZ Service](https://github.com/dsrvlabs/vatz)
3. VATZ Plugins (Official)
   - [vatz-plugin-sysutil](https://github.com/dsrvlabs/vatz-plugin-sysutil)
   - [vatz-plugin-cosmoshub](https://github.com/dsrvlabs/vatz-plugin-cosmoshub)
### **VATZ** service supports extension to 3rd party apps for alert notifications & metric analysis.
4. Dispatchers
5. Monitoring (Metric Exporter)

---

### 1. VATZ-Proto Repository (gRPC protocols)

**VATZ** is a total node management tool that is designed to be customizable and expandable through plug-in and gRPC protocol from the initial design stage. End-users can develop their own plugins to add new features with their own needs regardless of the programming language by using gRPC protocol.

### 2. VATZ Service

- This is a main service of **VATZ** that executes plugin APIs based on configs.

```
SAMPLE DEFAULT YAML
---
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

`vatz_protocol_info` & `plugins_infos` must be declared in default.yaml to get started with **VATZ** properly.

### 3. Plugins
**VATZ** Plugins are designed with maximum flexibility and expandability to perform following actions for Blockchain protocols nodes
- `Check`: Node & Machine status
- `Collect`: Node's metric + more
- `Execute`: Command on machine for certain behaviors (e.g, Restart Node)
- `Automation`: Node operation orchestration  
- `more` : You can develop your own plugins to manage your own nodes. 


### 4. Dispatchers(Notification)
**VATZ** Supports 3rd party apps for notification to alert users when there's any trouble on hardware or blockchain metrics. 
- [Discord](https://discord.com/)
- [Pagerduty](https://www.pagerduty.com/)
- [Telegram Messenger](https://telegram.org/)

### 5. Monitoring
The blockchain protocols have so many unique logs, and it brings a lot of data which causes difficulties in finding meaningful data by standardizing it to make it easier to view 
and most of the validator teams have trouble managing logs from running nodes due to log's varieties. <br/>
**VATZ**'s monitoring service is designed to find a way to manage all logs from nodes efficiently with minimum efforts and cost.

- `Available Now`
   - [Prometheus - Grafana](https://prometheus.io/docs/visualization/grafana/)
- `Upcoming soon`
   - Elastic - Kibana
   - Google drive - Big Query


#### - **AS-IS**
![monitoring_as_is](https://user-images.githubusercontent.com/6308023/210969218-f9548c35-ff3d-456f-9c70-0c175dfb24c9.png)

`VATZ` currently supports sending metrics for followings for Prometheus: <br/>
(Note: More metrics will be available in the future) <br>
- VATZ:`service` Liveness
- VATZ:`plugins` Liveness

#### - **TO-BE**
![monitoring_to_be](https://user-images.githubusercontent.com/6308023/210969235-4aa505ee-28cc-4e16-8129-843dbc4f2ca0.png)

**VATZ** will support to more monitoring and analysis 3rd party apps tool as shown in the diagram above.

---

# VATZ Project Design (v2)
The comprehensive design of VATZ v2 is presently in the development stage, with an anticipated release scheduled for 2024. 
Should you have any inquiries or feedback regarding VATZ v2, do not hesitate to contact us.