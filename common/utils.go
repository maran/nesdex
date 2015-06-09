package common

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/user"
)

const filechunk = 8192

func CalcMd5(path string) string {
	file, err := os.Open(path)

	if err != nil {
		panic(err.Error())
	}

	defer file.Close()

	// calculate the file size
	info, _ := file.Stat()

	filesize := info.Size()

	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))

	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)

		file.Read(buf)
		io.WriteString(hash, string(buf)) // append into the hash
	}
	hashStr := fmt.Sprintf("%x", hash.Sum(nil))
	return hashStr
}

func GetUserDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal("Could not get user dir", err)
	}
	return usr.HomeDir
}
