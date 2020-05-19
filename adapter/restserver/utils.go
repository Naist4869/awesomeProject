package restserver

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/big"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/Naist4869/awesomeProject/tool/pkcs7"
)

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(finalPath) != '/'
	if appendSlash {
		return finalPath + "/"
	}
	return finalPath
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			//debugPrint("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		//debugPrint("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too much parameters")
	}
}

// cleanPath is the URL version of path.Clean, it returns a canonical URL path
// for p, eliminating . and .. elements.
//
// The following rules are applied iteratively until no further processing can
// be done:
//	1. Replace multiple slashes with a single slash.
//	2. Eliminate each . path name element (the current directory).
//	3. Eliminate each inner .. path name element (the parent directory)
//	   along with the non-.. element that precedes it.
//	4. Eliminate .. elements that begin a rooted path:
//	   that is, replace "/.." by "/" at the beginning of a path.
//
// If the result of this process is an empty string, "/" is returned.
func cleanPath(p string) string {
	// Turn empty string into "/"
	if p == "" {
		return "/"
	}

	n := len(p)
	var buf []byte

	// Invariants:
	//      reading from path; r is index of next byte to process.
	//      writing to buf; w is index of next byte to write.

	// path must start with '/'
	r := 1
	w := 1

	if p[0] != '/' {
		r = 0
		buf = make([]byte, n+1)
		buf[0] = '/'
	}

	trailing := n > 1 && p[n-1] == '/'

	// A bit more clunky without a 'lazybuf' like the path package, but the loop
	// gets completely inlined (bufApp). So in contrast to the path package this
	// loop has no expensive function calls (except 1x make)

	for r < n {
		switch {
		case p[r] == '/':
			// empty path element, trailing slash is added after the end
			r++

		case p[r] == '.' && r+1 == n:
			trailing = true
			r++

		case p[r] == '.' && p[r+1] == '/':
			// . element
			r += 2

		case p[r] == '.' && p[r+1] == '.' && (r+2 == n || p[r+2] == '/'):
			// .. element: remove to last /
			r += 3

			if w > 1 {
				// can backtrack
				w--

				if buf == nil {
					for w > 1 && p[w] != '/' {
						w--
					}
				} else {
					for w > 1 && buf[w] != '/' {
						w--
					}
				}
			}

		default:
			// real path element.
			// add slash if needed
			if w > 1 {
				bufApp(&buf, p, w, '/')
				w++
			}

			// copy element
			for r < n && p[r] != '/' {
				bufApp(&buf, p, w, p[r])
				w++
				r++
			}
		}
	}

	// re-append trailing slash
	if trailing && w > 1 {
		bufApp(&buf, p, w, '/')
		w++
	}

	if buf == nil {
		return p[:w]
	}
	return string(buf[:w])
}

// internal helper to lazily create a buffer if necessary.
func bufApp(buf *[]byte, s string, w int, c byte) {
	if *buf == nil {
		if s[w] == c {
			return
		}

		*buf = make([]byte, len(s))
		copy(*buf, s[:w])
	}
	(*buf)[w] = c
}

// Sign 微信公众号 url 签名. https://www.programming-books.io/essential/go/hex-base64-encoding-b6e9fbb3165c4bcb907e469d86783aab
func Sign(strs ...string) (signature string) {
	sort.Strings(strs)
	tmpstr := strings.Join(strs, "")
	signature = fmt.Sprintf("%x", sha1.Sum([]byte(tmpstr)))
	return
}

// random  当前消息加密所用的 random, 16-bytes  cap为20 所以可以取出msgLen为[16:]
// rawXMLMsg  消息的明文文本, xml格式
// appId 当前消息加密所用的 AppId
func AESDecryptMsg(base64Msg []byte, aesKey []byte) (random, rawXMLMsg, appId []byte, err error) {
	var encryptedMsgLen int
	encryptedMsg := make([]byte, base64.StdEncoding.DecodedLen(len(base64Msg)))

	encryptedMsgLen, err = base64.StdEncoding.Decode(encryptedMsg, base64Msg)
	if err != nil {
		err = fmt.Errorf("AESDecryptMsg Decode base64Msg fail: %w", err)
		return
	}
	cipherText := encryptedMsg[:encryptedMsgLen]
	if len(cipherText) < pkcs7.BlockSize {
		err = fmt.Errorf("the length of ciphertext too short: %d", len(cipherText))
		return
	}
	if len(cipherText)&pkcs7.BlockMask != 0 {
		err = fmt.Errorf("ciphertext is not a multiple of the block size, the length is %d", len(cipherText))
		return
	}
	plaintext := make([]byte, len(cipherText)) // len(plaintext) >= BLOCK_SIZE
	// 解密
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		err = fmt.Errorf("AESDecryptMsg %w", err)
		return
	}
	mode := cipher.NewCBCDecrypter(block, aesKey[:16])
	mode.CryptBlocks(plaintext, cipherText)
	plaintext, err = pkcs7.Unpad(plaintext)
	if err != nil {
		err = fmt.Errorf("AESDecryptMsg %w", err)
		return
	}
	// 反拼接
	// len(plaintext) == 16+4+len(rawXMLMsg)+len(appId)
	if len(plaintext) <= 20 {
		err = fmt.Errorf("plaintext too short, the length is %d", len(plaintext))
		return
	}
	msgLen := binary.BigEndian.Uint32(plaintext[16:20])
	appIdOffset := 20 + msgLen
	if len(plaintext) <= int(appIdOffset) {
		err = fmt.Errorf("msg length too large: %d", msgLen)
		return
	}
	random = plaintext[:16:20]
	rawXMLMsg = plaintext[20:appIdOffset:appIdOffset]
	appId = plaintext[appIdOffset:]
	return
}

//AESEncryptMsg cipherText = AES_Encrypt[random(16B) + msg_len(4B) + rawXMLMsg + appId]
func AESEncryptMsg(random, rawXMLMsg []byte, appId string, aesKey []byte) (cipherText string, err error) {
	appIdOffset, plaintext := pkcs7.Pad(rawXMLMsg, appId)
	copy(plaintext[:16], random)
	binary.BigEndian.PutUint32(plaintext[16:20], uint32(len(rawXMLMsg)))
	copy(plaintext[20:], rawXMLMsg)
	copy(plaintext[appIdOffset:], appId)
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		err = fmt.Errorf("AESEncryptMsg %w", err)
		return
	}
	cipherTextBytes := make([]byte, len(plaintext))

	mode := cipher.NewCBCEncrypter(block, aesKey[:16])
	mode.CryptBlocks(cipherTextBytes, plaintext)
	cipherText = base64.StdEncoding.EncodeToString(cipherTextBytes)

	return
}

func makeNonce() (nonce string) {
	limit := big.NewInt(1)
	limit = limit.Lsh(limit, 64)
	n, err := rand.Int(rand.Reader, limit)
	if err != nil {
		return
	}
	nonce = n.String()
	return
}
