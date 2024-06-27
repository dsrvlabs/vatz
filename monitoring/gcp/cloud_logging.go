package gcp

import (
	"cloud.google.com/go/logging"
	"github.com/dsrvlabs/vatz/manager/config"
	tp "github.com/dsrvlabs/vatz/types"
	"github.com/dsrvlabs/vatz/utils"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"time"
)

const defaultRPCAddr = "http://localhost:19091"

var periodicCloudLog *gcpCloudLoggingEntry

type gcpCloudLoggingEntry struct {
	Protocol    string         `json:"protocol"`
	HostName    string         `json:"host_name"`
	PluginState tp.PluginState `json:"plugin_state"`
	Timestamp   string         `json:"timestamp"`
}

type cloudLogging struct {
	client           *logging.Client
	reminderSchedule []string
	reminderCron     *cron.Cron
}

func (gl *cloudLogging) Prep(cfg *config.Config) error {
	periodicCloudLog = setLogMessage(cfg.Vatz.ProtocolIdentifier, cfg.Vatz.NotificationInfo.HostName, tp.PluginState{})
	return nil
}

func (gl *cloudLogging) Process() error {
	for _, schedule := range gl.reminderSchedule {
		_, err := gl.reminderCron.AddFunc(schedule, func() {
			pluginStatus, pluginStatusErr := utils.GetPluginStatus(defaultRPCAddr)
			periodicCloudLog.PluginState = pluginStatus
			if pluginStatusErr != nil {
				log.Error().Str("module", "monitoring > gcp > cloud_logging").Msgf(" Execute(GetPluginStatus) error: %v", pluginStatusErr)
				return // Log the error and continue the function execution
			}
			err := gl.storeLog(periodicCloudLog)
			if err != nil {
				log.Error().Str("module", "monitoring > gcp > cloud_logging").Msgf(" Execute(Storing Log) error: %v", err)
			}
		})
		if err != nil {
			log.Error().Str("module", "monitoring > gcp > cloud_logging").Msgf("failed to add function to cron: %v", err)
			return err
		}
	}
	gl.reminderCron.Start()
	return nil
}

func (gl *cloudLogging) storeLog(logEntry *gcpCloudLoggingEntry) error {
	gcpLogger := gl.client.Logger(tp.MonitoringIdentifier)
	messageToSend := logging.Entry{
		Payload:  logEntry,
		Severity: logging.Info,
	}

	gcpLogger.Log(messageToSend)

	log.Info().Str("module", "monitoring").Msgf("Store Logs into Cloud logging for %s, %s", logEntry.Protocol, logEntry.HostName)
	return nil
}

func setLogMessage(protocol string, hostName string, ps tp.PluginState) *gcpCloudLoggingEntry {
	return &gcpCloudLoggingEntry{
		Timestamp:   time.Now().Format(time.RFC3339),
		Protocol:    protocol,
		HostName:    hostName,
		PluginState: ps,
	}
}
