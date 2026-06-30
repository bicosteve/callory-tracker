package db

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

// RegisterTLSConfig registers a custom TLS configuration with the MySQL driver
// using the CA certificate provided by Aiven (the ca.pem file you download from
// the Aiven service overview page).
//
// The returned name should be used as the value of the "tls" parameter in the
// DSN, e.g. "...?tls=aiven".
func RegisterTLSConfig(name, caCertPath string) error {
	if caCertPath == "" {
		return errors.New("ca certificate path is empty")
	}

	pem, err := os.ReadFile(caCertPath)
	if err != nil {
		return fmt.Errorf("could not read CA cert %q: %w", caCertPath, err)
	}

	rootCertPool := x509.NewCertPool()
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return errors.New("failed to append CA certificate to pool")
	}

	err = mysql.RegisterTLSConfig(name, &tls.Config{
		RootCAs:    rootCertPool,
		MinVersion: tls.VersionTLS12,
	})
	if err != nil {
		return fmt.Errorf("could not register TLS config: %w", err)
	}

	return nil
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn) // -> initialize connection pools for future use
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to db")

	// db.SetConnMaxLifetime(time.Minute * 3)
	// db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns(10)

	return db, nil
}
