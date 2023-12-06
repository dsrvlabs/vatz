<h1 align="center"> VATZ (Validators' A To Z) </h1>  

<br/>  
<div align="center" style="display:flex;">  
  <img width="800" alt="v_black_square_800px" src="https://user-images.githubusercontent.com/6308023/220511474-d403b287-2e51-4bbb-a13f-9e0e5bef8b38.svg">

  <p> 
    <br>
    <img alt="Version"  src="https://img.shields.io/badge/version-v1.0.0-blue.svg?cacheSeconds=2592000"  />    
    <a href="https://www.gnu.org/licenses/lgpl-3.0.en.html"  target="_blank"><img alt="License: GPL 3.0"  src="https://img.shields.io/badge/License-GPL 3.0-yellow.svg" /></a> 
  </p> 
</div> 

# What is VATZ?

**VATZ** is a tool for building, analyzing, and managing blockchain node infrastructure safely and efficiently. You can set up **VATZ** to manage existing or new blockchain nodes and integrate with popular services like PagerDuty, Discord, or Telegram as well as custom in-house solutions.

## How does VATZ Project work?

This project is primarily designed to check node states in real time and receive alert notifications for all blockchain protocols, including metrics that the protocol itself might not support.

To this end, it consists of 3 components:

1. **VATZ-Proto**: API specification for VATZ (SVC) and VATZ Plugins that implements protobuf allowing users to develop in the language they wish.
2. **VATZ (SVC)**: Service that executes plugin APIs based on configs, checks plugin health, and sends notifications to configured channels.
3. **VATZ** **Plugins (with SDK)**: Plugins that integrate with VATZ to support features like node status checks, metric collection, and command execution.

For further information, check [VATZ Project Design](docs/design.md)


## What is the key feature of **VATZ Project**?

### Multi-Protocol Support
Any protocol can be managed through **VATZ** with plugins. Even unsupported protocols can be integrated through simple plugin development.

### Infrastructure as Code
**VATZ** is described using a high-level configuration syntax. You can divide your plugins into modular components that can then be combined in different ways to behave through automation. 

### Monitoring
Various logs and node data are collected and exported by a node exporter, then monitored through the 3rd party applications like Grafana. 

### Data Analysis
**VATZ** helps build datasets for your protocols, and transfer your data into popular services like Prometheus, Kafka, Google BigQuery, and more. In this way, **VATZ** aims to optimize your node infrastructure, and operators get insight into dependencies in their infrastructure. (Exp. 2023-Q3)

### Change Automation
Complex node operation tasks can be executed through **VATZ** with minimal human interaction. (Exp. 2023-Q4)

# Our Mission
We're on a mission to both transform the way people experience blockchain technology, and help them shape it. 
As validators, we provide tools for low-cost and low-effort node management for anyone wanting to onboard next-generation blockchain technology.

--- 
# Usage of VATZ

## How to get started with VATZ
Check out the [Installation Guide](docs/installation.md) to install and start using VATZ.
- You can get started with simple scripts, Please check [install scripts instructions](script/simple_start_guide/readme.md)

## How to use **VATZ** CLIs
Refer to the [VATZ CLIs guide](docs/cli.md) to find available CLI arguments.

## Official Plugins
> Our team is developing official plugins for easier operation, including basic monitoring metrics. 


### 1. [vatz-plugin-sysutil](https://github.com/dsrvlabs/vatz-plugin-sysutil)
vatz-plugin-sysutil is **VATZ** plugin for system utilization monitoring, i.e.:
- CPU
- DISK
- Memory
- Network Traffic

### 2. [vatz-plugin-comoshub](https://github.com/dsrvlabs/vatz-plugin-cosmoshub)
vatz-plugin-comoshub is **VATZ** plugin for cosmoshub node monitoring for followings:
- Node Block Sync
- Node Liveness
- Peer Count
- Active Status
- Node Governance Alarm

## Community Plugins
> We encourage everyone to share their plugins to make node operating easier.

- Please, share your own VATZ plugins on [Community Plugins](docs/community_plugins.md)!


---
# Release Note

Please check the Release Note for details of the latest releases.
- [VATZ](https://github.com/dsrvlabs/vatz/releases)
- [vatz-plugin-comoshub](https://github.com/dsrvlabs/vatz-plugin-cosmoshub/releases)
- [vatz-plugin-sysutil](https://github.com/dsrvlabs/vatz-plugin-sysutil/releases)

---

# We welcome your feedback
We're constantly striving to enhance and build on open-source resources.
Feel free to share your thoughts or feedback with us regarding VATZ.
You can start by registering any [issues](https://github.com/dsrvlabs/vatz/issues) you might find! <br>
Let’s continue building VATZ together! 

## Contributing

**VATZ** welcomes contributions! If you are looking to contribute, please check the following documents.
- [Contributing](docs/contributing.md) explains what kinds of contributions we look for and how to contribute.
- [Project Workflow Instructions](docs/workflow.md) explains how to build and test.

## License

The `VATZ` library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html), also
included in our repository in the `LICENSE.LESSER` file.

The `VATZ` binaries (i.e. all code inside of the `cmd` directory) are licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also
included in our repository in the `LICENSE` file.


# Contact us
Please don’t hesitate to contact [us](mailto:validator@dsrvlabs.com) if you need any further information about **VATZ**.

## Who are we

A leading blockchain technology company, **[DSRV](https://www.dsrvlabs.com/)** validates for 40+ global networks and provides infrastructure solutions for next-level building. This includes All That Node (enterprise-grade NaaS supporting 24+ protocols) and WELLDONE Studio (multi-chain product suite for developers and retail users alike).

Our ethos is to adapt to what the market and community need; our mission to advance the next internet and enable every player to build what they envision.

[<img alt="Homepage" src="https://user-images.githubusercontent.com/63234878/210315637-2d30efdd-5b9e-463e-8731-571916a6e1e3.svg" width="50" height="50" />](https://www.dsrvlabs.com/)
[<img alt="Medium" src="https://user-images.githubusercontent.com/6308023/176984456-f82c5c67-ebf3-455c-8494-c64ebfd66c58.svg" width="50" height="50" />](https://medium.com/dsrv)
[<img alt="Github" src="https://user-images.githubusercontent.com/6308023/176984452-c73aa188-563a-4b93-8ad8-cd7974770275.svg" width="50" height="50" />](https://github.com/dsrvlabs)
[<img alt="Youtube" src="https://user-images.githubusercontent.com/6308023/176984454-52c20db5-6b8f-4c15-a621-dd4a0052e99f.svg" width="50" height="50" />](https://www.youtube.com/channel/UCWhv8Kd430cEMpEYBPtSPjA/featured)
[<img alt="Twitter" src="https://user-images.githubusercontent.com/6308023/176984455-d48b24a9-1eb4-4c38-b728-2f4a0ccff09b.svg" width="50" height="50" />](https://twitter.com/dsrvlabs)
