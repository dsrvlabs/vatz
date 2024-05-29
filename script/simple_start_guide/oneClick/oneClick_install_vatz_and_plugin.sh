#!/bin/bash
set -e
set -v

. .env

# Create vatz log folder
mkdir $LOG_PATH

# Compile VATZ
cd $VATZ_PATH
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

# Copy default_config.yaml to /root/git/vatz/default.yaml
cp /root/git/vatz/script/simple_start_guide/oneClick/default_config.yaml /root/git/vatz/default.yaml

# Configure /root/git/vatz/default.yaml
wget https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -O /usr/bin/yq && sudo chmod +x /usr/bin/yq
yq -i ".vatz_protocol_info.notification_info.host_name = \"$(cat /etc/hostname)\"" /root/git/vatz/default.yaml
yq -i ".vatz_protocol_info.notification_info.dispatch_channels[0].secret = \"$DISCORD_WEBHOOK\"" /root/git/vatz/default.yaml
yq -i ".vatz_protocol_info.protocol_identifier = \"$PROTOCOL\"" /root/git/vatz/default.yaml

# Copy vatz start script
cp /root/git/vatz/script/simple_start_guide/oneClick/.env /root/.vatz/.env_vatz
cp /root/git/vatz/script/simple_start_guide/oneClick/vatz_start.sh /root/dsrv/bin/vatz_start.sh
cp /root/git/vatz/script/simple_start_guide/oneClick/vatz_stop.sh /root/dsrv/bin/vatz_stop.sh

bash /root/dsrv/bin/vatz_start.sh
