package user_repository

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

type EncryptionOptions struct {
	Key string
}

func encrypt(stringToEncrypt string, keyString string) (string, error) {
	//Since the key is in string, we need to convert decode it to bytes
	key, err := hex.DecodeString(keyString)
	if err != nil {
		return "", err
	}
	plaintext := []byte(stringToEncrypt)
	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	//WithEncryption the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

func decrypt(encryptedString string, keyString string) (string, error) {
	key, err := hex.DecodeString(keyString)
	if err != nil {
		return "", err
	}
	enc, err := hex.DecodeString(encryptedString)
	if err != nil {
		return "", err
	}
	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	//Get the nonce size
	nonceSize := aesGCM.NonceSize()
	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", plaintext), nil
}

func (user *User) encryptUser(key string) error {
	if user == nil {
		return errors.New("user is nil")
	}
	if key == "" {
		return errors.New("key is empty")
	}
	if user.Name != "" {
		encName, err := encrypt(user.Name, key)
		if err != nil {
			return err
		}
		user.Name = encName
	}
	if user.Gender != "" {
		encGender, err := encrypt(user.Gender, key)
		if err != nil {
			return err
		}
		user.Gender = encGender
	}
	if user.Email != "" {
		encEmail, err := encrypt(user.Email, key)
		if err != nil {
			return err
		}
		user.Email = encEmail
	}
	if user.Role != "" {
		encRole, err := encrypt(user.Role, key)
		if err != nil {
			return err
		}
		user.Role = encRole
	}
	if user.Country != "" {
		encCountry, err := encrypt(user.Country, key)
		if err != nil {
			return err
		}
		user.Country = encCountry
	}
	if user.Image != "" {
		encImage, err := encrypt(user.Image, key)
		if err != nil {
			return err
		}
		user.Image = encImage
	}
	if user.Birthdate != "" {
		encBirthdate, err := encrypt(user.Birthdate, key)
		if err != nil {
			return err
		}
		user.Birthdate = encBirthdate
	}
	return nil
}

func (user *User) decryptUser(key string) error {
	if user == nil {
		return errors.New("user is nil")
	}
	if key == "" {
		return errors.New("key is empty")
	}
	if user.Name != "" {
		decName, err := decrypt(user.Name, key)
		if err != nil {
			return err
		}
		user.Name = decName
	}
	if user.Gender != "" {
		decGender, err := decrypt(user.Gender, key)
		if err != nil {
			return err
		}
		user.Gender = decGender
	}
	if user.Email != "" {
		decEmail, err := decrypt(user.Email, key)
		if err != nil {
			return err
		}
		user.Email = decEmail
	}
	if user.Role != "" {
		decRole, err := decrypt(user.Role, key)
		if err != nil {
			return err
		}
		user.Role = decRole
	}
	if user.Country != "" {
		decCountry, err := decrypt(user.Country, key)
		if err != nil {
			return err
		}
		user.Country = decCountry
	}
	if user.Image != "" {
		decImage, err := decrypt(user.Image, key)
		if err != nil {
			return err
		}
		user.Image = decImage
	}
	if user.Birthdate != "" {
		decBirthdate, err := decrypt(user.Birthdate, key)
		if err != nil {
			return err
		}
		user.Birthdate = decBirthdate
	}
	return nil
}
