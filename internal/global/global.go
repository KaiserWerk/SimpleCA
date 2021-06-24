package global

import (
	"github.com/KaiserWerk/CertMaker/internal/entity"
	"github.com/KaiserWerk/sessionstore"
	"net/http"
	"time"
)

const (
	CertificateMinDays     = 1
	CertificateMaxDays     = 182
	CertificateDefaultDays = 7
	CsrUploadMaxBytes      = 5 << 10
	ApiTokenLength         = 40
	ChallengeTokenLength   = 80

	TokenHeader               = "X-Api-Token"
	CertificateLocationHeader = "X-Certificate-Location"
	PrivateKeyLocationHandler = "X-Privatekey-Location"
	ChallengeLocationHeader   = "X-Challenge-Location"

	ApiPrefix = "/api/v1"

	OcspPath                      = ApiPrefix + "/ocsp"
	RootCertificateObtainPath     = ApiPrefix + "/root-certificate/obtain"
	CertificateRequestPath        = ApiPrefix + "/certificate/request"
	CertificateRequestWithCSRPath = ApiPrefix + "/certificate/request-with-csr"
	CertificateObtainPath         = ApiPrefix + "/certificate/%d/obtain"
	PrivateKeyObtainPath          = ApiPrefix + "/privatekey/%d/obtain"
	SolveChallengePath            = ApiPrefix + "/challenge/%d/solve"
	CertificateRevokePath         = ApiPrefix + "/certificate/%d/revoke"

	WellKnownPath = "/.well-known/certmaker-challenge/token.txt"

	RootCertificateFilename = "root-cert.pem"
	RootPrivateKeyFilename  = "root-key.pem"

	PemContentType = "application/x-pem-file"
)

var (
	DnsNamesToSkip = []string{"localhost", "127.0.0.1", "::1", "[::1]"}
)

var (
	config     *entity.Configuration
	sessMgr    *sessionstore.SessionManager
	httpClient http.Client
)

func init() {
	sessMgr = sessionstore.NewManager("CERTMAKERSESS")
	httpClient = http.Client{
		Timeout: 1500 * time.Millisecond,
	}
}

// GetClient returns a readily usable *http.Client
func GetClient() *http.Client {
	return &httpClient
}

// SetConfiguration sets the global configuration to a given object
func SetConfiguration(c *entity.Configuration) {
	config = c
}

// GetConfiguration fetches the global configuration
func GetConfiguration() *entity.Configuration {
	return config
}

// GetSessMgr fetches the global SessionManager
func GetSessMgr() *sessionstore.SessionManager {
	return sessMgr
}
