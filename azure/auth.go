package azure

type TokenCredentials struct {
	SubscriptionId string
	Token          string
}

type CertificateCredentials struct {
	SubscriptionId string
	Certificate    []byte
}

func NewTokenCredentials(subscriptionId string, token string) (TokenCredentials, error) {
	return TokenCredentials{subscriptionId, token}, nil
}

func NewCertificateCredentials(subscriptionId string, certificate []byte) (CertificateCredentials, error) {
	return CertificateCredentials{subscriptionId, certificate}, nil
}
