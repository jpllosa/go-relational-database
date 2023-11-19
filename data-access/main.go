package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-sql-driver/mysql"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

var db *sql.DB

func main() {
	cfg := mysql.NewConfig()
	cfg.User = "golang"
	cfg.Passwd = "golang"
	cfg.Net = "tcp"
	cfg.Addr = "localhost:3306"
	cfg.DBName = "recordings"

	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile("../pems/my-ca.pem")
	if err != nil {
		log.Fatalf("configuration: reading CA pem file: %v", err)
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		log.Fatalf("configuration: failed to append pem file: %v", err)
	}
	clientCert := make([]tls.Certificate, 0, 1)
	certs, err := tls.LoadX509KeyPair("../pems/my-client-cert.pem", "../pems/my-client-key.pem")
	if err != nil {
		log.Fatalf("configuration: failed to load key pair: %v", err)
	}
	clientCert = append(clientCert, certs)
	mysql.RegisterTLSConfig("secure", &tls.Config{
		RootCAs:      rootCertPool,
		Certificates: clientCert,
		MinVersion:   tls.VersionTLS10, //without this defaults tls1.3 which not supported by our mysql
		MaxVersion:   tls.VersionTLS11,
	})
	cfg.TLSConfig = "secure"

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Securely connected!")

	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		fmt.Errorf("database query: %v", err)
	}
	defer rows.Close()

	fmt.Printf("%2s %15s %15s %6s \n", "ID", "Title", "Artist", "Price")
	for rows.Next() {
		var alb Album
		err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)
		if err != nil {
			fmt.Errorf("row scan: %v", err)
		} else {
			fmt.Printf("%2d %15s %15s %6.2f \n", alb.ID, alb.Title, alb.Artist, alb.Price)
		}
	}
}
