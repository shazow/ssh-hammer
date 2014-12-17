package main

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"io"

	"golang.org/x/crypto/ssh"
)

func MakeKey() (ssh.Signer, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2014)
	if err != nil {
		return nil, err
	}
	return ssh.NewSignerFromKey(key)
}

type Keyring struct {
	keys []ssh.Signer
}

func (r *Keyring) Key(i int) (ssh.PublicKey, error) {
	if i >= len(r.keys) {
		return nil, nil
	}
	return r.keys[i].PublicKey(), nil
}

func (r *Keyring) Sign(i int, rand io.Reader, data []byte) ([]byte, error) {
	if i >= len(r.keys) {
		return nil, errors.New("Keyring: Invalid key index")
	}
	sig, err := r.keys[i].Sign(rand, data)
	if err != nil {
		return nil, err
	}
	return sig.Blob, nil
}

func (r *Keyring) Add(key ssh.Signer) {
	r.keys = append(r.keys, key)
}

// Make keyring with num random keys in it.
func NewKeyring(num int) *Keyring {
	r := Keyring{}
	for i := 0; i < num; i++ {
		key, err := MakeKey()
		if err != nil {
			logger.Errorf("Failed to make key: %s", err)
			return &r
		}
		r.Add(key)
	}
	return &r
}
