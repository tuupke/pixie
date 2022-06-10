package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"strconv"
	"testing"

	"go.dedis.ch/kyber/v3/pairing/bn256"
	"go.dedis.ch/kyber/v3/sign/bls"
	"go.dedis.ch/kyber/v3/util/random"
)

var (
	keySizes  = []int{128, 256, 512, 1024, 2048, 4096}
	plaintext = make([]byte, 1024)
)

func BenchmarkBLSKeyCreate(b *testing.B) {
	b.StopTimer()
	suite := bn256.NewSuite()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		private, public := bls.NewKeyPair(suite, random.New())

		private.MarshalBinary()
		public.MarshalBinary()
	}
}

func BenchmarkRSAKeyCreate(b *testing.B) {
	for _, bitLength := range keySizes {
		b.Run(strconv.Itoa(bitLength), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				privateKey, err := rsa.GenerateKey(rand.Reader, bitLength)
				if err != nil {
					b.Fatal(err)
				}

				_ = privateKey.PublicKey
			}
		})
	}
}

func BenchmarkBSLSign(b *testing.B) {
	b.StopTimer()
	suite := bn256.NewSuite()
	private, _ := bls.NewKeyPair(suite, random.New())
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := bls.Sign(suite, private, plaintext)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBLSVerify(b *testing.B) {
	b.StopTimer()
	suite := bn256.NewSuite()
	private, public := bls.NewKeyPair(suite, random.New())
	sign, err := bls.Sign(suite, private, plaintext)
	if err != nil {
		b.Fatal(err)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		err := bls.Verify(suite, public, plaintext, sign)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRSASign(b *testing.B) {
	hashed := sha256.Sum256(plaintext)
	for _, bitLength := range keySizes {
		privateKey, err := rsa.GenerateKey(rand.Reader, bitLength)
		if err != nil {
			b.Fatal(err)
		}

		b.Run(strconv.Itoa(bitLength), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkRSAVerify(b *testing.B) {
	hashed := sha256.Sum256(plaintext)

	for _, bitLength := range keySizes {
		privateKey, err := rsa.GenerateKey(rand.Reader, bitLength)
		if err != nil {
			b.Log(err)
			continue
		}

		sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
		if err != nil {
			b.Log(err)
			continue
		}

		public := &privateKey.PublicKey
		b.Run(strconv.Itoa(bitLength), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				err := rsa.VerifyPKCS1v15(public, crypto.SHA256, hashed[:], sign)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkRSASignatureLength(b *testing.B) {
	hashed := sha256.Sum256(plaintext)
	for _, bitLength := range keySizes {

		b.Run(strconv.Itoa(bitLength), func(b *testing.B) {

			privateKey, err := rsa.GenerateKey(rand.Reader, bitLength)
			if err != nil {
				b.Fatal(err)
			}
			sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
			if err != nil {
				b.Fatal(err)
			}

			b.Log(len(sign))
		})
	}
}
