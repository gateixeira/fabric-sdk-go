/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mocks

import (
	"net/http"

	"time"

	cfapi "github.com/cloudflare/cfssl/api"
	cfsslapi "github.com/cloudflare/cfssl/api"
	"github.com/hyperledger/fabric-sdk-go/internal/github.com/hyperledger/fabric-ca/api"
	"github.com/hyperledger/fabric-sdk-go/internal/github.com/hyperledger/fabric-ca/util"
	"github.com/hyperledger/fabric-sdk-go/pkg/context/api/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/logging"
)

var logger = logging.NewLogger("fabsdk/core")

// Matching key-cert pair. On enroll, the key will be
// imported into the key store, and the cert will be
// returned to the caller.

const privateKey = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgp4qKKB0WCEfx7XiB
5Ul+GpjM1P5rqc6RhjD5OkTgl5OhRANCAATyFT0voXX7cA4PPtNstWleaTpwjvbS
J3+tMGTG67f+TdCfDxWYMpQYxLlE8VkbEzKWDwCYvDZRMKCQfv2ErNvb
-----END PRIVATE KEY-----`

const ecert = `-----BEGIN CERTIFICATE-----
MIICGTCCAcCgAwIBAgIRALR/1GXtEud5GQL2CZykkOkwCgYIKoZIzj0EAwIwczEL
MAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG
cmFuY2lzY28xGTAXBgNVBAoTEG9yZzEuZXhhbXBsZS5jb20xHDAaBgNVBAMTE2Nh
Lm9yZzEuZXhhbXBsZS5jb20wHhcNMTcwNzI4MTQyNzIwWhcNMjcwNzI2MTQyNzIw
WjBbMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMN
U2FuIEZyYW5jaXNjbzEfMB0GA1UEAwwWVXNlcjFAb3JnMS5leGFtcGxlLmNvbTBZ
MBMGByqGSM49AgEGCCqGSM49AwEHA0IABPIVPS+hdftwDg8+02y1aV5pOnCO9tIn
f60wZMbrt/5N0J8PFZgylBjEuUTxWRsTMpYPAJi8NlEwoJB+/YSs29ujTTBLMA4G
A1UdDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAAMCsGA1UdIwQkMCKAIIeR0TY+iVFf
mvoEKwaToscEu43ZXSj5fTVJornjxDUtMAoGCCqGSM49BAMCA0cAMEQCID+dZ7H5
AiaiI2BjxnL3/TetJ8iFJYZyWvK//an13WV/AiARBJd/pI5A7KZgQxJhXmmR8bie
XdsmTcdRvJ3TS/6HCA==
-----END CERTIFICATE-----`

// The enrollment response from the server
type enrollmentResponseNet struct {
	// Base64 encoded PEM-encoded ECert
	Cert string
	// The server information
	ServerInfo serverInfoResponseNet
}

// The response to the GET /info request
type serverInfoResponseNet struct {
	// CAName is a unique name associated with fabric-ca-server's CA
	CAName string
	// Base64 encoding of PEM-encoded certificate chain
	CAChain string
}

// MockFabricCAServer is a mock for FabricCAServer
type MockFabricCAServer struct {
	address     string
	cryptoSuite core.CryptoSuite
	running     bool
}

// Start fabric CA mock server
func (s *MockFabricCAServer) Start(address string, cryptoSuite core.CryptoSuite) error {

	if s.running {
		return nil
	}

	s.address = address
	s.cryptoSuite = cryptoSuite

	// Register request handlers
	http.HandleFunc("/register", s.register)
	http.HandleFunc("/enroll", s.enroll)
	http.HandleFunc("/reenroll", s.enroll)

	server := &http.Server{
		Addr:      s.address,
		TLSConfig: nil,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic("HTTP Server: Failed to start")
		}
	}()
	time.Sleep(1 * time.Second)
	logger.Infof("HTTP Server started on %s", s.address)

	s.running = true
	return nil

}

func (s *MockFabricCAServer) addKeyToKeyStore(privateKey []byte) error {
	// Import private key that matches the cert we will return
	// from this mock service, so it can be looked up by SKI from the cert
	_, err := util.ImportBCCSPKeyFromPEMBytes([]byte(privateKey), s.cryptoSuite, false)
	if err != nil {
		return err
	}
	return nil
}

// Register user
func (s *MockFabricCAServer) register(w http.ResponseWriter, req *http.Request) {
	resp := &api.RegistrationResponseNet{RegistrationResponse: api.RegistrationResponse{Secret: "mockSecretValue"}}
	cfsslapi.SendResponse(w, resp)
}

// Enroll user
func (s *MockFabricCAServer) enroll(w http.ResponseWriter, req *http.Request) {
	s.addKeyToKeyStore([]byte(privateKey))
	resp := &enrollmentResponseNet{Cert: util.B64Encode([]byte(ecert))}
	fillCAInfo(&resp.ServerInfo)
	cfapi.SendResponse(w, resp)
}

// Fill the CA info structure appropriately
func fillCAInfo(info *serverInfoResponseNet) {
	info.CAName = "MockCAName"
	info.CAChain = util.B64Encode([]byte("MockCAChain"))
}
