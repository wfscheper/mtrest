// Copyright © 2017 Walter Scheper <walter.scheper@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mtrest

import (
	"testing"
)

func assertMediaTypeEqual(a, b *MediaType) bool {
	if a.Type != b.Type {
		return false
	}
	if a.SubType != b.SubType {
		return false
	}
	if a.Q != b.Q {
		return false
	}
	if len(a.Params) != len(b.Params) {
		return false
	}
	for key, value := range a.Params {
		if b.Params[key] != value {
			return false
		}
	}
	return true
}

func TestEncoding(t *testing.T) {
	tests := []struct {
		mt       string
		expected string
	}{
		{"application/json", "json"},
		{"application/yaml", "yaml"},
		{"application/xml", "xml"},
		{"application/vnd.foo+json", "json"},
		{"application/vnd.foo+yaml", "yaml"},
		{"application/vnd.foo+xml", "xml"},
	}
	for idx, test := range tests {
		m, _ := NewMediaType(test.mt)
		a := m.Encoding()
		if a != test.expected {
			t.Errorf("%d: Wanted %s, got %s", idx, test.expected, a)
		}
	}
}

func TestNewMediaType(t *testing.T) {
	tests := []struct {
		in       string
		expected *MediaType
	}{
		{"text", &MediaType{
			Type:     "text",
			SubType:  "",
			Params:   map[string]string{},
			Q:        1.0,
			Unparsed: "text",
		}},
		{"text/plain", &MediaType{
			Type:     "text",
			SubType:  "plain",
			Params:   map[string]string{},
			Q:        1.0,
			Unparsed: "text/plain",
		}},
		{"text/*", &MediaType{
			Type:     "text",
			SubType:  "*",
			Params:   map[string]string{},
			Q:        1.0,
			Unparsed: "text/*",
		}},
		{"*/*", &MediaType{
			Type:     "*",
			SubType:  "*",
			Params:   map[string]string{},
			Q:        1.0,
			Unparsed: "*/*",
		}},
		{"text/plain; version=1", &MediaType{
			Type:     "text",
			SubType:  "plain",
			Params:   map[string]string{"version": "1"},
			Q:        1.0,
			Unparsed: "text/plain; version=1",
		}},
		{"text/plain; q=0.3", &MediaType{
			Type:     "text",
			SubType:  "plain",
			Params:   map[string]string{"q": "0.3"},
			Q:        0.3,
			Unparsed: "text/plain; q=0.3",
		}},
	}
	for idx, test := range tests {
		mt, err := NewMediaType(test.in)
		if err != nil {
			t.Errorf("%d: exected nil, got %q", idx, err)
		}
		if !assertMediaTypeEqual(test.expected, mt) {
			t.Errorf("%d: Expected %q, got %q", idx, test.expected, mt)
		}
	}
}

func TestNewMediaTypeErrors(t *testing.T) {
	tests := []struct {
		in       string
		expected string
	}{
		{"/", "mime: no media type"},
		{"text/", "mime: expected token after slash"},
		{"text/plain/a", "mime: unexpected content after media subtype"},
		{"text/plain; q=foo", "Error parsing quality factor: 'foo'"},
	}
	for idx, test := range tests {
		_, err := NewMediaType(test.in)
		if err == nil {
			t.Errorf("%d: expected '%s', got nil", idx, test.expected)
		} else if err.Error() != test.expected {
			t.Errorf("%d: expected '%s', got '%q'", idx, test.expected, err)
		}
	}
}

func TestString(t *testing.T) {
	tests := []string{
		"a/b",
		"a/b; p=1",
		"a/b+c",
		"a/b+c; p=1",
		"a/b; q=0.5",
		"a/b; p=1; q=0.5",
		"a/b+c; q=0.5",
		"a/b+c; p=1; q=0.5",
	}
	for i, test := range tests {
		m, _ := NewMediaType(test)
		if actual := m.String(); actual != test {
			t.Errorf("%d: expected '%s', got '%s'", i, test, actual)
		}
	}
}

func BenchmarkNewMediaType(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewMediaType("text/plain; q=0.8; version=1")
	}
}
