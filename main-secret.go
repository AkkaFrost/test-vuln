package main

import (
	"fmt"
	"net/http"
	"net/smtp"
)

// ==================== Cloud Provider Credentials ====================

const (
	// AWS
	AWSAccessKeyID     = "AKIAIOSFODNN7EXAMPLE"
	AWSSecretAccessKey = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	AWSSessionToken    = "FwoGZXIvYXdzEBYaDHqa0AP1z2EXAMPLE//////////wEaDE5TSlVaQ0pzVk1REXAMPLETOKEN"

	// GCP
	GCPServiceAccountKey = `{
		"type": "service_account",
		"project_id": "my-production-project",
		"private_key_id": "key123abc456def",
		"private_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA2Z3VS5JJcds3xfn/ygWyF8PBWLEBNhTzPIGKJB0Wk\n-----END RSA PRIVATE KEY-----\n",
		"client_email": "deploy@my-production-project.iam.gserviceaccount.com",
		"client_id": "123456789012345678901",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token"
	}`

	// Azure
	AzureTenantID     = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	AzureClientID     = "12345678-abcd-efgh-ijkl-9876543210ab"
	AzureClientSecret = "xYz~8Q~abcDEFghiJKLmnoPQRstuvWXyz1234567"
	AzureStorageKey   = "DefaultEndpointsProtocol=https;AccountName=prodstorage;AccountKey=4dY8Gf+EXAMPLEKEY1234567890abcdefghijklmnopqrstuvwxyz==;EndpointSuffix=core.windows.net"
)

// ==================== Database Credentials ====================

var (
	PostgresConnString = "postgres://admin:P@ssw0rd!Pr0d@prod-db.us-east-1.rds.amazonaws.com:5432/production?sslmode=disable"
	MySQLConnString    = "root:MyS3cretDBPa$$@tcp(mysql-prod.internal:3306)/app_production"
	MongoURI           = "mongodb://appuser:M0ng0Pr0dP@ss!@mongo-cluster.internal:27017/production?authSource=admin&replicaSet=rs0"
	RedisURL           = "redis://:R3d1sS3cret!@redis-prod.internal:6379/0"
	ElasticSearchURL   = "https://elastic:El@st1cPr0d!@es-prod.internal:9200"
)

// ==================== API Keys & Tokens ====================

const (
	// OpenAI
	OpenAIAPIKey = "sk-proj-abc123def456ghi789jkl012mno345pqr678stu901vwx234"

	// Stripe
	StripeSecretKey      = "sk_live_51HG3CMExampleKeyDoNotUse1234567890abcdefghijklmnop"
	StripeWebhookSecret  = "whsec_1234567890abcdefghijklmnopqrstuvwxyz"

	// Twilio
	TwilioAccountSID = "AC1234567890abcdef1234567890abcdef"
	TwilioAuthToken  = "your_auth_token_1234567890abcdef"

	// SendGrid
	SendGridAPIKey = "SG.abcdefghijklmnop.qrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUV"

	// GitHub
	GitHubPAT          = "ghp_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdef123456"
	GitHubAppPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA0m59l2u9iDnMbrXHfqkOrn2daqYXisRTwKRAJBIYGA0Eexam
plekeydata1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRS
-----END RSA PRIVATE KEY-----`

	// GitLab
	GitLabToken = "glpat-xxxxxxxxxxxxxxxxxxxx"

	// Slack
	SlackBotToken    = "xoxb-123456789012-1234567890123-AbCdEfGhIjKlMnOpQrStUvWx"
	SlackSigningSecret = "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4"
	SlackWebhookURL  = "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"

	// Datadog
	DatadogAPIKey = "abcdef1234567890abcdef1234567890"
	DatadogAppKey = "abcdef1234567890abcdef1234567890abcdef12"

	// PagerDuty
	PagerDutyAPIKey = "u+ExAmPlEaPiKeY1234567890"

	// Sentry
	SentryDSN = "https://examplePublicKey@o0.ingest.sentry.io/0"
)

// ==================== OAuth / JWT / Auth ====================

const (
	JWTSigningKey       = "super-secret-jwt-signing-key-do-not-leak-2024"
	OAuth2ClientSecret  = "GOCSPX-abcdefghijklmnopqrstuvwxyz12"
	Auth0ClientSecret   = "abcDEFghiJKL123456789_mnoPQRstuVWX"
	SessionEncryptionKey = "32-byte-long-encryption-key!1234"
)

// ==================== SSH & Certificates ====================

var SSHPrivateKey = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACBExampleKeyDataHere1234567890abcdefghijklmnopqrstuvwxyz
-----END OPENSSH PRIVATE KEY-----`

var TLSPrivateKey = `-----BEGIN EC PRIVATE KEY-----
MHQCAQEEIODg7EXAMPLE1234567890abcdefghijklmnopqrstuvwxyzABCDEFGH
IJKLMNOPQRS=
-----END EC PRIVATE KEY-----`

var TLSCertificate = `-----BEGIN CERTIFICATE-----
MIIDXTCCAkWgAwIBAgIJAJC1EXAMPLE1234567890abcdefghijklmnopqrstuvwx
yzABCDEFGHIJKLMNOPQRSTUVWXYZ
-----END CERTIFICATE-----`

// ==================== Encryption Keys ====================

const (
	AES256Key        = "0123456789abcdef0123456789abcdef"
	HMACSecret       = "hmac-shared-secret-for-webhook-verification"
	EncryptionIV     = "1234567890abcdef"
	MasterKeyHex     = "4a6f686e446f654d61737465724b6579313233343536373839306162636465"
)

// ==================== Third-Party Service Credentials ====================

const (
	// Payment processors
	PayPalClientID     = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPp"
	PayPalClientSecret = "EeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTt"

	// Email / SMTP
	SMTPHost     = "smtp.company.com"
	SMTPUser     = "noreply@company.com"
	SMTPPassword = "Sm7pP@ssw0rd!2024"

	// Docker Registry
	DockerRegistryPassword = "dckr_pat_EXAMPLETOKEN1234567890"

	// NPM
	NPMToken = "npm_EXAMPLETOKEN1234567890abcdefghij"

	// PyPI
	PyPIToken = "pypi-AgEIcHlwaS5vcmcExampleToken1234567890"

	// Terraform Cloud
	TerraformToken = "atlasv1.ExAmPlEtOkEn1234567890abcdefghijklmnop.qrstuvwxyz"

	// Vault
	VaultToken = "hvs.EXAMPLETOKEN1234567890abcdefghijklmnop"
	VaultAddr  = "https://vault.internal.company.com:8200"
)

// ==================== Internal / Infra ====================

const (
	KubernetesToken  = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmRlZmF1bHQ6ZGVwbG95ZXIifQ.EXAMPLESIGNATURE"
	ConsulACLToken   = "01234567-89ab-cdef-0123-456789abcdef"
	GrafanaAPIKey    = "eyJrIjoiT0EXAMPLE1234567890abcdefghijklmnop"
	NewRelicLicenseKey = "eu01xx1234567890abcdef1234567890abcdefNRAL"
)

// ==================== Usage ====================

func sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", SMTPUser, SMTPPassword, SMTPHost)
	msg := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))
	return smtp.SendMail(SMTPHost+":587", auth, SMTPUser, []string{to}, msg)
}

func callOpenAI(prompt string) (*http.Response, error) {
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", nil)
	req.Header.Set("Authorization", "Bearer "+OpenAIAPIKey)
	return http.DefaultClient.Do(req)
}

func callStripe(amount int) (*http.Response, error) {
	req, _ := http.NewRequest("POST", "https://api.stripe.com/v1/charges", nil)
	req.Header.Set("Authorization", "Bearer "+StripeSecretKey)
	return http.DefaultClient.Do(req)
}

func main() {
	fmt.Println("This file is for demonstration purposes only.")
	fmt.Println("It contains example secrets that should be detected by scanners.")
	fmt.Printf("DB: %s\n", PostgresConnString)
	fmt.Printf("AWS Key: %s\n", AWSAccessKeyID)
	fmt.Printf("Stripe: %s\n", StripeSecretKey)
}
