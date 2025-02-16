package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "myjpg"
	pathKey := CASPathTransformFunc(key)
	expectedOriginalKey := "717d5b5fe0e1a5f88216179791868b6ca042c155"
	expectedPathName := "717d5/b5fe0/e1a5f/88216/17979/1868b/6ca04/2c155"

	if pathKey.PathName != expectedPathName {
		t.Errorf("expected %s, got %s", expectedPathName, pathKey.PathName)
	}

	if pathKey.Filename != expectedOriginalKey {
		t.Errorf("expected %s, got %s", expectedOriginalKey, pathKey.Filename)
	}
}

func TestStoreDelete(t *testing.T) {
	cfg := StoreConfig{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(cfg)
	key := "special"

	data := []byte("some data")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

func TestStore(t *testing.T) {
	cfg := StoreConfig{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(cfg)
	key := "special"

	data := []byte("some data")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if ok := s.Has(key); !ok {
		t.Errorf("expected %s to exist", key)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}
	b, err := io.ReadAll(r)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(b, data) {
		t.Errorf("expected %s, got %s", string(data), string(b))
	}
	fmt.Println("SUCCESS READ: " + string(b))

	if err = s.Delete(key); err != nil {
		t.Error(err)
	}
}
