package generate_password

import (
	"crypto/rand"
	"errors"
	"math/big"
)

var (
	ErrBadLen   = errors.New("incorrect pass length")
	ErrGenerate = errors.New("something went wrong while generate pass")
)

// letters - список допустимых символов в пароле
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GeneratePassword(n int) (string, error) {
	if n < 0 { // хотел <= 0 поставить, но в тестах прописано что длина 0 - ок
		return "", ErrBadLen
	}

	pass := make([]byte, 0, n)
	idxMax := big.NewInt(int64(len(letters)))

	for range n {
		idx, err := rand.Int(rand.Reader, idxMax)
		if err != nil {
			return "", ErrGenerate
		}

		pass = append(pass, letters[idx.Int64()])
	}

	return string(pass), nil
}
