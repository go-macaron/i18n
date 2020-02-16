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
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/macaron.v1"
)

func TestI18n(t *testing.T) {
	t.Run("no language", func(t *testing.T) {
		defer func() {
			assert.Equal(t, "no language is specified", recover())
		}()

		m := macaron.New()
		m.Use(I18n(Options{}))
	})

	t.Run("languages and names not match", func(t *testing.T) {
		defer func() {
			assert.Equal(t, "length of langs is not same as length of names", recover())
		}()

		m := macaron.New()
		m.Use(I18n(Options{
			Langs: []string{"en-US"},
		}))
	})

	t.Run("invalid directory", func(t *testing.T) {
		defer func() {
			assert.Equal(t, errors.New("fail to set message file(en-US): open 404/locale_en-US.ini: no such file or directory"), recover())
		}()

		m := macaron.New()
		m.Use(I18n(Options{
			Directory: "404",
			Langs:     []string{"en-US"},
			Names:     []string{"English"},
		}))
	})

	t.Run("with correct options", func(t *testing.T) {
		m := macaron.New()
		m.Use(I18n(Options{
			Files: map[string][]byte{"locale_en-US.ini": []byte("")},
			Langs: []string{"en-US"},
			Names: []string{"English"},
		}))
		m.Get("/", func() {})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		m.ServeHTTP(resp, req)
	})

	t.Run("set by Accept-Language", func(t *testing.T) {
		m := macaron.New()
		m.Use(I18n(Options{
			Langs: []string{"en-US", "zh-CN", "it-IT"},
			Names: []string{"English", "简体中文", "Italiano"},
		}))
		m.Get("/", func(l Locale) {
			assert.Equal(t, "it-IT", l.Language())
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Accept-Language", "it")
		m.ServeHTTP(resp, req)
	})

	t.Run("set to default language", func(t *testing.T) {
		m := macaron.New()
		m.Use(I18n(Options{
			Langs: []string{"en-US", "zh-CN", "it-IT"},
			Names: []string{"English", "简体中文", "Italiano"},
		}))
		m.Get("/", func(l Locale) {
			assert.Equal(t, "en-US", l.Language())
		})

		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Accept-Language", "ru")
		m.ServeHTTP(resp, req)
	})
}

func TestRedirect(t *testing.T) {
	m := macaron.New()
	m.Use(I18n(Options{
		Langs:    []string{"en-US"},
		Names:    []string{"English"},
		Redirect: true,
	}))
	m.Get("/", func() {})

	tests := []struct {
		url    string
		expURL string
	}{
		{
			url:    "/?lang=en-US",
			expURL: "/",
		}, {
			url:    "//example.com?lang=en-US",
			expURL: "/example.com",
		}, {
			url:    "/abc/../../../example.com?lang=en-US",
			expURL: "/example.com",
		}, {
			url:    "/../abc/../example.com?lang=en-US",
			expURL: "/example.com",
		},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", test.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.RequestURI = test.url
			m.ServeHTTP(resp, req)

			assert.Equal(t, http.StatusFound, resp.Code)
			assert.Equal(t, "<a href=\""+test.expURL+"\">Found</a>.\n\n", resp.Body.String())
		})
	}
}
