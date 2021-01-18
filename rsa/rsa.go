package main

import (
	"bufio"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"github.com/hackaio/pk"
	"os"
)

type rsaEncoder struct {
	PublicKey *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

func NewEncoder(pubKey rsa.PublicKey,privKey rsa.PrivateKey) pk.Encoder {
	return &rsaEncoder{
		PublicKey:  &pubKey,
		PrivateKey: &privKey,
	}
}

func (r rsaEncoder) Encode(password string) ([]byte, error) {
	encryptedBytes, err := rsa.EncryptOAEP(
			sha256.New(),
			rand.Reader,
			r.PublicKey,
			[]byte(password),
			[]byte("passwords"),
		)

	return encryptedBytes, err
}

func (r rsaEncoder) Decode(encoded []byte) (string, error) {
	// The first argument is an optional random data generator (the rand.Reader we used before)
	// we can set this value as nil
	// The OEAPOptions in the end signify that we encrypted the data using OEAP, and that we used
	// SHA256 to hash the input.
	decryptedBytes, err := r.PrivateKey.
		Decrypt(nil, encoded, &rsa.OAEPOptions{Hash: crypto.SHA256,Label: []byte("passwords")})

	return string(decryptedBytes), err

}

type signer struct {
	PublicKey *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

func NewSigner(pubKey rsa.PublicKey,privKey rsa.PrivateKey) pk.Signer {
	return &signer{
		PublicKey:  &pubKey,
		PrivateKey: &privKey,
	}
}


func (s signer) Sign(password string) (digest []byte ,signature []byte, err error) {
	msg := []byte(password)

	// Before signing, we need to hash our message
	// The hash is what we actually sign
	msgHash := sha256.New()
	_, err = msgHash.Write(msg)
	if err != nil {
		return nil,nil, err
	}
	digest = msgHash.Sum(nil)

	// In order to generate the signature, we provide a random number generator,
	// our private key, the hashing algorithm that we used, and the hash sum
	// of our message
	signature, err = rsa.SignPSS(
		rand.Reader,
		s.PrivateKey,
		crypto.SHA256,
		digest,
		nil)

	if err != nil {
		return nil,nil, err
	}

	return digest,signature,nil

}

func (s signer) Verify(digest []byte, signature []byte) (err error) {
	// To verify the signature, we provide the public key, the hashing algorithm
	// the hash sum of our message and the signature we generated previously
	// there is an optional "options" parameter which can omit for now
	err = rsa.VerifyPSS(s.PublicKey, crypto.SHA256, digest, signature, nil)
	if err != nil {
		return err
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

	pemfile, err := os.Create(fileName)
	checkError(err)
	defer pemfile.Close()

	err = pem.Encode(pemfile, pemkey)
	checkError(err)
}

func loadPrivKeyFromPEM(filename string)(key *rsa.PrivateKey, err error){

	//All right! Now we have our RSA key pair created and exported
	//to a PEM file.
	//It’s time to see how we can import them from this same file.
	//Before starting, we need to open PEM file:
	privateKeyFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)
	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)
	data, _ := pem.Decode([]byte(pembytes))
	privateKeyFile.Close()

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKeyImported, err
}

func loadPubKeyFromPEM(filename string)(key *rsa.PublicKey, err error){
	keyFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	fileInfo, _ := keyFile.Stat()
	var size int64 = fileInfo.Size()
	fileBytes := make([]byte, size)
	buffer := bufio.NewReader(keyFile)
	_, err = buffer.Read(fileBytes)
	data, _ := pem.Decode([]byte(fileBytes))
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




func main() {;
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


	privkey, err := loadPrivKeyFromPEM("private.pem")
	pubkey, err := loadPubKeyFromPEM("public.pem")

	checkError(err)

	rsaEnc := rsaEncoder{
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

	passSigner := signer{
		PublicKey:  pubkey,
		PrivateKey: privkey,
	}

	msg := "verifiable message"

	digest,signature, err := passSigner.Sign(msg)

	err = passSigner.Verify(digest,signature)

	// If we don't get any error from the `VerifyPSS` method, that means our
	// signature is valid
	fmt.Println("signature verified")
}

