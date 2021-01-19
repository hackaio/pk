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

var (
	ErrInternalError   = errors.New("internal error, possible db compromise")
	ErrCriticalFailure = errors.New("could not perform critical operation")
)

type rsaEncoderSigner struct {
	PublicKey *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}


var _ pk.Encoder = (*rsaEncoderSigner)(nil)
var _ pk.Signer = (*rsaEncoderSigner)(nil)

func NewEncoderSigner(credentialsDir string) (pk.EncoderSigner,error){

	privateKeyPEMFile := filepath.Join(credentialsDir, "private.pem")
	publicKeyPEMFile := filepath.Join(credentialsDir, "public.pem")

	_, err := os.Stat(privateKeyPEMFile)

	if os.IsExist(err){
		//create publicKey and privateKey
		err1 :=  initCredentials(credentialsDir)
		if err1 != nil {
			return nil, err
		}
	}
	_,err = os.Stat(publicKeyPEMFile)

	if err != nil {
		return nil, err
	}


	pubKey, privateKey, err := loadCredentials(credentialsDir)

	if err != nil {
		return nil, err
	}
	return &rsaEncoderSigner{
		PublicKey:  pubKey,
		PrivateKey: privateKey,
	},nil
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
		Decrypt(nil, encoded, &rsa.OAEPOptions{Hash: crypto.SHA256,Label: []byte("passwords")})

	return string(decryptedBytes), err

}

func (r rsaEncoderSigner) Sign(password string) (digest []byte ,signature []byte, err error) {
	msg := []byte(password)

	// Before signing, we need to hash our message
	// The hash is what we actually sign
	msgHash := sha256.New()
	_, err = msgHash.Write(msg)
	if err != nil {
		return nil,nil, errors.Wrap(err, ErrCriticalFailure)
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
		return nil,nil, errors.Wrap(err,ErrCriticalFailure)
	}

	return digest,signature,nil
}

func (r rsaEncoderSigner) Verify(password string, dbDigest []byte, dbSignature []byte) (err error) {

	digest, _, err := r.Sign(password)


	//Compares the digest of recovered password and that retrieved from db
	comp := bytes.Compare(digest,dbDigest)

	if comp != 0{
		return ErrInternalError
	}

	//If the digest compares we verify it with the stored signature
	// To verify the signature, we provide the public key, the hashing algorithm
	// the hash sum of our message and the signature we generated previously
	// there is an optional "options" parameter which can omit for now
	err = rsa.VerifyPSS(r.PublicKey, crypto.SHA256, digest, dbSignature, nil)
	if err != nil {
		return ErrInternalError
	}

	return nil

}


func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	checkError(err)
}

func savePublicPEMKey(fileName string, pubkey rsa.PublicKey) {
	asn1Bytes, err := asn1.Marshal(pubkey)
	checkError(err)

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemFile, err := os.Create(fileName)
	checkError(err)
	defer pemFile.Close()

	err = pem.Encode(pemFile, pemkey)
	checkError(err)
}

func privateKeyFromPEM(filename string)(key *rsa.PrivateKey, err error){

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

func pubKeyFromPEM(filename string)(key *rsa.PublicKey, err error){
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

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
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

func initCredentials(credentialsDir string) (err error) {
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





func main() {

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Printf("can not find home: %v\n",err)
		os.Exit(1)
	}
	appHome := filepath.Join(home, ".pk", "creds")
	err = os.MkdirAll(appHome, 0700)
	if err != nil {
		fmt.Printf("can not create dir: %v\n",err)
		os.Exit(1)
	}


	fmt.Printf("%v\n",appHome)

	_, err = NewCipher(appHome)



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

	rsaEnc := rsaEncoderSigner{
		PublicKey:  pubkey,
		PrivateKey: privkey,
	}

	bt,err := rsaEnc.Encode("mypasswordphrase")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n",bt)

	btStr, err := rsaEnc.Decode(bt)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n",btStr)

	msg := "message"

	digest,signature, err := rsaEnc.Sign(msg)

	err = rsaEnc.Verify("message",digest,signature)

	if err != nil {
		panic(err)
		os.Exit(1)
	}
	// If we don't get any error from the `VerifyPSS` method, that means our
	// signature is valid
	fmt.Println("signature verified")
}



