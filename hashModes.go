package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
)

var HASHMODES = map[string]func(string, string, string) bool{
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

	"BCRYPT": BCRYPT,

	"SHA1":            SHA1,
	"SHA1PLAIN":       SHA1,
	"SHA1X1PLAIN":     SHA1,
	"SHA1(SALTPLAIN)": SHA1SALTPLAIN,
	"SHA1SALTPLAIN":   SHA1SALTPLAIN,
	"SHA1DASH":        SHA1DASH,

	"SHA256(PLAINSALT)": SHA256PLAINSALT,
}

func SHA1(hash, plain, salt string) bool {
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

func BCRYPT(hash, plain, salt string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

func MD5(h, plain, salt string) bool {
	plainHash := md5.New()
	io.WriteString(plainHash, plain)
	return h == fmt.Sprintf("%x", plainHash.Sum(nil))
}

func MD5PLAINSALT(h, plain, salt string) bool {
	plainHash := md5.New()
	io.WriteString(plainHash, plain+salt)
	return h == fmt.Sprintf("%x", plainHash.Sum(nil))
}

func MD5SALTPLAIN(h, plain, salt string) bool {
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
	hash := md5.New()
	io.WriteString(hash, plain)
	plainHash := fmt.Sprintf("%x", hash.Sum(nil))

	hash = md5.New()
	hashedPlainNSalt := plainHash + salt
	io.WriteString(hash, hashedPlainNSalt)

	return h == fmt.Sprintf("%x", hash.Sum(nil))
}

func MYBB(h, plain, salt string) bool {
	hash := md5.New()
	io.WriteString(hash, salt)
	saltHash := fmt.Sprintf("%x", hash.Sum(nil))

	hash = md5.New()
	hashedPlainNSalt := saltHash + plain
	io.WriteString(hash, hashedPlainNSalt)

	return h == fmt.Sprintf("%x", hash.Sum(nil))
}
