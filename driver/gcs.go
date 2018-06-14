package driver

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"os"

	"cloud.google.com/go/storage"
	"github.com/DennisDenuto/boshver-resource/version"
	"google.golang.org/api/option"
)

type GCSDriver struct {


	InitialVersion version.BoshVersion

	Servicer   IOServicer
	BucketName string
	Key        string
}

func (d *GCSDriver) Bump(b version.Bump) (version.BoshVersion, error) {
	versions, err := d.Check(nil)

	if err != nil {
		return version.BoshVersion{}, err
	}

	if len(versions) == 0 {
		return version.BoshVersion{}, nil
	}

	newVersion := b.Apply(versions[0])
	err = d.Set(newVersion)

	if err != nil {
		return version.BoshVersion{}, err
	}
	return newVersion, nil
}

func (d *GCSDriver) Set(v version.BoshVersion) error {
	w, err := d.Servicer.PutObject(d.BucketName, d.Key)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = w.Write([]byte(v.String()))
	return err
}

func (d *GCSDriver) Check(cursor *version.BoshVersion) ([]version.BoshVersion, error) {
	r, err := d.Servicer.GetObject(d.BucketName, d.Key)

	switch err {
	case storage.ErrObjectNotExist:
		if cursor == nil {
			return []version.BoshVersion{d.InitialVersion}, nil
		}
		return []version.BoshVersion{}, nil
	case nil:
	default:
		return nil, err
	}
	defer r.Close()

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	v, err := version.Parse(string(b))
	if err != nil {
		return nil, fmt.Errorf("parsing number in bucket: %s", err)
	}

	if cursor == nil || v.Compare(*cursor) > 0 {
		return []version.BoshVersion{v}, nil
	}

	return nil, nil
}

type IOServicer interface {
	GetObject(bucketName, objectName string) (io.ReadCloser, error)
	PutObject(bucketName, objectName string) (io.WriteCloser, error)
}

type GCSIOServicer struct {
	JSONCredentials string
}

func (s *GCSIOServicer) GetObject(bucketName, objectName string) (io.ReadCloser, error) {
	temp, err := ioutil.TempFile("", "auth-credentials.json")
	if err != nil {
		return nil, err
	}

	_, err = temp.WriteString(s.JSONCredentials)
	if err != nil {
		return nil, err
	}
	defer os.Remove(temp.Name())
	ctx := context.Background()

	authOption := option.WithCredentialsFile(temp.Name())
	client, err := storage.NewClient(ctx, authOption)

	if err != nil {
		return nil, err
	}

	bkt := client.Bucket(bucketName)
	obj := bkt.Object(objectName)

	return obj.NewReader(context.Background())
}

func (s *GCSIOServicer) PutObject(bucketName, objectName string) (io.WriteCloser, error) {
	temp, err := ioutil.TempFile("", "auth-credentials.json")
	if err != nil {
		return nil, err
	}

	_, err = temp.WriteString(s.JSONCredentials)
	if err != nil {
		return nil, err
	}
	defer os.Remove(temp.Name())
	ctx := context.Background()

	authOption := option.WithCredentialsFile(temp.Name())
	client, err := storage.NewClient(ctx, authOption)

	if err != nil {
		return nil, err
	}

	bkt := client.Bucket(bucketName)
	obj := bkt.Object(objectName)

	return obj.NewWriter(context.Background()), nil
}
