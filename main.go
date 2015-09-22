package main

import (
	"flag"
	"log"

	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
)

var (
	bucketName string
	fileName   string
)

func init() {
	flag.StringVar(&bucketName, "b", "fliglio", "Bucket Name")
}

func main() {

	flag.Parse()

	// The AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables are used.
	auth, err := aws.EnvAuth()
	if err != nil {
		log.Printf("%s", err)
		return
	}
	// Open Bucket
	s := s3.New(auth, aws.USEast)
	bucket := s.Bucket(bucketName)

	data := []byte("Hello, World")
	err = bucket.Put("/sample.txt", data, "text/plain", s3.BucketOwnerFull)
	if err != nil {
		log.Printf("%+v", err)
	}
}
