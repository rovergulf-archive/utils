package tlsconf

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
)

func CreateTLSClientConfig(certFile, keyFile, caFile string, verifySsl bool) (*tls.Config, error) {
	if certFile == "" && keyFile == "" && caFile == "" {
		return nil, nil
	}

	var err error
	t := &tls.Config{
		InsecureSkipVerify: verifySsl,
	}

	if certFile != "" && keyFile != "" {
		t, err = CreateTLSConfig(certFile, keyFile)
		if err != nil {
			return nil, err
		}
	}

	if caFile != "" {
		caCert, err := ioutil.ReadFile(caFile)
		if err != nil {
			return nil, err
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t.RootCAs = caCertPool
	}

	// will be nil by default if nothing is provided
	return t, nil
}

func CreateTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	if certFile == "" && keyFile == "" {
		return nil, nil
	}

	t := &tls.Config{}

	if certFile != "" && keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}
		t.Certificates = []tls.Certificate{cert}
	}

	// will be nil by default if nothing is provided
	return t, nil
}
