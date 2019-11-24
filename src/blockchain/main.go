	package main

	import (
		//"MerklePatriciaTree/p5/Blockchain_Application_P5/data"
		//"MerklePatriciaTree/p5/Blockchain_Application_P5/p5"
		"crypto/rand"
		"crypto/rsa"
		"fmt"
		"github.com/edgexfoundry/device-simple/src/blockchain/data"
		p5 "github.com/edgexfoundry/device-simple/src/blockchain/transaction"
		"log"
		"net/http"
		"os"
	)

	func main() {

		//NimaKey := data.GenerateKeyFirst()
		//RozitaKey := data.GenerateKeyFirst()
		//
		//fmt.Println("Private Key : ", RozitaKey.PrivateKey)
		//fmt.Println("Public key ", RozitaKey.PublicKey)
		//fmt.Println("Private Key : ", NimaKey.PrivateKey)
		//fmt.Println("Public key ", NimaKey.PublicKey)
		//
		//
		//message := "Hi I am Rozita !!!!!"
		//cipherTexttoNima, hash, label, _:= data.Encrypt(message, NimaKey.PublicKey)
		//fmt.Println("cipherTexttoNima is:", cipherTexttoNima )
		//
		//signature, opts, hashed, newhash, _:= data.Sign(cipherTexttoNima, RozitaKey.PrivateKey)
		//fmt.Println("Rozita Signature is:", signature)
		//
		//plainTextfromRozita, _ := data.Decrypt(cipherTexttoNima, hash , label ,NimaKey.PrivateKey)
		//fmt.Println("plainTextfrom Rozita is:", plainTextfromRozita)
		//
		//isVerified, _ := data.Verification (RozitaKey.PublicKey, opts, hashed, newhash, signature)
		//fmt.Println("Is Verified is:", isVerified)

		// generate private key
		privatekey, err := rsa.GenerateKey(rand.Reader, 1024)

		if err != nil {
			fmt.Println(err.Error)
			os.Exit(1)
		}

		privatekey.Precompute()
		err = privatekey.Validate()
		if err != nil {
			fmt.Println(err.Error)
			os.Exit(1)
		}

		var publickey *rsa.PublicKey
		publickey = &privatekey.PublicKey

		msg := "The secret message!"

		// EncryptPKCS1v15
		encryptedPKCS1v15 := data.EncryptPKCS(publickey, msg)

		// DecryptPKCS1v15
		decryptedPKCS1v15 := data.DecryptPKCS(privatekey , encryptedPKCS1v15)
		fmt.Printf("PKCS1v15 decrypted [%x] to \n[%s]\n", encryptedPKCS1v15, decryptedPKCS1v15)
		// SignPKCS1v15
		message := "message to be signed"
		h, hashed, signature := data.SignPKCS(message , privatekey)
		fmt.Printf("PKCS1v15 Signature : %x\n", signature)
		//VerifyPKCS1v15
		verified, err := data.VerifyPKCS(publickey, h, hashed, signature)
		fmt.Println("Signature verified: ", verified)



		router := p5.NewRouter()
		if len(os.Args) > 1 {
			log.Fatal(http.ListenAndServe(":" + os.Args[1], router))
		} else {
			log.Fatal(http.ListenAndServe(":6686", router))
		}

	}

