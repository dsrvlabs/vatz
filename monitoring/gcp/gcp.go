package gcp

import (
	"cloud.google.com/go/logging"
	"context"
	"errors"
	"github.com/dsrvlabs/vatz/manager/config"
	tp "github.com/dsrvlabs/vatz/types"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"sync"
	"time"
)

type GCP interface {
	Prep(cfg *config.Config) error
	Process() error
}

var (
	gcpSingletons          []GCP
	gcpOnce                sync.Once
	validCredentialOptions = map[tp.CredentialOption]bool{
		tp.ApplicationDefaultCredentials: true,
		tp.ServiceAccountCredentials:     true,
		tp.APIKey:                        true,
		tp.OAuth2:                        true,
	}
)

func GetGCP(cfg config.MonitoringInfo) []GCP {
	gcpOnce.Do(func() {
		loggPrepInfo := cfg.GCP.GCPCloudLogging
		if loggPrepInfo.Enabled && isValidCredentialOption(loggPrepInfo.GCPCredentialInfo.CredentialsType) {
			gcpClient, err := getClient(context.Background(), loggPrepInfo.GCPCredentialInfo.ProjectID, tp.CredentialOption(loggPrepInfo.GCPCredentialInfo.CredentialsType), loggPrepInfo.GCPCredentialInfo.Credentials)
			if err != nil {
				log.Error().Str("module", "monitoring > Init").Msgf("get GCP client for Logging Error: %s", err)
				return
			}
			gcpSingletons = append(gcpSingletons, &cloudLogging{
				client:           gcpClient,
				reminderCron:     cron.New(cron.WithLocation(time.UTC)),
				reminderSchedule: loggPrepInfo.GCPCredentialInfo.CheckerSchedule,
			})
		}
	})
	return gcpSingletons
}

func getClient(ctx context.Context, projectID string, credType tp.CredentialOption, credentials string) (*logging.Client, error) {
	var client *logging.Client
	var err error

	switch credType {
	case tp.ApplicationDefaultCredentials:
		client, err = logging.NewClient(ctx, projectID)
	case tp.ServiceAccountCredentials:
		client, err = logging.NewClient(ctx, projectID, option.WithCredentialsFile(credentials))
	case tp.APIKey:
		client, err = logging.NewClient(ctx, projectID, option.WithAPIKey(credentials))
	case tp.OAuth2:
		tokenSource, err := google.DefaultTokenSource(ctx, logging.WriteScope)
		if err != nil {
			return nil, err
		}
		client, err = logging.NewClient(ctx, projectID, option.WithTokenSource(tokenSource))
		if err != nil {
			return nil, err
		}
	default:
		err = errors.New("invalid credential type")
	}
	return client, err
}

func isValidCredentialOption(option string) bool {
	credOption := tp.CredentialOption(option)
	_, isValid := validCredentialOptions[credOption]
	return isValid
}
