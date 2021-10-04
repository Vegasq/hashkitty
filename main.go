package main

import (
	"archive/zip"
	"bufio"
	"encoding/hex"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"hashkitty/modes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Record struct {
	Hash     string
	Salt     string
	Plain    string
	HashType string
	RAW      string

	Origin    string
	AddedDate time.Time
}

func getArchives() []string {
	hdoLocation, err := os.ReadFile("archives.location")
	if err != nil {
		panic("Failed to read connection uri")
	}

	var files []string
	err = filepath.Walk(string(hdoLocation), func(path string, info os.FileInfo, err error) error {
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
	var plain = rightPartLine[60:]
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
	validator := modes.HASHMODES[hashType]
	if validator != nil {
		return validator(hash, plain, salt)
	}
	return false
}

func parseLine(ch *chan Record, line string, path string) {
	r := Record{}
	values := strings.SplitN(line, " ", 2)
	if len(values) != 2 {
		fmt.Println("Unexpected line", line, len(values))
		return
	}
	r.HashType = values[0]

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
	//isValid := validate(r.HashType, hashCont, plainCont, saltCont)
	//if isValid == false && len(saltCont) > 0 {
	//	isValid = validate(r.HashType, hashCont, plainCont, saltCont[0:len(saltCont)-1])
	//	if isValid {
	//		saltCont = saltCont[0 : len(saltCont)-1]
	//	}
	//}
	//
	//if isValid == false {
	//	fmt.Println(r.HashType, "h", hashCont, "s", saltCont, "p", plainCont, "l", line)
	//}

	r.Hash = hashCont
	r.Plain = plainCont
	r.Salt = saltCont
	r.RAW = line

	r.Origin = path
	r.AddedDate = time.Now().UTC()
	*ch <- r
}

func readFile(ch *chan Record, f *zip.File, path string) {
	flReader, err := f.Open()

	if err != nil {
		fmt.Println("Failed to open file from archive", f)
	}

	fmt.Printf("Contents of %s:\n", f.Name)
	scanner := bufio.NewScanner(flReader)
	for scanner.Scan() {
		parseLine(ch, scanner.Text(), path)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	err = flReader.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func readArchive(ch *chan Record, path string) {
	r, err := zip.OpenReader(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func(r *zip.ReadCloser) {
		err := r.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(r)

	for _, f := range r.File {
		readFile(ch, f, path)
	}
}

func recordsReader(ch *chan Record, mh *MongoHolder) {
	var rec Record
	var batch []Record
	for {
		rec = <-*ch
		batch = append(batch, rec)
		if len(batch) >= 999_000 {
			mh.Insert(batch)
			//recordsWriter(mh, batch)
			batch = []Record{}
		}
	}
}

type Processed struct {
	Path string
}

func main() {
	records := make(chan Record)
	mh := NewMongoHolder()
	go recordsReader(&records, mh)
	files := getArchives()
	for i := 0; i < len(files); i++ {
		cur, err := mh.PROCESSEDB.Find(*mh.Context, bson.M{"path": files[i]})
		if err != nil {
			fmt.Println(err)
		}
		exist := cur.Next(*mh.Context)
		if exist == false {
			fmt.Println(files[i])
			readArchive(&records, files[i])
			mh.PROCESSEDB.InsertOne(*mh.Context, Processed{files[i]})
		} else {
			fmt.Println("SKIP", files[i])
			//return
		}
	}
}
