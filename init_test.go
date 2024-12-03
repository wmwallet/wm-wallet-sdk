package sdk

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	fn := func(path string) []byte {
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("fail: %n", err)
		}
		return b
	}

	tests := []struct {
		name    string
		ops     []Option
		wantErr bool
	}{
		{
			name: "byFilePath",
			ops: []Option{
				WithCertPath("config/test_ca.crt", "config/test_client.crt", "config/test_client.key"),
				WithCustomer("a"),
				WithSecretPath("test/public_key.pem"),
			},
			wantErr: false,
		},
		{
			name: "byByte",
			ops: []Option{
				WithCertBytes(fn("config/test_ca.crt"), fn("config/test_client.crt"), fn("config/test_client.key")),
				WithSecretBytes(fn("config/public_key.pem")),
				WithCustomer("a"),
			},
			wantErr: false,
		},
		{
			name: "testMode",
			ops: []Option{
				WithCustomer("a"),
				WithSecretPath("config/public_key.pem"),
				WithTest(true),
			},
			wantErr: false,
		},
		{
			name:    "nil cfg",
			ops:     nil,
			wantErr: true,
		},
		{
			name: "err cfg",
			ops: []Option{
				WithCertBytes(fn("config/test_ca.crt"), nil, nil),

				WithCustomer("a"),
			},
			wantErr: true,
		},
		{
			name: "ca path err",
			ops: []Option{
				WithCertPath(("test_ca.crt"), ("test_client.crt"), ("test_client.key")),
				WithSecretPath("test/public_key.pem"),

				WithCustomer("a"),
			},
			wantErr: true,
		},
		{
			name: "ca err",
			ops: []Option{
				WithCertBytes(fn("config/test_ca.crt")[:10], nil, nil),
				WithSecretPath("config/public_key.pem"),

				WithCustomer("a"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Init(tt.ops...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClientAndServer(t *testing.T) {
	cfg := &config{
		caCertPath: "config/test_ca.crt",
		certPath:   "config/test_client.crt",
		keyPath:    "config/test_client.key",
		customer:   "a",
	}
	go func() {
		caCert, err := os.ReadFile(cfg.caCertPath)
		if err != nil {
			t.Error(err)
			return
		}

		// 创建 CA 证书池
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		serverCert, err := tls.LoadX509KeyPair("config/test_server.crt", "config/test_server.key")
		if err != nil {
			t.Error(err)
			return
		}

		config := &tls.Config{
			Certificates: []tls.Certificate{serverCert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    caCertPool,
		}

		fn := func(w http.ResponseWriter, r *http.Request) {
			_, err = fmt.Fprintf(w, "ok")
			if err != nil {
				return
			}
		}
		// 创建 HTTPS 服务器
		server := &http.Server{
			Addr:      ":8443",
			TLSConfig: config,
			Handler:   http.HandlerFunc(fn),
		}
		err = server.ListenAndServeTLS("", "")
		if err != nil {
			return
		}
	}()

	ops := []Option{WithCertPath(("config/test_ca.crt"), ("config/test_client.crt"), ("config/test_client.key")),
		WithCustomer("a"),
		WithSecretPath("config/public_key.pem"),
	}
	cli, err := Init(ops...)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := cli.client.Get("https://localhost:8443/")
	if err != nil {
		t.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "ok" {
		t.Fatal(string(body))
	}
	t.Log(string(body))
}
