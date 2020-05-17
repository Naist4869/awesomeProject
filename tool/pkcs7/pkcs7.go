package pkcs7

import "fmt"

const (
	BlockSize  = 32            // PKCS#7
	BLOCK_MASK = BlockSize - 1 // BLOCK_SIZE 为 2^n 时, 可以用 mask 获取针对 BLOCK_SIZE 的余数
)

func Pad(x []byte) []byte {
	numPadBytes := 32 - len(x)%32
	padByte := byte(numPadBytes)
	tmp := make([]byte, len(x)+numPadBytes)
	copy(tmp, x)
	for i := 0; i < numPadBytes; i++ {
		tmp[len(x)+i] = padByte
	}
	return tmp
}

func Unpad(x []byte) ([]byte, error) {
	// last byte is number of suffix bytes to remove
	n := int(x[len(x)-1])
	if n < 1 || n > BlockSize {
		return nil, fmt.Errorf("the amount to pad is incorrect: %d", n)
	}
	return x[:len(x)-n], nil
}
