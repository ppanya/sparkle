package sparkle

import (
	"fmt"
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/sparkle/external/line"
	sparklecrypto "github.com/octofoxio/sparkle/pkg/crypto"
	"net/url"
)

var (
	LocalMongoDBURL   = "mongodb://localhost:27017,localhost:27018,localhost:27019/?replicaSet=rs"
	LocalDatastoreURL = ""

	LocalHostURL           = "http://localhost:3009"
	LocalSparkleServiceURL = "localhost:7911"
	LocalSpikeServiceURL   = "localhost:8901"
)

// Collection name
const (
	UserCollectionName     = "user"
	IdentityCollectionName = "identity"
	SessionCollectionName  = "session"
)

type Config struct {
	Database    Database
	EmailSender EmailSender
	Fs          foundation.FileSystem

	DefaultEmailConfirmationTemplate string
	Host                             *url.URL // hostname for views
	Address                          *url.URL // GRPC address for binding service
	DefaultEmailSenderAddress        string

	TokenSigner sparklecrypto.TokenSigner

	UserCollectionName     string
	IdentityCollectionName string
	SessionCollectionName  string
	LineClient             line.LineClient
}

func NewConfig(system foundation.FileSystem) *Config {
	b, err := system.GetObject("./default-email-template.hbs")
	if err != nil {
		fmt.Println("new config error!")
		panic(err)
	}
	return &Config{
		DefaultEmailConfirmationTemplate: string(b),
		EmailSender:                      &ConsoleEmailSender{},
		Fs:                               system,
		DefaultEmailSenderAddress:        "twilight-sparkle9822@canterlot.edu",
		IdentityCollectionName:           IdentityCollectionName,
		UserCollectionName:               UserCollectionName,
		SessionCollectionName:            SessionCollectionName,
		LineClient:                       line.NewDefaultLineClient(),
	}
}

func (c *Config) UseJWTSignerWithHMAC256(secret string) *Config {
	c.TokenSigner = sparklecrypto.NewJWT(secret)
	return c
}
func (c *Config) SetTokenSigner(signer sparklecrypto.TokenSigner) *Config {
	c.TokenSigner = signer
	return c
}

func (c *Config) GetAddress() url.URL {
	if c.Address == nil {
		panic("config.Address is nil")
	}
	return *c.Address
}

// SetAddress set grpc service address
func (c *Config) SetAddress(value string) *Config {
	host, err := url.Parse(value)
	if err != nil {
		panic(err)
	}
	c.Address = host
	return c
}

func (c *Config) GetHost() url.URL {
	return *c.Host
}

// SeteHost // set view host (confirm email page and etc.)
func (c *Config) SetHost(value string) *Config {
	host, err := url.Parse(value)
	if err != nil {
		panic(err)
	}
	c.Host = host
	return c
}

func (c *Config) UseLocalFileSystem(path string) *Config {
	c.Fs = foundation.NewFileSystem(path, foundation.StaticMode_LOCAL)
	return c
}
func (c *Config) UseStatikFileSystem() *Config {
	c.Fs = foundation.NewFileSystem("", foundation.StaticMode_Statik)
	return c
}

func (c *Config) SetDatabase(value Database) *Config {
	c.Database = value
	return c
}

func (c *Config) SetDefaultEmailTemplate(value string) *Config {
	c.DefaultEmailConfirmationTemplate = value
	return c
}

func (c *Config) SetDefaultEmailSenderAddress(value string) *Config {
	c.DefaultEmailSenderAddress = value
	return c
}
