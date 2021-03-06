package main

import (
	"encoding/base64"
	"flag"
	"log"

	uuid "github.com/satori/go.uuid"
	"github.com/synapse-garden/sg-proto/auth"
	"github.com/synapse-garden/sg-proto/rest"

	"github.com/boltdb/bolt"
)

// CLI flags
var (
	Address  = flag.String("addr", "127.0.0.1", "the address to host on")
	Port     = flag.String("port", ":8080", "the port to listen on")
	DBAddr   = flag.String("db", "my.db", "the database to use")
	CertFile = flag.String("cert", "", "the certificate file to use")
	KeyFile  = flag.String("key", "cert.key", "the certificate key to use")
	ConfFile = flag.String("cfg", "conf.toml", "the config file to use")

	SourceLocation = flag.String(
		"source",
		"https://github.com/synapse-garden/sg-proto",
		"where the source is hosted",
	)

	DevMode  = flag.Bool("dev", false, "start in developer mode")
	RegenKey = flag.Bool("regen-key", false, "re-create the admin key")
)

// Source constants
const (
	SourceLicense = "Affero GPL V3"
	Licensee      = "SynapseGarden 2016"
)

func main() {
	flag.Parse()

	db, err := bolt.Open(*DBAddr, 0600, nil)
	if err != nil {
		log.Fatalf("unable to open Bolt database: %s", err.Error())
	}

	source := rest.SourceInfo{
		Location:   *SourceLocation,
		License:    SourceLicense,
		LicensedTo: Licensee,
	}

	var key auth.Token
	if *RegenKey {
		key = auth.Token(uuid.NewV4().Bytes())
		log.Printf("new admin key generated: %#q",
			base64.StdEncoding.EncodeToString(key))
	}
	switch {
	case *SourceLocation == "" && !*DevMode:
		log.Fatal("must provide a source location using -source")
	case *DevMode && *CertFile == "":
		devServeInsecure(db, key, *Address, *Port, source)
	case *DevMode:
		devServeSecure(db, key, *Address, *Port, *CertFile, *KeyFile, source)
	case *CertFile == "":
		serveInsecure(db, key, *Address, *Port, source)
	default:
		serveSecure(db, key, *Address, *Port, *CertFile, *KeyFile, source)
	}
}
