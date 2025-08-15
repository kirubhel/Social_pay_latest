package procedure

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

func Encrypt(publicKeyPEM string, data []byte) (string, error) {
	// Decode base64 PEM
	pemBytes, err := base64.StdEncoding.DecodeString(publicKeyPEM)
	if err != nil {
		return "", err
	}

	// Parse PEM block
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return "", errors.New("failed to parse PEM block")
	}

	// Parse public key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// Encrypt data
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, pub.(*rsa.PublicKey), data)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// Decrypt decrypts data using RSA private key
func Decrypt(privateKeyPEM string, encryptedData string) ([]byte, error) {
	// Decode base64 PEM
	pemBytes, err := base64.StdEncoding.DecodeString(privateKeyPEM)
	if err != nil {
		return nil, err
	}

	// Parse PEM block
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}

	// Parse private key
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// Decode base64 encrypted data
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	// Decrypt data
	return rsa.DecryptPKCS1v15(rand.Reader, privKey, ciphertext)
}
