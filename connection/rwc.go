package connection

import (
	"crypto/ed25519"
	"net"

	tmcrypto "github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/protoio"
	p2pconn "github.com/tendermint/tendermint/p2p/conn"
)

// ReadWriteConn holds the secret connection, the reader and the writer.
type ReadWriteConn struct {
	// SecretConn holds the secret connection for communication between
	// Pairmint and Tendermint.
	SecretConn *p2pconn.SecretConnection

	// Reader is used to read from the TCP stream.
	Reader protoio.ReadCloser

	// Writer is used to write to the TCP stream.
	Writer protoio.WriteCloser
}

// NewReadWriteConn returns a new instance of ReadWriteConn.
func NewReadWriteConn() *ReadWriteConn {
	return &ReadWriteConn{
		SecretConn: new(p2pconn.SecretConnection),
	}
}

// RetrySecretDial dials the given address until success and returns
// a secret connection.
func RetrySecretDial(protocol, address string, privkey ed25519.PrivateKey) (*p2pconn.SecretConnection, error) {
	var conn net.Conn
	var err error

	for {
		if conn, err = net.Dial(protocol, address); err == nil {
			break
		}
	}

	secretConn, err := p2pconn.MakeSecretConnection(conn, tmcrypto.PrivKey(privkey))
	if err != nil {
		return nil, err
	}

	return secretConn, nil
}