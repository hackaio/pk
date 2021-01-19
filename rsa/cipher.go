package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"github.com/hackaio/pk"
	"github.com/hackaio/pk/pkg/errors"
	"os"
	"path/filepath"
)

type cipher struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func loadCredentials(credentialsDir string) (
	publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey, err error) {
	privateKeyPEMFile := filepath.Join(credentialsDir, "private.pem")
	privateKey, err = loadPrivKeyFromPEM(privateKeyPEMFile)
	if err != nil {
		return nil, nil, err
	}

	publicKeyPEMFile := filepath.Join(credentialsDir, "public.pem")
	publicKey, err = loadPubKeyFromPEM(publicKeyPEMFile)

	if err != nil {
		return nil, nil, err
	}

	return publicKey, privateKey, nil

}


func NewCipher(credentialsDir string) (c pk.Cipher, err error) {

	privateKeyPEMFile := filepath.Join(credentialsDir, "private.pem")
	publicKeyPEMFile := filepath.Join(credentialsDir, "public.pem")

	_, err = os.Stat(privateKeyPEMFile)

	if os.IsExist(err){
		//create publicKey and privateKey
		err1 :=  initCipher(credentialsDir)
		if err1 != nil {
			return nil, err
		}
	}
	_,err = os.Stat(publicKeyPEMFile)

	if err != nil {
		return nil, err
	}


	/*if _, err := os.Stat(publicKeyPEMFile); err == nil {
		// path/to/whatever exists

	} */

	publicKey, privateKey, err := loadCredentials(credentialsDir)

	if err != nil {
		return nil, err
	}

	return &cipher{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}

func initCipher(credentialsDir string) (err error) {
	//Check if File exists if not create the credentials and save them
	//then load the credentials
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// The public key is a part of the *rsa.PrivateKey struct
	publicKey := privateKey.PublicKey

	privateKeyPEMFile := filepath.Join(credentialsDir, "private.pem")
	publicKeyPEMFile := filepath.Join(credentialsDir, "public.pem")

	savePEMKey(privateKeyPEMFile, privateKey)
	savePublicPEMKey(publicKeyPEMFile, publicKey)

	return nil
}

func (c cipher) Encode(password string) ([]byte, error) {
	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		c.PublicKey,
		[]byte(password),
		[]byte("passwords"),
	)

	return encryptedBytes, err
}

func (c cipher) Decode(encoded []byte) (string, error) {
	// The first argument is an optional random data generator (the rand.Reader we used before)
	// we can set this value as nil
	// The OEAPOptions in the end signify that we encrypted the data using OEAP, and that we used
	// SHA256 to hash the input.
	decryptedBytes, err := c.PrivateKey.
		Decrypt(nil, encoded, &rsa.OAEPOptions{Hash: crypto.SHA256, Label: []byte("passwords")})

	return string(decryptedBytes), err

}

func (c cipher) Sign(password string) (digest []byte, signature []byte, err error) {
	msg := []byte(password)

	// Before signing, we need to hash our message
	// The hash is what we actually sign
	msgHash := sha256.New()
	_, err = msgHash.Write(msg)
	if err != nil {
		return nil, nil, errors.Wrap(err, ErrCriticalFailure)
	}
	digest = msgHash.Sum(nil)

	// In order to generate the signature, we provide a random number generator,
	// our private key, the hashing algorithm that we used, and the hash sum
	// of our message
	signature, err = rsa.SignPSS(
		rand.Reader,
		c.PrivateKey,
		crypto.SHA256,
		digest,
		nil)

	if err != nil {
		return nil, nil, errors.Wrap(err, ErrCriticalFailure)
	}

	return digest, signature, nil

}

func (c cipher) Verify(password string, dbDigest []byte, dbSignature []byte) error {
	digest, _, err := c.Sign(password)

	//Compares the digest of recovered password and that retrieved from db
	comp := bytes.Compare(digest, dbDigest)

	if comp != 0 {
		return ErrInternalError
	}

	//If the digest compares we verify it with the stored signature
	// To verify the signature, we provide the public key, the hashing algorithm
	// the hash sum of our message and the signature we generated previously
	// there is an optional "options" parameter which can omit for now
	err = rsa.VerifyPSS(c.PublicKey, crypto.SHA256, digest, dbSignature, nil)
	if err != nil {
		return ErrInternalError
	}

	return nil
}
