package data

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	//"sync"
)

type Signature struct {
	Signature []byte
	hashed []byte
	h crypto.Hash
}

/* Event()
*
* /getQueryEvent API
* /postQueryEvent API
* To search the event info
 */
func GenerateKeyPair(bitSize int) (*rsa.PrivateKey) {
	reader := rand.Reader
	key, err := rsa.GenerateKey(reader, bitSize)
	checkError(err)
	return key
}

/* Event()
*
* /getQueryEvent API
* /postQueryEvent API
* To search the event info
 */
func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

/*PrivateKeyToBytes()
*
* To PublicKeyToBytes public key to bytes
*
 */
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

/* PublicKeyToBytes()
*
* To PublicKeyToBytes public key to bytes
*
 */
func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		fmt.Println(err)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}


/* BytesToPublicKey()
*
* To concet data structure
*
 */
func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			fmt.Println(err)
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		fmt.Println(err)
	}
	return key
}


/* BytesToPublicKey()
*
* To encrypt transaction
*
 */
func BytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			fmt.Println(err)
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		fmt.Println(err)
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		fmt.Println("Not Ok")
	}
	return key
}

/* EncryptWithPublicKey()
*
* To encrypt transaction
*
*/
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		fmt.Println(err)
	}
	return ciphertext
}

/* DecryptWithPrivateKey()
*
* To decrypt transaction
*
 */
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
	}
	return plaintext
}

/* EncryptPSS()
*
* To encrypt transaction
*
 */
func EncryptPSS (messageJson string, pubLicKey *rsa.PublicKey) ([]byte, hash.Hash, []byte, error){
	message := []byte(messageJson)
	label := []byte("")
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(
		hash,
		rand.Reader,
		pubLicKey,
		message,
		label,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return ciphertext, hash, label, err
}

/* DecryptPSS()
*
* To Decrypt transaction
*
 */
func DecryptPSS (ciphertext []byte, hash hash.Hash, label []byte,privateKey *rsa.PrivateKey) (string, error){

	plainText, err := rsa.DecryptOAEP(
		hash,
		rand.Reader,
		privateKey,
		ciphertext,
		label,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	plainTextJson := string(plainText)
	return plainTextJson, err
}

/* SignPSS()
*
* To sign transaction
*
 */
func SignPSS (message []byte, privateKey *rsa.PrivateKey) ([]byte, rsa.PSSOptions, []byte, crypto.Hash, error){
	//messageByte := []byte(message)
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	PSSmessage := message
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPSS(
		rand.Reader,
		privateKey,
		newhash,
		hashed,
		&opts,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return signature, opts, hashed, newhash, err
}


/* VerificationPSS()
*
* To verify signature
*
 */
func VerificationPSS (publicKey *rsa.PublicKey, opts rsa.PSSOptions, hashed []byte, newhash crypto.Hash, signature []byte) (bool,error){
	isVerify := false
	err := rsa.VerifyPSS(
		publicKey,
		newhash,
		hashed,
		signature,
		&opts,
	)
	if err != nil {
		fmt.Println("Verify Signature failed!!!")
		isVerify = false
		os.Exit(1)
	} else {
		fmt.Println("Verify Signature successful...")
		isVerify = true
	}
	return isVerify, err
}

/* GenerateKeyString()
*
* Generate Key string
*
 */
func GenerateKeyString() (*rsa.PrivateKey,VerificationKeyJson){

	privateKey, err := rsa.GenerateKey(rand.Reader, 2014)
	if err != nil {
		return privateKey,NewVerificationKeyJson("", "")
	}

	privateKeyDer := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDer,
	}
	privateKeyPem := string(pem.EncodeToMemory(&privateKeyBlock))

	publicKey := privateKey.PublicKey
	publicKeyDer, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return privateKey,NewVerificationKeyJson("", "")
	}

	publicKeyBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   publicKeyDer,
	}
	publicKeyPem := string(pem.EncodeToMemory(&publicKeyBlock))

	//fmt.Println(privateKeyPem)
	//fmt.Println(publicKeyPem)

	return privateKey,NewVerificationKeyJson(publicKeyPem, privateKeyPem)
}

/* VerificationKeyJson struct
*
* Generate data structure
*
 */
type VerificationKeyJson struct {
	PublicKey   string   `json:"publickey"` //DataStructure
	PrivateKey  string  `json:"privatekey"`
}

/* NewVerificationKeyJson()
*
* Generate new key
*
 */
func NewVerificationKeyJson(publicKey string, privateKey string) VerificationKeyJson {
	return VerificationKeyJson{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
}

/* ExportRsaPrivateKeyAsPemStr()
*
* to Export RsaPrivateKey As PemString
*
 */
func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {
	privkey_bytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkey_bytes,
		},
	)
	return string(privkey_pem)
}

/* ParseRsaPrivateKeyFromPemStr()
*
* to Parse RsaPrivate Key From PemString
*
 */
func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

/* ExportRsaPublicKeyAsPemStr()
*
* to Export RsaPublicKey as PemString
*
 */
func ExportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) (string, error) {
	pubkey_bytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return "", err
	}
	pubkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkey_bytes,
		},
	)

	return string(pubkey_pem), nil
}

/* ParseRsaPublicKeyFromPemStr()
*
* to Parse RsaPublicKey From PemString
*
 */
func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("Key type is not RSA")
}


/* EncryptPKCS()
*
* to encrypt transaction
*
 */
func EncryptPKCS (publickey *rsa.PublicKey, msg string) []byte{
	msgByte := []byte(msg)
	encryptedPKCS1v15, errPKCS1v15 := rsa.EncryptPKCS1v15(rand.Reader, publickey, msgByte)

	if errPKCS1v15 != nil {
		fmt.Println(errPKCS1v15)
		os.Exit(1)
	}

	fmt.Printf("PKCS1v15 encrypted [%s] to \n[%x]\n", string(msg), encryptedPKCS1v15)
	return encryptedPKCS1v15
}

/* DecryptPKCS
*
* to decrypt transaction
*
 */
func DecryptPKCS (privatekey *rsa.PrivateKey, encryptedPKCS1v15 []byte) []byte{
	decryptedPKCS1v15, err := rsa.DecryptPKCS1v15(rand.Reader, privatekey, encryptedPKCS1v15)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("PKCS1v15 decrypted [%x] to \n[%s]\n", encryptedPKCS1v15, decryptedPKCS1v15)
	fmt.Println()
	return decryptedPKCS1v15
}

/* VerifyPKCS()
*
* to sign transaction
*
 */
func SignPKCS (message string , privatekey *rsa.PrivateKey)(crypto.Hash,[]byte, []byte ){
	var h crypto.Hash
	messageByte := []byte(message)
	hash := md5.New()
	io.WriteString(hash, string(messageByte))
	hashed := hash.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privatekey, h, hashed)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("PKCS1v15 Signature : %x\n", signature)
	return h, hashed, signature
}

/* VerifyPKCS()
*
* to Verify signature
*
 */
func VerifyPKCS (publickey *rsa.PublicKey, h crypto.Hash, hashed []byte, signature []byte) (bool, error){
	verified := false
	err := rsa.VerifyPKCS1v15(publickey, h, hashed, signature)

	if err != nil {
		fmt.Println("VerifyPKCS1v15 failed")
		os.Exit(1)
		verified = false
	} else {
		fmt.Println("VerifyPKCS1v15 successful")
		verified = true
	}
	return verified, err
}
