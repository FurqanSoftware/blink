package fs

import (
	"bytes"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"time"
)

type FileServer struct {
	fs           fs.FS
	tryIndexHTML bool
}

func (s FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path

	f, err := s.fs.Open(name[1:])
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}

	if d.IsDir() && s.tryIndexHTML {
		// Try index.html
		index := strings.TrimSuffix(name, "/") + "/index.html"
		ff, err := s.fs.Open(index)
		if err == nil {
			defer ff.Close()
			dd, err := ff.Stat()
			if err == nil {
				d = dd
				f = ff
			}
		}
	}

	if d.IsDir() {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	b := bytes.Buffer{}
	_, err = io.Copy(&b, f)
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}
	http.ServeContent(w, r, d.Name(), time.Time{}, bytes.NewReader(b.Bytes()))
}
