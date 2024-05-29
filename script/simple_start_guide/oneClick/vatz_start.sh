#!/bin/bash
set -e
set -v

. /root/.vatz/.env_vatz

cd $VATZ_PATH

./vatz plugin start --plugin cpu_monitor --args "-port $CPU_MONITOR_PORT" --log $LOG_PATH/cpu_monitor.logs
./vatz plugin start --plugin mem_monitor --args "-port $MEM_MONITOR_PORT" --log $LOG_PATH/mem_monitor.logs
./vatz plugin start --plugin disk_monitor --args "-port $DISK_MONITOR_PORT" --log $LOG_PATH/disk_monitor.logs

./vatz plugin start --plugin node_block_sync --args "-port $NODE_BLOCK_SYNC_PORT" --log $LOG_PATH/node_block_sync.logs
./vatz plugin start --plugin node_is_alived --args "-port $NODE_IS_ALIVE_PORT" --log $LOG_PATH/node_is_alived.logs
./vatz plugin start --plugin node_peer_count --args "-port $NODE_PEER_COUNT_PORT" --log $LOG_PATH/node_peer_count.logs
./vatz plugin start --plugin node_active_status --args "-port $NODE_ACTIVE_STATUS_PORT -valoperAddr $VALOPER_ADDRESS" --log $LOG_PATH/node_active_status.logs
./vatz plugin start --plugin node_governance_alarm --args "-port $NODE_GOVERNANCE_ALARM_PORT -voterAddr $VOTER_ADDRESS" --log $LOG_PATH/node_governance_alarm.logs

./vatz start --config default.yaml >> /var/log/vatz/vatz.log 2>&1 &

echo "true"
