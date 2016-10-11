// Copyright 2014 The Macaron Authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package i18n

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/insionng/vodka"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/macaron.v1"
)

func Test_Version(t *testing.T) {
	Convey("Check package version", t, func() {
		So(Version(), ShouldEqual, _VERSION)
	})
}

func Test_I18n(t *testing.T) {
	Convey("Use i18n middleware", t, func() {
		Convey("No langauge", func() {
			defer func() {
				So(recover(), ShouldNotBeNil)
			}()

			v := vodka.New()
			v.Use(I18n(Options{}))
		})

		Convey("Languages and names not match", func() {
			defer func() {
				So(recover(), ShouldNotBeNil)
			}()

			v := vodka.New()
			v.Use(I18n(Options{
				Langs: []string{"en-US"},
			}))
		})

		Convey("Invalid directory", func() {
			defer func() {
				So(recover(), ShouldNotBeNil)
			}()

			v := vodka.New()
			v.Use(I18n(Options{
				Directory: "404",
				Langs:     []string{"en-US"},
				Names:     []string{"English"},
			}))
		})

		Convey("With correct options", func() {
			v := vodka.New()
			v.Use(I18n(Options{
				Files: map[string][]byte{"locale_en-US.ini": []byte("")},
				Langs: []string{"en-US"},
				Names: []string{"English"},
			}))
			v.Get("/", func(self vodka.Context) error { return nil })

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/", nil)
			So(err, ShouldBeNil)
			m := macaron.New()
			m.ServeHTTP(resp, req)
		})

		Convey("Set by redirect of URL parameter", func() {
			v := vodka.New()
			v.Use(I18n(Options{
				Langs:    []string{"en-US"},
				Names:    []string{"English"},
				Redirect: true,
			}))
			v.Get("/", func(self vodka.Context) error { return self.String(http.StatusOK, "") })

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/?lang=en-US", nil)
			So(err, ShouldBeNil)
			req.RequestURI = "/?lang=en-US"
			m := macaron.New()
			m.ServeHTTP(resp, req)
		})

		Convey("Set by Accept-Language", func() {
			v := vodka.New()
			v.Use(I18n(Options{
				Langs: []string{"en-US", "zh-CN"},
				Names: []string{"English", "简体中文"},
			}))
			v.Get("/", func(self vodka.Context) error {
				l := self.Get("i18n")
				So(l.(Locale).Language(), ShouldEqual, "en-US")
				return nil
			})

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/", nil)
			So(err, ShouldBeNil)
			req.Header.Set("Accept-Language", "en-US")
			m := macaron.New()
			m.ServeHTTP(resp, req)
		})
	})
}
