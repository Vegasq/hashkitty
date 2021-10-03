package main

import (
	"archive/zip"
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Record struct {
	Hash     string
	Salt     string
	Plain    string
	HashType string
}

func getArchives() []string {
	hdoLocation := "E:\\HashesOrg Archive\\Leaks"
	var files []string
	err := filepath.Walk(hdoLocation, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".zip") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Faled to read files", err)
	}
	return files
}

func parseBrokenBcrypt(rightPartLine string) (string, string) {
	hash := rightPartLine[0:60]
	plain := rightPartLine[60 : len(rightPartLine)]
	return hash, plain
}

func readHash(line string) (string, string) {
	if strings.Contains(line, ":") {
		parts := strings.SplitN(line, ":", 2)
		return parts[0], parts[1]
	}

	if line[0] == '$' && line[1] == '2' && strings.Contains(line, ":") == false {
		return parseBrokenBcrypt(line)
	}

	return "", ""
}

func decodeHashcatHexPlain(hexPlain string) string {
	if strings.Contains(hexPlain, "$HEX[") {
		hexPlain = strings.Replace(hexPlain, "$HEX[", "", 1)
		hexPlain = strings.Replace(hexPlain, "]", "", 1)
		bPlain, err := hex.DecodeString(hexPlain)
		if err != nil {
			panic("Failed to unhex")
		}
		hexPlain = string(bPlain)
	}
	return hexPlain
}

func validate(hashType, hash, plain, salt string) bool {
	validator := HASHMODES[hashType]
	if validator != nil {
		return validator(hash, plain, salt)
	}
	return false
}

func parseLine(line string) {
	r := Record{}
	values := strings.SplitN(line, " ", 2)
	if len(values) != 2 {
		fmt.Println("Unexpected line", line, len(values))
		return
	}
	r.HashType = values[0]

	if r.HashType != "BCRYPT" {
		return
	}

	hashCont, leftovers := readHash(values[1])

	var saltCont, plainCont string
	for i := 0; i < len(leftovers); i++ {
		if leftovers[i] == ':' {
			plainCont += string(leftovers[i])
			saltCont += plainCont
			plainCont = ""
		} else {
			plainCont += string(leftovers[i])
		}
	}
	if len(saltCont) > 0 {
		saltCont = saltCont[0 : len(saltCont)-1]
	}
	plainCont = decodeHashcatHexPlain(plainCont)

	// Validate hash/salt/plain/mode + temp fix for double cols
	isValid := validate(r.HashType, hashCont, plainCont, saltCont)
	if isValid == false && len(saltCont) > 0 {
		isValid = validate(r.HashType, hashCont, plainCont, saltCont[0:len(saltCont)-1])
		if isValid {
			saltCont = saltCont[0 : len(saltCont)-1]
		}
	}

	if isValid == false {
		fmt.Println(r.HashType, "h", hashCont, "s", saltCont, "p", plainCont, "l", line)
	}

	r.Hash = hashCont
	r.Plain = plainCont
	r.Salt = saltCont

}

func sneakPeek(f *zip.File) bool {
	flReader, err := f.Open()
	defer flReader.Close()

	if err != nil {
		fmt.Println("Failed to open file from archive", f)
	}

	doubleCols := 0
	scanner2 := bufio.NewScanner(flReader)
	for i := 0; i < 1000; i++ {
		scanner2.Scan()
		line := scanner2.Text()
		fmt.Println(line)
		if strings.Contains(line, "::") {
			doubleCols += 1
		}
	}

	return doubleCols > 500
}

func readFile(f *zip.File) {
	flReader, err := f.Open()

	if err != nil {
		fmt.Println("Failed to open file from archive", f)
	}

	fmt.Printf("Contents of %s:\n", f.Name)
	scanner := bufio.NewScanner(flReader)
	for scanner.Scan() {
		parseLine(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	flReader.Close()
}

func readArchive(path string) {
	// Open a zip archive for reading.
	r, err := zip.OpenReader(path)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {
		readFile(f)
	}
}

func main() {
	//a, b := parseBrokenBcrypt("$2y$07$aa00x37f5mgo8krilsvhuebesmextTV6633fYMGFLrjBtuQkE4AWCphilips")
	//fmt.Println(a, b)
	//return
	files := getArchives()
	for i := 0; i < len(files); i++ {
		if i < 5 {
			continue
		}
		fmt.Println(files[i], "\n-----\n")
		readArchive(files[i])

		if i == 100 {
			return
		}
	}
}
