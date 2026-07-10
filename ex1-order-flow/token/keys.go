package token

import (
	"crypto/ed25519"
	"crypto/sha256"
)

func PrivateKey(role string) ed25519.PrivateKey {
	sum := sha256.Sum256([]byte("ex1-order-flow:" + role))
	return ed25519.NewKeyFromSeed(sum[:])
}

func PublicKey(role string) ed25519.PublicKey {
	return PrivateKey(role).Public().(ed25519.PublicKey)
}
