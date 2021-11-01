/*
Package algos

Collection of hash algorithms and their map to IDs in hashcat.
*/
package algos

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"github.com/GerardSoleCa/wordpress-hash-go"
	"github.com/matthewhartstonge/argon2"
	"golang.org/x/crypto/bcrypt"
	"io"
	"strings"
)

var HASHCATALGOS = map[uint]func(string, string, string) bool{
	0:    MD5,
	10:   MD5PLAINSALT,
	11:   MD5PLAINSALT,
	20:   MD5SALTPLAIN,
	100:  SHA1,
	400:  WORDPRESS,
	2611: VBULLETIN,
	2711: VBULLETIN,
	2811: MYBB,
	3200: BCRYPT,

	99001: ARGON2,
}

var HASHALGOS = map[string]func(string, string, string) bool{
	"IPB":                      IPB,
	"MD5(MD5(SALT)MD5(PLAIN))": IPB,
	"VBULLETIN":                VBULLETIN,
	"MD5(MD5(PLAIN)SALT)":      VBULLETIN,
	"MYBB":                     MYBB,
	"MD5":                      MD5,
	"MD5PLAIN":                 MD5,
	"MD5X1PLAIN":               MD5,
	"MD5(PLAINSALT)":           MD5PLAINSALT,
	"MD5(MD5(SALT)PLAIN)":      MYBB,
	"MD5SALTPLAIN":             MD5SALTPLAIN,
	"OSCOMMERCE":               MD5SALTPLAIN,
	"JOOMLA":                   MD5PLAINSALT,
	"WORDPRESS":                WORDPRESS,

	"BCRYPT": BCRYPT,
	"ARGON2": ARGON2,

	"SHA1":            SHA1,
	"SHA1PLAIN":       SHA1,
	"SHA1X1PLAIN":     SHA1,
	"SHA1(SALTPLAIN)": SHA1SALTPLAIN,
	"SHA1SALTPLAIN":   SHA1SALTPLAIN,
	"SHA1DASH":        SHA1DASH,

	"SHA256(PLAINSALT)": SHA256PLAINSALT,
}

func ARGON2(hash, plain, _ string) bool {
	ok, _ := argon2.VerifyEncoded([]byte(plain), []byte(hash))
	return ok
}

func WORDPRESS(hash, plain, _ string) bool {
	return wphash.CheckPassword(plain, hash)
}

func SHA1(hash, plain, _ string) bool {
	h := sha1.New()
	io.WriteString(h, fmt.Sprintf("%s", plain))
	return hash == fmt.Sprintf("%x", h.Sum(nil))
}

func SHA1DASH(hash, plain, salt string) bool {
	h := sha1.New()
	io.WriteString(h, fmt.Sprintf("--%s--%s--", salt, plain))
	return hash == fmt.Sprintf("%x", h.Sum(nil))
}

func SHA1SALTPLAIN(hash, plain, salt string) bool {
	h := sha1.New()
	io.WriteString(h, fmt.Sprintf("%s%s", salt, plain))
	return hash == fmt.Sprintf("%x", h.Sum(nil))
}

func BCRYPT(hash, plain, _ string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

func MD5(h, plain, _ string) bool {
	h = strings.ToLower(h)
	plainHash := md5.New()
	io.WriteString(plainHash, plain)
	return h == fmt.Sprintf("%x", plainHash.Sum(nil))
}

func MD5PLAINSALT(h, plain, salt string) bool {
	h = strings.ToLower(h)
	plainHash := md5.New()
	io.WriteString(plainHash, plain+salt)
	return h == fmt.Sprintf("%x", plainHash.Sum(nil))
}

func MD5SALTPLAIN(h, plain, salt string) bool {
	h = strings.ToLower(h)
	plainHash := md5.New()
	io.WriteString(plainHash, salt+plain)
	return h == fmt.Sprintf("%x", plainHash.Sum(nil))
}

func SHA256PLAINSALT(h, plain, salt string) bool {
	plainHash := sha256.New()
	io.WriteString(plainHash, plain+salt)
	return h == fmt.Sprintf("%x", plainHash.Sum(nil))
}

func IPB(h, plain, salt string) bool {
	// md5(md5($salt).md5($pass))
	saltHash := md5.New()
	io.WriteString(saltHash, salt)

	plainHash := md5.New()
	io.WriteString(plainHash, plain)

	saltPlainCombo := fmt.Sprintf("%x%x", saltHash.Sum(nil), plainHash.Sum(nil))
	saltPlainHash := md5.New()
	io.WriteString(saltPlainHash, saltPlainCombo)

	return h == fmt.Sprintf("%x", saltPlainHash.Sum(nil))
}

func VBULLETIN(h, plain, salt string) bool {
	// md5(md5($pass).$salt)
	hash := md5.New()
	io.WriteString(hash, plain)
	plainHash := fmt.Sprintf("%x", hash.Sum(nil))

	hash = md5.New()
	hashedPlainNSalt := plainHash + salt
	io.WriteString(hash, hashedPlainNSalt)

	return h == fmt.Sprintf("%x", hash.Sum(nil))
}

func MYBB(h, plain, salt string) bool {
	// Can't find such algo in hashcat?
	// md5(md5($salt).$pass)

	hash := md5.New()
	io.WriteString(hash, salt)
	saltHash := fmt.Sprintf("%x", hash.Sum(nil))

	hash = md5.New()
	hashedPlainNSalt := saltHash + plain
	io.WriteString(hash, hashedPlainNSalt)

	return h == fmt.Sprintf("%x", hash.Sum(nil))
}
