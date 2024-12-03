package sdk

import "errors"

var (
	ConfigNotExistsErr    = errors.New("config not exists")
	ConfigMissCertErr     = errors.New("config miss cert")
	ConfigMissCustomerErr = errors.New("config miss customer")
)
