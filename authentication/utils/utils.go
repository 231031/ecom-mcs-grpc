package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/231031/ecom-mcs-grpc/authentication/model"
	"github.com/golang-jwt/jwt"
)

func ParseArgon2Hash(encodedHash string) (*model.Argon2Configuration, error) {
	components := strings.Split(encodedHash, "$")
	if len(components) != 6 {
		return nil, errors.New("invalid hash format structure")
	}

	// Validate algorithm identifier
	if !strings.HasPrefix(components[1], "argon2id") {
		return nil, errors.New("unsupported algorithm variant")
	}

	// Extract version information
	var version int
	fmt.Sscanf(components[2], "v=%d", &version)

	// Parse configuration parameters
	config := &model.Argon2Configuration{}
	fmt.Sscanf(components[3], "m=%d,t=%d,p=%d",
		&config.MemoryCost, &config.TimeCost, &config.Threads)

	// Decode salt component
	salt, err := base64.RawStdEncoding.DecodeString(components[4])
	if err != nil {
		fmt.Sprintln("error decoding salt:", err)
		return nil, err
	}
	config.Salt = salt

	// Decode hash component
	hash, err := base64.RawStdEncoding.DecodeString(components[5])
	if err != nil {
		fmt.Sprintln("error decoding hash:", err)
		return nil, err
	}
	config.HashRaw = hash
	config.KeyLength = uint32(len(hash))

	return config, nil
}

// generateSalt creates a cryptographically secure random salt
func GenerateSalt(saltLength uint32) ([]byte, error) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func ConfigGenerateKey(cfg *model.Config) *model.TokenConfig {
	tokenCfg := &model.TokenConfig{
		TokenIDExpirationSecs: 10 * 60,
		RefreshExpirationSecs: 48 * 3600,
		RefreshSecret:         cfg.SecretKey,
	}

	priv, err := ioutil.ReadFile(cfg.FilePriPath)
	if err != nil {
		fmt.Sprintln("failed to read private pem file", err)
		return tokenCfg
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		fmt.Sprintln("failed to parse private pem to rsa", err)
		return tokenCfg
	}
	tokenCfg.PrivateKey = privateKey

	pub, err := ioutil.ReadFile(cfg.FilePubPath)
	if err != nil {
		fmt.Sprintln("failed to read public pem file", err)
		return tokenCfg
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		fmt.Sprintln("failed to parse public pem to rsa", err)
		return tokenCfg
	}
	tokenCfg.PublicKey = publicKey

	return tokenCfg
}
