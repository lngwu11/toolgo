package utils

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"time"
)

// NewClientTLSConfig
// rootCert 根证书路径
// certFile 客户端证书路径,客户端证书经过根证书签名
// keyFile  客户端私钥路径
func NewClientTLSConfig(rootCert, certFile, keyFile string) (config *tls.Config, err error) {
	var cert tls.Certificate
	var certBytes []byte

	cert, err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		err = errors.Wrapf(err, "无法加载证书 [certFile=%v] [keyFile=%v]", certFile, keyFile)
		return
	}

	config = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	if len(rootCert) > 0 {
		certBytes, err = ioutil.ReadFile(rootCert)
		if err != nil {
			err = errors.Wrapf(err, "无法读取root根证书 [rootCert=%v]", rootCert)
			return
		}
		rootCertPool := x509.NewCertPool()
		ok := rootCertPool.AppendCertsFromPEM(certBytes)
		if !ok {
			err = errors.Wrapf(err, "无法加载root根证书 [rootCert=%v]", rootCert)
			return
		}
		config.InsecureSkipVerify = false
		config.RootCAs = rootCertPool
	} else {
		config.InsecureSkipVerify = true
		config.RootCAs = nil
	}

	return config, nil
}

// NewServerTLSConfig
// rootCert 根证书路径
// certFile 客户端证书路径,客户端证书经过根证书签名
// keyFile  客户端私钥路径
func NewServerTLSConfig(rootCert, certFile, keyFile string) (config *tls.Config, err error) {
	var cert tls.Certificate
	var certBytes []byte

	cert, err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		err = errors.Wrapf(err, "无法加载证书 [certFile=%v] [keyFile=%v]", certFile, keyFile)
		return
	}

	setTCPKeepAlive := func(clientHello *tls.ClientHelloInfo) (cfg *tls.Config, err error) {
		// Check that the underlying connection really is TCP.
		if tcpConn, ok := clientHello.Conn.(*net.TCPConn); ok {
			if err = tcpConn.SetKeepAlive(true); err != nil {
				err = errors.Wrapf(err, "SetKeepAlive for tls")
			}
			if err = tcpConn.SetKeepAlivePeriod(30 * time.Second); err != nil {
				err = errors.Wrapf(err, "SetKeepAlivePeriod for tls")
			}
		} else {
			err = errors.Errorf("TLS over non-TCP connection")
		}

		// ignore err
		// Make sure to return nil, nil to let the caller fall back on the default behavior.
		return nil, nil
	}

	config = &tls.Config{
		Certificates:       []tls.Certificate{cert},
		GetConfigForClient: setTCPKeepAlive,
		/*CipherSuites: []uint16{
			//tls.TLS_RSA_WITH_RC4_128_SHA,
			//tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			//tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
			//tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,

			// TLS 1.3 cipher suites.
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},*/
		//MinVersion: tls.VersionTLS13,
	}

	if len(rootCert) > 0 {
		certBytes, err = ioutil.ReadFile(rootCert)
		if err != nil {
			err = errors.Wrapf(err, "无法读取root根证书 [rootCert=%v]", rootCert)
			return
		}
		rootCertPool := x509.NewCertPool()
		ok := rootCertPool.AppendCertsFromPEM(certBytes)
		if !ok {
			err = errors.Wrapf(err, "无法加载root根证书 [rootCert=%v]", rootCert)
			return
		}
		config.ClientAuth = tls.RequireAndVerifyClientCert
		config.ClientCAs = rootCertPool
	} else {
		config.ClientAuth = tls.NoClientCert
		config.ClientCAs = nil
	}

	return config, nil
}
