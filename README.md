<h1 align="center"> VATZ (Validators' A To Z) </h1>  

<br/>  
<div align="center" style="display:flex;">  
  <img width="400" alt="v_black_square_800px" src="https://user-images.githubusercontent.com/63234878/209922215-cccf7b88-de12-42ac-9714-8a5e3f194d7f.png">

  <p> 
    <br>
    <img alt="Version"  src="https://img.shields.io/badge/version-v1.0--beta--rc.1-blue.svg?cacheSeconds=2592000"  />    
    <a href="https://www.apache.org/licenses/LICENSE-2.0"  target="_blank"><img alt="License: Apache 2.0"  src="https://img.shields.io/badge/License-Apache 2.0-yellow.svg" /></a> 
  </p> 
</div> 

# What is VATZ?

**VATZ** is a tool for building, analyzing, and managing blockchain node infrastructure safely and efficiently. You can set up VATZ to manage existing or new blockchain nodes and integrate with popular services PagerDuty, Discord, Telegram and more as well as custom in-house solutions.

## How does **VATZ Project** work?

VATZ project is primarily designed to check the node states in real-time and receive alert notifications of all blockchain protocols, including metrics that the protocol itself does not support. 

To this end, it consists of 3 components:

1. VATZ : VATZ executes plugin APIs based on configs, checks the health of plugins, and sends notifications to configured channels
2. Plugins (SDK): Various features like checking node status, collecting node metrics and executing certain commands, can be integrated to VATZ through separate plugins. 
3. Monitoring: Various logs and data of nodes are exported by a node exporter and monitored through the 3rd party applications like Grafana. 

For further information, check [VATZ Project Design](docs/design.md)


## What is the key feature of **VATZ Project**?

  ### Multi Protocol Support
   **VATZ** is NOT limited Protocol Type where it categorizes on chain protocol. Any Protocol can be managed through VATZ with plugins, even unsupported protocols can be integrated through simple plugin development.
  ### Infrastructure as Code
  **VATZ** is described using a high-level configuration syntax. You can divide your plugins into modular components that can then be combined in different ways to behave through automation.
  ### Data Analysis
   **VATZ** helps to build datasets for your managing protocols and transfer your data into popular services Prometheus, Kafaka, Google BigQuery and more. Because of this, VATZ aims to set your Node infrastructure as efficiently as possible, and operators get insight into dependencies in their infrastructure. (2023-Q3)
  ### Change automation
  Complex sets of node's operational tasks can be done through **VATZ** with minimal human interaction. (2023-Q4)
  

--- 
# Usage of VATZ

## How to get started with **VATZ**?
Please follow [Installation guide](docs/installation.md) to install and start VATZ.

## How to use **VATZ** CLIs?
Please check [VATZ CLIs guide](./docs/cli.md) to find available CLI arguments.

## Official Plugins
> We are developing official plugins together for easier operation including basic monitoring metrics.

### 1. [vatz-plugin-sysutil](https://github.com/dsrvlabs/vatz-plugin-sysutil)
vatz-plugin-sysutil is **VATZ** plugin for system utilization monitoring such as 
- CPU
- DISK
- Memory

### 2. [vatz-plugin-comoshub](https://github.com/dsrvlabs/vatz-plugin-cosmoshub)
vatz-plugin-comoshub is **VATZ** plugin for cosmoshub node monitoring for followings:
- Node Block Sync
- Node Liveness
- Peer Count
- Active Status
- Node Governance Alarm

---
# Release Note

Please check the Release Note for details of the latest releases.
- [VATZ](https://github.com/dsrvlabs/vatz/releases)
- [vatz-plugin-comoshub](https://github.com/dsrvlabs/vatz-plugin-cosmoshub/releases)
- [vatz-plugin-sysutil](https://github.com/dsrvlabs/vatz-plugin-sysutil/releases)

# Our Mission

We're on a mission to transform the way people experience blockchain technology and let them contribute and become a part of its technology.
As Validators, we provide tools to people to manage their own nodes with low cost and less effort for anyone who would like to join future blockchain technology.

---

# Feel free to share your feedback
We're constantly striving to make better open-source all together.
Please, share your thoughts or any feedback regarding **VATZ** Project.
You can start with registering an [issue](https://github.com/dsrvlabs/vatz/issues), if there's one you think. <br>
Contribute to **VATZ** project too!!

## Contributing

**VATZ** welcomes contributions! If you are looking to contribute, please check the following documents.
- [Contributing](docs/contributing.md) explains what kinds of contributions we welcome and how to contribute.
- [Project Workflow Instructions](docs/workflow.md) explains how to build and test.


# Contact us
Please, contact [us](mailto:validator@dsrvlabs.com) if you need any further information about **VATZ**.

**DSRV** is a blockchain infrastructure company that provides powerful and easy-to-use solutions to enable developers and enterprises to become architects of the future. Visit [DSRV](https://www.dsrvlabs.com/), if you are interested in various products we build for the Web 3.0 developers.

[<img alt="Homepage" src="https://user-images.githubusercontent.com/63234878/210315637-2d30efdd-5b9e-463e-8731-571916a6e1e3.svg" width="50" height="50" />](https://www.dsrvlabs.com/)
[<img alt="Medium" src="https://user-images.githubusercontent.com/6308023/176984456-f82c5c67-ebf3-455c-8494-c64ebfd66c58.svg" width="50" height="50" />](https://medium.com/dsrv)
[<img alt="Github" src="https://user-images.githubusercontent.com/6308023/176984452-c73aa188-563a-4b93-8ad8-cd7974770275.svg" width="50" height="50" />](https://github.com/dsrvlabs)
[<img alt="Youtube" src="https://user-images.githubusercontent.com/6308023/176984454-52c20db5-6b8f-4c15-a621-dd4a0052e99f.svg" width="50" height="50" />](https://www.youtube.com/channel/UCWhv8Kd430cEMpEYBPtSPjA/featured)
[<img alt="Twitter" src="https://user-images.githubusercontent.com/6308023/176984455-d48b24a9-1eb4-4c38-b728-2f4a0ccff09b.svg" width="50" height="50" />](https://twitter.com/dsrvlabs)
