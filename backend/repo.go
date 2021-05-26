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

	f, err := os.Open("/tmp/topify-list.csv")
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

	log.Print("updated CSV file uploaded")
	return nil
}
