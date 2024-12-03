package sdk

type config struct {
	caCert []byte
	cert   []byte
	key    []byte

	caCertPath string
	certPath   string
	keyPath    string

	isTest bool

	customer string

	secretPath  string
	secretBytes []byte
}
