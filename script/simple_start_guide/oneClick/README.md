# How to install VATZ & Vatz Official plugins with scripts
> Those following scripts helps user install VATZ in simple way to get started fast. <br>
> Please, be advised that you need to update several variable in scripts to match with your current setting such as cloned VATZ path

## Here's simple instruction to get started
### 1. Go to script folder
- Users can simply setup and run VATZ with following scripts.

### 2. Setup .env file
- 
Users only need to input the protocol name, Discord webhook, and voter address in the .env file. The valoper address will be set via API call assuming the moniker name is "DSRV". The paths for the vatz and log are pre-configured with default settings.
```
PROTOCOL="Validator Node: <Protocol Name>"
DISCORD_WEBHOOK="webhook address"
VATZ_PATH=/root/git/vatz
LOG_PATH=/var/log/vatz

CPU_MONITOR_PORT=9001
MEM_MONITOR_PORT=9002
DISK_MONITOR_PORT=9003

NODE_BLOCK_SYNC_PORT=10001
NODE_IS_ALIVE_PORT=10002
NODE_PEER_COUNT_PORT=10003
NODE_ACTIVE_STATUS_PORT=10004
NODE_GOVERNANCE_ALARM_PORT=10005
VALOPER_ADDRESS=$(curl -s localhost:1317/cosmos/staking/v1beta1/validators | jq -r --arg moniker "DSRV" '.validators[] | select(.description.moniker == $moniker) | .operator_address')
VOTER_ADDRESS=voter address
```

### 3. To run the `oneClick_install_vatz_and_plugin.sh` script
- Simply executing this script will start the Vatz application. 
The script handles the installation of Vatz, plugins, configuration of the default.yaml file, and starts the Vatz service. 
Additionally, `vatz_start.sh` and `vatz_stop.sh` scripts will be copied to `/root/dsrv/bin`.
```
$ bash oneClick_install_vatz_and_plugin.sh
```


