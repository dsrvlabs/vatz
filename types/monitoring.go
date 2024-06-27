package types

/*
 * gcp Types
 *
 */

type CredentialOption string

const (
	MonitoringIdentifier                           = "vatz"
	ApplicationDefaultCredentials CredentialOption = "ADC"
	ServiceAccountCredentials     CredentialOption = "SAC"
	APIKey                        CredentialOption = "APIKey"
	OAuth2                        CredentialOption = "OAuth"
)
