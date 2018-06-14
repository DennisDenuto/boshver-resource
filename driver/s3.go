package driver

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/DennisDenuto/boshver-resource/version"
)

type Servicer interface {
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
	PutObject(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
}

type S3Driver struct {
	InitialVersion version.BoshVersion

	Svc                  Servicer
	BucketName           string
	Key                  string
	ServerSideEncryption string
}

func (driver *S3Driver) Bump(bump version.Bump) (version.BoshVersion, error) {
	var currentVersion version.BoshVersion

	resp, err := driver.Svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(driver.BucketName),
		Key:    aws.String(driver.Key),
	})
	if err == nil {
		bucketNumberPayload, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return version.BoshVersion{}, err
		}
		defer resp.Body.Close()

		payloadStr := strings.TrimSpace(string(bucketNumberPayload))
		currentVersion, err = version.Parse(payloadStr)
		if err != nil {
			return version.BoshVersion{}, err
		}
	} else if s3err, ok := err.(awserr.RequestFailure); ok && s3err.StatusCode() == 404 {
		currentVersion = driver.InitialVersion
	} else {
		return version.BoshVersion{}, err
	}

	newVersion := bump.Apply(currentVersion)

	err = driver.Set(newVersion)
	if err != nil {
		return version.BoshVersion{}, err
	}

	return newVersion, nil
}

func (driver *S3Driver) Set(newVersion version.BoshVersion) error {
	params := &s3.PutObjectInput{
		Bucket:      aws.String(driver.BucketName),
		Key:         aws.String(driver.Key),
		ContentType: aws.String("text/plain"),
		Body:        bytes.NewReader([]byte(newVersion.String())),
		ACL:         aws.String(s3.ObjectCannedACLPrivate),
	}

	if len(driver.ServerSideEncryption) > 0 {
		params.ServerSideEncryption = aws.String(driver.ServerSideEncryption)
	}

	_, err := driver.Svc.PutObject(params)
	return err
}

func (driver *S3Driver) Check(cursor *version.BoshVersion) ([]version.BoshVersion, error) {
	var bucketNumber string

	resp, err := driver.Svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(driver.BucketName),
		Key:    aws.String(driver.Key),
	})
	if err == nil {
		bucketNumberPayload, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []version.BoshVersion{}, err
		}
		defer resp.Body.Close()

		bucketNumber = string(bucketNumberPayload)
	} else if s3err, ok := err.(awserr.RequestFailure); ok && s3err.StatusCode() == 404 {
		if cursor == nil {
			return []version.BoshVersion{driver.InitialVersion}, nil
		} else {
			return []version.BoshVersion{}, nil
		}
	} else {
		return nil, err
	}

	bucketVersion, err := version.Parse(bucketNumber)
	if err != nil {
		return nil, fmt.Errorf("parsing number in bucket: %s", err)
	}

	if cursor == nil || bucketVersion.Compare(*cursor) > 0 {
		return []version.BoshVersion{bucketVersion}, nil
	}

	return []version.BoshVersion{}, nil
}
