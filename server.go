package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/op/go-logging"
	"gopkg.in/tomb.v2"
)

type Server struct {
	config *Config
	log    *logging.Logger
	tomb   tomb.Tomb
}

func (s *Server) Start() error {
	log = s.log

	server := &http.Server{
		Addr:    s.config.Listen,
		Handler: http.HandlerFunc(s.handler),
	}

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return err
	}
	s.tomb.Go(func() error {
		err := server.Serve(listener)
		select {
		case <-s.tomb.Dying():
			return nil
		default:
			return err
		}
	})

	s.tomb.Go(func() error {
		<-s.tomb.Dying()
		return listener.Close()
	})

	// http.HandleFunc("/", s.handler)
	// if err := http.ListenAndServe(s.config.Listen, nil); err != nil {
	// 	return err
	// }

	return nil
}

func (s *Server) Stop() error {
	s.tomb.Kill(nil)
	return s.tomb.Wait()
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("method [%s] uri [%s] ", r.Method, r.RequestURI)
	if r.Method != "POST" {
		responceError(w, fmt.Errorf("Bad method"))
		return
	}
	var (
		l        *Location
		fileName string
		err      error
	)

	if l, err = s.getLocation(r); err != nil {
		responceError(w, err)
		log.Errorf("FormParse %v\n", err)
		return
	}
	if fileName, err = s.save(l, r); err != nil {
		responceError(w, err)
		log.Errorf("SaveFile %v\n", err)
		return
	}
	log.Debugf("Successfully save file [%s]", fileName)

	if err = s.executeScritps(l, fileName); err != nil {
		responceError(w, fmt.Errorf("Error in after upload script"))
		log.Errorf("ExecuteScript %v\n", err)
		return
	}
	responceOK(w, fileName)
}

func (s *Server) executeScritps(location *Location, fileName string) error {
	if len(location.BashExec) < 1 {
		log.Debugf("No scripts for [%s]", fileName)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(location.BashExecTimeout)*time.Second)
	defer cancel()

	script := strings.Replace(location.BashExec, "%filename%", fileName, -1)
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("location [%s] script status [%v] output [%s]", location.SavePath, err, string(out))
	}
	log.Debugf("Successfully execute script on file [%s]\n%s", fileName, string(out))
	return nil
}

func (s *Server) getLocation(r *http.Request) (*Location, error) {
	var (
		location *Location
	)
	for _, location = range s.config.Locations {
		if r.RequestURI == location.URI {
			return location, nil
		}
	}
	return nil, fmt.Errorf("Bad Location [%s]", r.RequestURI)
}

func (s *Server) save(l *Location, r *http.Request) (string, error) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		return "", err
	}
	defer file.Close()

	if !l.fileNameRe.MatchString(handler.Filename) {
		return "", fmt.Errorf("Bad filename [%s] re [%s]", handler.Filename, l.FileNameRegexp)
	}

	filePath := filepath.Join(l.SavePath, handler.Filename)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()
	io.Copy(f, file)

	return handler.Filename, nil
}

func responceError(w http.ResponseWriter, err error) {
	body := `
<html>
<head><title>400 Bad Request</title></head>
<body bgcolor="white">
<center><h1>400 Bad Request</h1></center>
<center>%s</center>
<hr><center>candy-upload</center>
</body>
</html>
`
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, body, err.Error())
}

func responceOK(w http.ResponseWriter, s string) {
	body := `
<html>
<head><title>201 Created</title></head>
<body bgcolor="white">
<center><h1>201 Created</h1></center>
<center>%s</center>
<hr><center>candy-upload</center>
</body>
</html>
`
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, body, s)
}
