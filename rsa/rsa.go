/*
 * Copyright © 2021 PIUS ALFRED me.pius1102@gmail.com
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"bufio"
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"github.com/hackaio/pk"
	"github.com/hackaio/pk/pkg/errors"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
)

type rsaEncoderSigner struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

var _ pk.Encoder = (*rsaEncoderSigner)(nil)
var _ pk.Signer = (*rsaEncoderSigner)(nil)

func NewEncoderSigner(homeDir string) (es pk.EncoderSigner, err error) {

	credsDir := filepath.Join(homeDir, pk.AppDir, pk.CredDir)
	privateKeyPEMFile := filepath.Join(credsDir, "private.pem")
	publicKeyPEMFile := filepath.Join(credsDir, "public.pem")

	fmt.Printf("scanning creds at %v and %v ..... \n", privateKeyPEMFile, publicKeyPEMFile)

	_, errNoPrivatePEMFile := os.Stat(privateKeyPEMFile)
	_, errNoPubKeyPEMFile := os.Stat(publicKeyPEMFile)

	if errNoPrivatePEMFile == nil && errNoPubKeyPEMFile == nil {
		fmt.Printf("creds files exists... load em up\n")
		pubKey, privateKey, err := loadCredentials(credsDir)
		es = rsaEncoderSigner{
			PublicKey:  pubKey,
			PrivateKey: privateKey,
		}

		if err != nil {
			return nil, err
		}

	} else if os.IsNotExist(errNoPubKeyPEMFile) || os.IsNotExist(errNoPrivatePEMFile) {
		fmt.Printf("creds files do not exists... time to craft\n")
		err = initCredentials(homeDir)
	}

	pubKey, privateKey, err := loadCredentials(credsDir)

	if err != nil {
		return nil, err
	}
	return &rsaEncoderSigner{
		PublicKey:  pubKey,
		PrivateKey: privateKey,
	}, nil
}

func (r rsaEncoderSigner) Encode(password string) ([]byte, error) {
	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		r.PublicKey,
		[]byte(password),
		[]byte("passwords"),
	)

	return encryptedBytes, err
}

func (r rsaEncoderSigner) Decode(encoded []byte) (string, error) {
	// The first argument is an optional random data generator (the rand.Reader we used before)
	// we can set this value as nil
	// The OEAPOptions in the end signify that we encrypted the data using OEAP, and that we used
	// SHA256 to hash the input.
	decryptedBytes, err := r.PrivateKey.
		Decrypt(nil, encoded, &rsa.OAEPOptions{Hash: crypto.SHA256, Label: []byte("passwords")})

	return string(decryptedBytes), err

}

func (r rsaEncoderSigner) Sign(password string) (digest []byte, signature []byte, err error) {
	msg := []byte(password)

	// Before signing, we need to hash our message
	// The hash is what we actually sign
	msgHash := sha256.New()
	_, err = msgHash.Write(msg)
	if err != nil {
		return nil, nil, errors.Wrap(err, pk.ErrCriticalFailure)
	}
	digest = msgHash.Sum(nil)

	// In order to generate the signature, we provide a random number generator,
	// our private key, the hashing algorithm that we used, and the hash sum
	// of our message
	signature, err = rsa.SignPSS(
		rand.Reader,
		r.PrivateKey,
		crypto.SHA256,
		digest,
		nil)

	if err != nil {
		return nil, nil, errors.Wrap(err, pk.ErrCriticalFailure)
	}

	return digest, signature, nil
}

func (r rsaEncoderSigner) Verify(password string, dbDigest []byte, dbSignature []byte) (err error) {

	digest, _, err := r.Sign(password)

	//Compares the digest of recovered password and that retrieved from db
	comp := bytes.Compare(digest, dbDigest)

	if comp != 0 {
		return pk.ErrInternalError
	}

	//If the digest compares we verify it with the stored signature
	// To verify the signature, we provide the public key, the hashing algorithm
	// the hash sum of our message and the signature we generated previously
	// there is an optional "options" parameter which can omit for now
	err = rsa.VerifyPSS(r.PublicKey, crypto.SHA256, digest, dbSignature, nil)
	if err != nil {
		return pk.ErrInternalError
	}

	return nil

}

func savePEMKey(fileName string, key *rsa.PrivateKey) error {
	outFile, err := os.Create(fileName)
	if err != nil {
		errMsg := errors.New(fmt.Sprintf("could not create %v file\n", fileName))
		return errors.Wrap(errMsg, err)
	}

	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	if err != nil {
		errMsg := errors.New("could not encode file to pem format")
		return errors.Wrap(errMsg, err)
	}

	return nil
}

func savePublicPEMKey(fileName string, pubkey rsa.PublicKey) (err error) {
	asn1Bytes, err := asn1.Marshal(pubkey)
	if err != nil {
		errMsg := errors.New("could not marshal the public key bytes")
		return errors.Wrap(errMsg, err)
	}

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemFile, err := os.Create(fileName)
	if err != nil {
		errMsg := errors.New(fmt.Sprintf("could not create %v file\n", fileName))
		return errors.Wrap(errMsg, err)
	}
	defer pemFile.Close()

	err = pem.Encode(pemFile, pemkey)
	if err != nil {
		errMsg := errors.New(fmt.Sprintf("could not encode %v file to pem format\n", fileName))
		return errors.Wrap(errMsg, err)
	}

	return nil
}

func privateKeyFromPEM(filename string) (key *rsa.PrivateKey, err error) {

	//All right! Now we have our RSA key pair created and exported
	//to a PEM file.
	//It’s time to see how we can import them from this same file.
	//Before starting, we need to open PEM file:
	privateKeyFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	pemFileInfo, _ := privateKeyFile.Stat()
	var size = pemFileInfo.Size()
	pemBytes := make([]byte, size)
	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pemBytes)
	data, _ := pem.Decode(pemBytes)
	_ = privateKeyFile.Close()

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKeyImported, err
}

func pubKeyFromPEM(filename string) (key *rsa.PublicKey, err error) {
	keyFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	fileInfo, _ := keyFile.Stat()
	var size = fileInfo.Size()
	fileBytes := make([]byte, size)
	buffer := bufio.NewReader(keyFile)
	_, err = buffer.Read(fileBytes)
	data, _ := pem.Decode(fileBytes)
	_ = keyFile.Close()

	publicKeyFromFile, err := x509.ParsePKCS1PublicKey(data.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKeyFromFile, err
}

func loadCredentials(credentialsDir string) (
	publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey, err error) {
	privateKeyPEMFile := filepath.Join(credentialsDir, "private.pem")
	privateKey, err = privateKeyFromPEM(privateKeyPEMFile)
	if err != nil {
		return nil, nil, err
	}

	publicKeyPEMFile := filepath.Join(credentialsDir, "public.pem")
	publicKey, err = pubKeyFromPEM(publicKeyPEMFile)

	if err != nil {
		return nil, nil, err
	}

	return publicKey, privateKey, nil

}

func initCredentials(homeDir string) (err error) {

	fmt.Printf("init credentials .....\n")
	//Check if File exists if not create the credentials and save them
	//then load the credentials
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// The public key is a part of the *rsa.PrivateKey struct
	publicKey := privateKey.PublicKey

	privateKeyPEMFile := filepath.Join(homeDir, pk.AppDir, pk.CredDir, "private.pem")
	publicKeyPEMFile := filepath.Join(homeDir, pk.AppDir, pk.CredDir, "public.pem")

	fmt.Printf("saving creds at: %v and %v\n", privateKeyPEMFile, publicKeyPEMFile)

	err = savePEMKey(privateKeyPEMFile, privateKey)
	err = savePublicPEMKey(publicKeyPEMFile, publicKey)

	return err
}

func main() {

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Printf("can not find home: %v\n", err)
		os.Exit(1)
	}

	appHomePath := filepath.Join(home, pk.AppDir)
	appCredsPath := filepath.Join(appHomePath, pk.CredDir)
	appDBPath := filepath.Join(appHomePath, pk.DBDir)
	err = os.MkdirAll(appCredsPath, 0777)
	err = os.MkdirAll(appDBPath, 0777)
	if err != nil {
		fmt.Printf("can not create dir: %v\n", err)
		os.Exit(1)
	}

	es, err := NewEncoderSigner(home)

	if err != nil {
		fmt.Printf("could not create encoder signer due to %v\n", err)
		panic(err)
	}
	/*


		// The GenerateKey method takes in a reader that returns random bits, and
		// the number of bits
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			panic(err)
		}

		// The public key is a part of the *rsa.PrivateKey struct
		publicKey := privateKey.PublicKey

		savePEMKey("private.pem",privateKey)
		savePublicPEMKey("public.pem",publicKey)


		privkey, err := privateKeyFromPEM("private.pem")
		pubkey, err := pubKeyFromPEM("public.pem")

		checkError(err)

		es := rsaEncoderSigner{
			PublicKey:  pubkey,
			PrivateKey: privkey,
		}
	*/
	bt, err := es.Encode("mypasswordphrase")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", bt)

	btStr, err := es.Decode(bt)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", btStr)

	msg := "message"

	digest, signature, err := es.Sign(msg)

	err = es.Verify("message", digest, signature)

	if err != nil {
		panic(err)
		os.Exit(1)
	}
	// If we don't get any error from the `VerifyPSS` method, that means our
	// signature is valid
	fmt.Println("signature verified")
}
