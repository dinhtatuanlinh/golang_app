package models

import (
	"cloud.google.com/go/storage"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type File struct {
	File string `json:"file"`
}

type ClientUploader struct {
	Cl         *storage.Client
	ProjectID  string
	BucketName string
	UploadPath string
}

var uploader *ClientUploader

func (c *ClientUploader) UploadFile() error {
	// get project base folder
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	credential := basepath + "/test-files/dev-ab-kyc-vcall.json"
	fmt.Println(credential)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	//os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "localhost:3030")
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credential))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	f, err := os.Open("test-files/202111300752049781-42a207ab-52b1-4367-bbd0-d6d39bb9720d.jpeg")
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	ctx, cancel = context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	o := client.Bucket("dev-ab-kyc-vcall").Object("linh")
	o = o.If(storage.Conditions{DoesNotExist: true})

	// Upload an object with storage.Writer.
	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}
