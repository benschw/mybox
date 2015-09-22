package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {

	// DO NOT PUT credentials in code for production usage!
	// see https://www.socketloop.com/tutorials/golang-setting-up-configure-aws-credentials-with-official-aws-sdk-go
	// on setting creds from environment or loading from file

	// the file location and load default profile
	creds := credentials.NewEnvCredentials()

	_, err := creds.Get()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	config := aws.NewConfig().WithRegion("us-east-1").WithEndpoint("s3.amazonaws.com").WithCredentials(creds)

	s3client := s3.New(config)

	bucketName := "fliglio" // <-- change this to your bucket name

	fileToUpload := "test.txt"

	file, err := os.Open(fileToUpload)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()

	buffer := make([]byte, size)
	// read file content to buffer
	file.Read(buffer)

	fileBytes := bytes.NewReader(buffer) // convert to io.ReadSeeker type

	fileType := http.DetectContentType(buffer)

	path := "/examplefolder/" + file.Name() // target file and location in S3

	params := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName), // required
		Key:           aws.String(path),       // required
		Body:          fileBytes,
		ContentLength: &size,
		ContentType:   aws.String(fileType),
		Metadata: map[string]*string{
			"Key": aws.String("MetadataValue"), //required
		},
		// see more at http://godoc.org/github.com/aws/aws-sdk-go/service/s3#S3.PutObject
	}

	log.Println("here")
	result, err := s3client.PutObject(params)
	log.Println("here")

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Printf("%+v", awsErr)
		} else {
			log.Printf("%+v", err)
		}
	}

	fmt.Printf("%+v", result)
}
