package bla

import "golang.org/x/net/webdav"

func loadWebDav(s *Handler) {

	fs := webdav.Dir(s.Cfg.RootPath)
	ls := webdav.NewMemLS()

	handler := &webdav.Handler{
		Prefix:     "/fs",
		FileSystem: fs,
		LockSystem: ls,
	}
	a := NewAuthRateByIPHandler(s.Cfg.HostName, handler, s.Cfg.UserName,
		s.Cfg.Password, 3)
	s.webfs = a
}
