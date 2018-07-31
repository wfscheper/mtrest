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

package headers

import (
	"reflect"
	"testing"

	"github.com/wfscheper/mtrest"
)

func TestNewAccepts(t *testing.T) {
	m, err := mtrest.NewMediaType("a/b")
	if err != nil {
		t.Fatal(err)
	}
	n, err := mtrest.NewMediaType("c/d")
	if err != nil {
		t.Fatal(err)
	}

	expected := make(Accepts, 0)
	expected = append(expected, m)
	if actual, _ := NewAccepts("a/b"); !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %+v, got %+v", expected, actual)
	}

	expected = append(expected, n)
	inputs := []string{"a/b,c/d", " a/b,c/d", "a/b ,c/d", "a/b, c/d", "a/b,c/d "}
	for i, input := range inputs {
		if actual, _ := NewAccepts(input); !reflect.DeepEqual(expected, actual) {
			t.Fatalf("%d: expected %+v, got %+v", i, expected, actual)
		}
	}
}

func TestNewAcceptsErrors(t *testing.T) {
	_, err := NewAccepts("a/")
	if err == nil || err.Error() != "mime: expected token after slash" {
		t.Fatalf("expected %s, got %+v", "mime: expected token after slash", err)
	}
	_, err = NewAccepts("/a")
	if err == nil || err.Error() != "mime: no media type" {
		t.Fatalf("expected %s, got %+v", "mime: no media type", err)
	}
}

func TestBestMatch(t *testing.T) {
	textPlain, _ := mtrest.NewMediaType("text/plain")
	offers := []*mtrest.MediaType{
		&mtrest.ApplicationJson,
		&mtrest.ApplicationYaml,
		textPlain,
	}
	applicationYamlQS, _ := mtrest.NewMediaType("application/yaml; q=0.4")
	applicationJsonQS, _ := mtrest.NewMediaType("application/json; q=0.8")
	textPlainQS, _ := mtrest.NewMediaType("text/plain; q=0.2")
	qsOffers := []*mtrest.MediaType{
		applicationYamlQS,
		applicationJsonQS,
		textPlainQS,
	}
	tests := []struct {
		accepts  string
		offers   []*mtrest.MediaType
		expected string
	}{
		{"application/json,application/yaml", offers, "application/json"},
		{"application/json,application/yaml", qsOffers, "application/json; q=0.8"},
		{"application/json; q=0.001, application/yaml", offers, "application/yaml"},
		{"application/json; q=0.001, application/yaml", qsOffers, "application/yaml; q=0.4"},
		{"*/*, application/yaml", offers, "application/yaml"},
		{"*/*, application/yaml", qsOffers, "application/yaml; q=0.4"},
		{"application/*, application/yaml", offers, "application/yaml"},
		{"application/*, application/yaml", qsOffers, "application/yaml; q=0.4"},
		{"*/*, text/html", offers, "application/json"},
		{"*/*, text/html", qsOffers, "application/json; q=0.8"},
		{"application/*, text/html", offers, "application/json"},
		{"application/*, text/html", qsOffers, "application/json; q=0.8"},
	}
	for i, test := range tests {
		accepts, _ := NewAccepts(test.accepts)
		actual := accepts.BestMatch(test.offers)
		if test.expected != actual.String() {
			t.Errorf("%d: expected %+v, got %+v", i, test.expected, actual.String())
		}
	}
}
