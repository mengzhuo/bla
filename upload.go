package bla

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type TinyMceUploadResponse struct {
	Location string `json:"location"`
}

func (s *Server) UploadMedia(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(Cfg.MaxUploadSize)
	f, h, err := r.FormFile("name")
	if err != nil {
		log.Print(err)
		w.Write([]byte("500 internal error"))
		return
	}

	extAllow := false
	for _, ext := range []string{".jpg", ".png", ".gif", ".jpeg"} {
		if filepath.Ext(h.Filename) == ext {
			extAllow = true
			break
		}
	}

	if !extAllow {
		w.Write([]byte("403 extend not allow"))
		return
	}
	now := time.Now()
	uploadPath := filepath.Join(Cfg.UploadPath, now.Format("2006-01"))
	os.MkdirAll(uploadPath, 0755)

	rf, err := os.Create(filepath.Join(uploadPath, santiSpace(h.Filename)))
	if err != nil {
		log.Print(err)
		return
	}
	_, err = io.Copy(rf, f)
	if err == nil {
		rsp := &TinyMceUploadResponse{filepath.Join(Cfg.BasePath, filepath.Base(Cfg.UploadPath),
			now.Format("2006-01"), h.Filename)}
		enc := json.NewEncoder(w)
		enc.Encode(rsp)
	} else {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
