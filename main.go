package main

import (
	"fmt"
	"os"
	"crypto/sha256"
	"bytes"
)

const (
	BLOCK_SIZE = 32
	NUMBER_BLOCKS = 256
	HASH_SIZE = 32
)

// SecretKey defines the structure of secret key in the Lamport signature scheme
type SecretKey struct{
	Key			[2][NUMBER_BLOCKS * BLOCK_SIZE]byte
}

// PublicKey defines the structure of public key in the Lamport signature scheme
type PublicKey struct{
	Key			[512 * BLOCK_SIZE]byte
}

// Signature defines the structure of a signature in the Lamport signature scheme
type Signature struct{
	Signature	[256 * NUMBER_BLOCKS]byte
}


// CheckError function receives an error type and checks wether the function returned an error
// or not
func CheckError (e error){
	if e != nil {
		panic (e)
	}
}

// GetHash function returns the sha256 hash of the message received as argument
func GetHash (byteArray []byte) [HASH_SIZE]byte{
	if byteArray == nil {
		panic ("Unable to hash null message")
	}
	return sha256.Sum256(byteArray[:])
}

// GenerateKeys does not receives arguments and returns a KeyPair structure containing a key pair
// to perform Lamport signature scheme
func GenerateKeys () (SecretKey, PublicKey){

	var secretKey [2][NUMBER_BLOCKS * BLOCK_SIZE]byte
	var publicKey []byte
	var publicKey16K [512 * BLOCK_SIZE]byte
	var index uint16
	var index2 uint8 

	// In Lamport's signature scheme, the user generates 256 pairs of random number of 32 bytes to
	// use as the secret key. To do that, we read file /dev/urandom to get random bytes that will
	// be used as secret key
	fileRandomness, errorException := os.Open("/dev/urandom")
	CheckError(errorException)

	for index2 = 0; index2 < 2; index2 ++{
		fileRandomness.Read(secretKey[index2][:])
	}

	// In Lamport's signature scheme, the publiic key is generated by hashing each of the 512 bits
	// number generated
	for index2 = 0; index2 < 2; index2 ++{
		for index = 0; index < 256; index ++{
			blockHash := sha256.Sum256(secretKey[index2][(index*32):(index*32 + 32)])
			publicKey = append(publicKey[:], blockHash[:]...)
		}
	}

	fileRandomness.Close()

	copy (publicKey16K[:], publicKey)
	SecretKey := SecretKey{secretKey}
	PublicKey := PublicKey{publicKey16K}

	return SecretKey, PublicKey

}

// Sign function receives the message to be signed, the key-pair to sign the message and returns
// a Lamport signature
func Sign (message string, secretKey SecretKey) Signature{
	var byteIterator uint8
	var bitMask uint8
	var secretKeyIterator uint16
	var signature []byte

	if message == ""{
		panic ("No message to sign")
	}
	messageBytes := []byte(message)
	messageHash := GetHash(messageBytes)

	byteIterator = 0
	
	// secretKeyIterator iterates through the bytes of the secret key to generate the Lamport
	// signature. The variable secretKeyIterator iterates through blocks in the secret key.
	// Hence, secretKeyIterator increments 32 in size (1 block size) for every bit int the 
	// message hash. Thus:
	// byteIterator[0] is from bit 0 to bit 7
	// bit 0 -> defines byte 0 to 31
	// bit 1 -> defines byte 32 to 63
	// ...
	// byteIterator[1] is from bit 8 to bit 15
	// bit 8 -> defines byte 256 to 287
	secretKeyIterator = 0

	for byteIterator = 0; byteIterator < 32; byteIterator ++{
		for bitMask = 1; bitMask <= 128; bitMask = bitMask * 2{
			if (bitMask & messageHash[byteIterator]) == 0{
				signature = append(signature[:], secretKey.Key[0][secretKeyIterator:(secretKeyIterator + 32)]...)
			}else {
				signature = append(signature[:], secretKey.Key[1][secretKeyIterator:(secretKeyIterator + 32)]...)
			}
			secretKeyIterator += 32
			if bitMask == 128{
				break;
			}
		}
	}

	Signature := Signature {}
	copy (Signature.Signature[:], signature)
	return Signature
	
}

// Verify function receives a message, a Lamport's signature and the public key of the signer
// and verifies wether the provided signature is valid or not.
func Verify (message string, signature Signature, publicKey PublicKey) bool{
	var byteIterator uint8
	var bitMask uint8
	var signatureIterator uint16

	if message == ""{
		panic ("No message to verify")
	}

	messageBytes := []byte(message)
	messageHash := GetHash(messageBytes)

	byteIterator = 0
	bitMask = 0
	signatureIterator = 0
	for byteIterator = 0; byteIterator < 32; byteIterator ++{
		for bitMask = 1; bitMask <= 128; bitMask = bitMask * 2{
			if (bitMask & messageHash[byteIterator]) == 0{
				signatureBlockHash := GetHash(signature.Signature[signatureIterator:(signatureIterator+32)])
				if bytes.Compare(publicKey.Key[signatureIterator:(signatureIterator + 32)], signatureBlockHash[:]) != 0{
					return false
				}
			}else {
				signatureBlockHash := GetHash(signature.Signature[signatureIterator:(signatureIterator+32)])
				if bytes.Compare(publicKey.Key[signatureIterator + 8192:(signatureIterator + 32) + 8192], signatureBlockHash[:]) != 0{
					return false
				}
			}
			signatureIterator += 32
			if bitMask == 128{
				break;
			}
		}
	} 

	return true
}

func main () {
	SecretKey, PublicKey := GenerateKeys()
	Signature := Sign("Hi",SecretKey)
	Verify ("Hi", Signature, PublicKey)
}