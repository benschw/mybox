package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

func main() {
	auth, err := aws.EnvAuth()
	if err != nil {
		log.Fatal(err)
	}
	client := s3.New(auth, aws.USEast)

	resp, err := client.ListBuckets()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Print(fmt.Sprintf("%T %+v", resp.Buckets[0], resp.Buckets[0]))

	bucket := client.Bucket("fliglio")

	//data := []byte("Hello World")
	data, fileType, path, err := foo("test.txt")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = bucket.Put(path, data, fileType, s3.BucketOwnerFull)
	if err != nil {
		log.Fatal(err)
		return
	}

	out, err := bucket.Get(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("%s", out)

}

func foo(fileToUpload string) ([]byte, string, string, error) {

	file, err := os.Open(fileToUpload)
	if err != nil {
		return nil, "", "", err
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()

	buffer := make([]byte, size)
	file.Read(buffer)

	fileType := http.DetectContentType(buffer)

	return buffer, fileType, file.Name(), nil
}
