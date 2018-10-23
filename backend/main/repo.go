package main

import (
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
)

type Repo interface {
	Upload() error
}

type repo struct {
	st  *storage.Client
	bkt *storage.BucketHandle
	obj *storage.ObjectHandle
	w   *storage.Writer
}

func NewRepo(st *storage.Client, bkt *storage.BucketHandle, obj *storage.ObjectHandle, w *storage.Writer) repo {
	return repo{

		st:  st,
		bkt: bkt,
		obj: obj,
		w:   w,
	}
}

func (r *repo) Upload() error {

	f, err := os.Open("topify-list.csv")
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = io.Copy(r.w, f); err != nil {
		return err
	}
	if err := r.w.Close(); err != nil {
		return err
	}

	log.Print("info", "file uploaded")
	return nil
}

func (r *repo) Download() error {

	// wc := r.w
	// wc.ContentType = "text/plain"
	// wc.Metadata = map[string]string{
	// 		"x-goog-meta-foo": "foo",
	// 		"x-goog-meta-bar": "bar",
	// }
	// d.cleanUp = append(d.cleanUp, fileName)

	// if _, err := wc.Write([]byte("abcde\n")); err != nil {
	// 		d.errorf("createFile: unable to write data to bucket %q, file %q: %v", d.bucketName, fileName, err)
	// 		return
	// }

	// if err := wc.Close(); err != nil {
	// 		d.errorf("createFile: unable to close bucket %q, file %q: %v", d.bucketName, fileName, err)
	// 		return
	// }

	return nil
}
