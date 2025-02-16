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
	expectedFileName := "717d5b5fe0e1a5f88216179791868b6ca042c155"
	expectedPathName := "717d5/b5fe0/e1a5f/88216/17979/1868b/6ca04/2c155"

	if pathKey.PathName != expectedPathName {
		t.Errorf("expected %s, got %s", expectedPathName, pathKey.PathName)
	}

	if pathKey.Filename != expectedFileName {
		t.Errorf("expected %s, got %s", expectedFileName, pathKey.Filename)
	}
}

func TestStoreDelete(t *testing.T) {
	s := newStoreHelper()
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
	s := newStoreHelper()
	defer teardown(t, s)

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key-%d", i)
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

		if err := s.Delete(key); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); ok {
			t.Errorf("expected %s to not exist", key)
		}

	}
}

func newStoreHelper() *Store {
	cfg := StoreConfig{
		PathTransformFunc: CASPathTransformFunc,
	}
	return NewStore(cfg)
}

func teardown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}
