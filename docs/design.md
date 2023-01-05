# VATZ Project Design

![Vatz Project Design](https://user-images.githubusercontent.com/6308023/179885451-6d40505b-8b31-41d3-8dff-25220e00be1c.png)

> **VATZ** is mainly designed to check the node status in real time and get the alert notification of all blockchain protocols, including metrics that doesn't supported by the protocol itself. Features for helping node operators such as automation that enable node manage orchestration and controlling VATZ over CLI commands are going to be available in near future.

3 major services are planned in **VATZ** project as follows: <br>
(Will be upgraded or added for the future)

   1. Manager
   2. SDK
   3. Monitoring

---

## Proto Repository (gRPC protocols)

**VATZ** is a total node management tool that is designed to be customizable and expandable through plug-in and gRPC protocol from the initial design stage. End-users can develop their own plugins to add new features with their own needs regardless of the programming language by using gRPC protocol.


## Protocol Node

### 1. Manager

- This is a main service of **VATZ** that executes plugin APIs based on configs.

```
SAMPLE DEFAULT YAML
---
vatz_protocol_info:
  protocol_identifier: "VATZ"
  port: 9090
  notification_info:
    discord_secret: "xxxxxxx"
plugins_infos:
  default_verify_interval: 15
  default_execute_interval: 30
  default_plugin_name: "vatz-plugin"
  plugins:
    - plugin_name: "sample1"
      plugin_address: "localhost"
      plugin_port: 9091
      executable_methods:
        - method_name: "sampleMethod1"
    - plugin_name: "sample2"
      plugin_address: "localhost"
      verify_interval: 7
      execute_interval: 9
      plugin_port: 10002
      executable_methods:
        - method_name: "sampleMethod2"
```

`vatz_protocol_info` & `plugins_infos` must be declared in default.yaml to get started with **VATZ** properly.

### 2. Plugins

Plugins that allow you to perform followings per protocols
   - `Check`: Node & Machine status
   - `Collect`: Node's metric + more
   - `Execute`: Command on machine for certain behaviors (e.g, Restart Node)
   

### 3. Monitoring
The blockchain protocols have so many unique logs, and it brings a lot of data which causes difficulties in finding meaningful data by standardizing it to make it easier to view.<br/>
The most validator teams have trouble managing logs from running nodes due to log's varieties.<br/>
**VATZ**'s monitoring service is designed to find a way to manage all logs from nodes efficiently with minimum cost.<br/>

**AS-IS** <br/>
![monitoring-as-is](../docs/images/monitoring_as_is.png)

`VATZ` currently supports sending metrics for followings for Prometheus: <br/>

- VATZ:`service` liveness
- VATZ:`plugins` liveness

More metrics will be available in the future.

**TO-BE**<br/>
![monitoring-as-is](../docs/images/monitoring_to_be.png)

We can also add monitoring services as shown in the diagram above.

## 3rd Party Applications
We are trying to provide functions that can be easily integrated with the 3rd party applications most of the Validator teams are currently using now.
 
> 1. [Grafana](https://grafana.com/)
> 2. [Kibana](https://www.elastic.co/)
> 3. (TBD)
