#!/bin/bash
set -e
set -v

DEFAULT_LOG_PATH=/var/log/vatz
DEFAULT_VATZ_PATH=/root/vatz

# Create vatz log folder
## Check if $LOG_PATH exists
if [ -d "$DEFAULT_LOG_PATH" ]; then
  echo "$DEFAULT_LOG_PATH already exists. Skipping creation."
else
## Create $LOG_PATH if it doesn't exist
  mkdir $DEFAULT_LOG_PATH
fi

# Compile VATZ
cd $DEFAULT_VATZ_PATH
make

## You will see binary named vatz

# Initialize VATZ
./vatz init

# Install vatz-plugin-sysutil
./vatz plugin install github.com/dsrvlabs/vatz-plugin-sysutil/plugins/cpu_monitor cpu_monitor
./vatz plugin install github.com/dsrvlabs/vatz-plugin-sysutil/plugins/mem_monitor mem_monitor
./vatz plugin install github.com/dsrvlabs/vatz-plugin-sysutil/plugins/disk_monitor disk_monitor

# Install vatz-plugin-cosmoshub
./vatz plugin install github.com/dsrvlabs/vatz-plugin-cosmoshub/plugins/node_block_sync node_block_sync
./vatz plugin install github.com/dsrvlabs/vatz-plugin-cosmoshub/plugins/node_is_alived node_is_alived
./vatz plugin install github.com/dsrvlabs/vatz-plugin-cosmoshub/plugins/node_peer_count node_peer_count
./vatz plugin install github.com/dsrvlabs/vatz-plugin-cosmoshub/plugins/node_active_status node_active_status
./vatz plugin install github.com/dsrvlabs/vatz-plugin-cosmoshub/plugins/node_governance_alarm node_governance_alarm

# Check plugin list
./vatz plugin list
