package repository

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"go-blockchain/wallet/domain/entity"
	"go-blockchain/wallet/domain/repository"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type walletRepository struct{}

func NewWalletRepository() repository.WalletRepository {
	return &walletRepository{}
}

func NewWallet() *entity.Wallet {
	// 1. Creating ECDSA private key (32 bytes) public key (64 bytes)
	w := new(entity.Wallet)
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.PrivateKey = privateKey
	w.PublicKey = &w.PrivateKey.PublicKey
	// 2. Perform SHA-256 hashing on the public key (32 bytes).
	h2 := sha256.New()
	h2.Write(w.PublicKey.X.Bytes())
	h2.Write(w.PublicKey.Y.Bytes())
	digest2 := h2.Sum(nil)
	// 3. Perform RIPEMD-160 hashing on the result of SHA-256 (20 bytes).
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)
	// 4. Add version byte in front of RIPEMD-160 hash (0x00 for Main Network).
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])
	// 5. Perform SHA-256 hash on the extended RIPEMD-160 result.
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)
	// 6. Perform SHA-256 hash on the result of the previous SHA-256 hash.
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)
	// 7. Take the first 4 bytes of the second SHA-256 hash for checksum.
	chsum := digest6[:4]
	// 8. Add the 4 checksum bytes from 7 at the end of extended RIPEMD-160 hash from 4 (25 bytes).
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], chsum[:])
	// 9. Convert the result from a byte string into base58.
	address := base58.Encode(dc8)
	w.BlockchainAddress = address
	return w
}

func (wr *walletRepository) PrivateKey(w *entity.Wallet) *ecdsa.PrivateKey {
	return w.PrivateKey
}

func (wr *walletRepository) PrivateKeyStr(w *entity.Wallet) string {
	return fmt.Sprintf("%x", w.PrivateKey.D.Bytes())
}

func (wr *walletRepository) PublicKey(w *entity.Wallet) *ecdsa.PublicKey {
	return w.PublicKey
}

func (wr *walletRepository) PublicKeyStr(w *entity.Wallet) string {
	return fmt.Sprintf("%x%x", w.PublicKey.X.Bytes(), w.PublicKey.Y.Bytes())
}

func (wr *walletRepository) BlockchainAddress(w *entity.Wallet) string {
	return w.BlockchainAddress
}

func (wr *walletRepository) MarshalJSON(w *entity.Wallet) ([]byte, error) {
	return json.Marshal(struct {
		PrivateKey        string `json:"private_key"`
		PublicKey         string `json:"public_key"`
		BlockchainAddress string `json:"blockchain_address"`
	}{
		PrivateKey:        wr.PrivateKeyStr(w),
		PublicKey:         wr.PublicKeyStr(w),
		BlockchainAddress: wr.BlockchainAddress(w),
	})
}
