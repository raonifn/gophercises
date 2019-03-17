package urlshort

import (
	"fmt"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
)

type dbMapper struct {
	db       *bolt.DB
	fallback http.Handler
}

func runFind(path string, b *bolt.Bucket) string {
	url := b.Get([]byte(path))
	return string(url)
}

func (dbm *dbMapper) find(path string) (string, error) {
	var url string
	err := dbm.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("redirect"))
		url = runFind(path, b)
		return nil
	})
	return url, err
}

func (dbm *dbMapper) handle(w http.ResponseWriter, r *http.Request) {
	url, err := dbm.find(r.URL.Path)
	if err != nil {
		msg := fmt.Sprintf("%v", err)
		w.Write([]byte(msg))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	if url == "" {
		dbm.fallback.ServeHTTP(w, r)
		return
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func DBHandler(file string) (HandlerStacker, error) {
	db, err := bolt.Open(file, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("redirect"))
		return err
	})
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("redirect"))
		return b.Put([]byte("/db"), []byte("https://github.com/boltdb/bolt"))
	})

	dbMapper := dbMapper{db: db}

	return func(fallback http.Handler) http.Handler {
		dbMapper.fallback = fallback
		return http.HandlerFunc(dbMapper.handle)
	}, nil
}
