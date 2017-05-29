package main

import (
	"io/ioutil"
	"testing"

	"github.com/op/go-logging"
	. "github.com/smartystreets/goconvey/convey"
)

func TestServer(t *testing.T) {
	log := logging.MustGetLogger("module")
	logBackend := logging.AddModuleLevel(logging.NewLogBackend(ioutil.Discard, "", 0))
	log.SetBackend(logBackend)

	s := &Server{
		Log: log,
	}
	s.SetLogger()
	Convey("Test executeScripts", t, func() {

		Convey("Good script", func() {
			location := &Location{
				Desciption:      "testCase",
				SavePath:        "/tmp",
				BashExecTimeout: 1,
				BashExec: "true \n" +
					"true \n" +
					"true \n",
			}
			err := s.executeScritps(location, "filename")
			So(err, ShouldBeNil)
		})
		Convey("Bad script", func() {
			location := &Location{
				Desciption:      "testCase",
				SavePath:        "/tmp",
				BashExecTimeout: 1,
				BashExec: "true \n" +
					"false \n" +
					"true \n",
			}
			err := s.executeScritps(location, "filename")
			So(err, ShouldNotBeNil)
		})
		Convey("Long script", func() {
			location := &Location{
				Desciption:      "testCase",
				SavePath:        "/tmp",
				BashExecTimeout: 1,
				BashExec:        "sleep 10",
			}
			err := s.executeScritps(location, "filename")
			So(err, ShouldNotBeNil)
		})
	})
}
