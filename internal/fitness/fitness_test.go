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

package fitness

import (
	"testing"

	"github.com/wfscheper/mtrest"
)

func TestCmp(t *testing.T) {
	tests := []struct {
		title    string
		a, b     *Score
		expected int
	}{
		{"Zero-value fitness score cmp", &Score{}, &Score{}, 0},
		{"Score 1 vs zero-value fitness score", &Score{1, 0.0, 0}, &Score{}, 1},
		{"Zero-value vs score 1", &Score{}, &Score{1, 0.0, 0}, -1},
		{"Fitness score 1 vs 1", &Score{1, 0.0, 0}, &Score{1, 0.0, 0}, 0},
		{"Fitness score 110 vs 100", &Score{110, 0.0, 0}, &Score{100, 0.0, 0}, 1},
		{"Fitness score 100 vs 110", &Score{100, 0.0, 0}, &Score{110, 0.0, 0}, -1},
		{"Quality 1.0 vs 0.0", &Score{0, 1.0, 0}, &Score{0, 0.0, 0}, 1},
		{"Quality 0.0 vs 1.0", &Score{0, 0.0, 0}, &Score{0, 1.0, 0}, -1},
		{"Quality 0.8 vs 0.4", &Score{0, 0.8, 0}, &Score{0, 0.4, 0}, 1},
		{"Quality 0.4 vs 0.8", &Score{0, 0.4, 0}, &Score{0, 0.8, 0}, -1},
		{"Index 1 vs 0", &Score{0, 0.0, 1}, &Score{0, 0.0, 0}, -1},
		{"Index 0 vs 1", &Score{0, 0.0, 0}, &Score{0, 0.0, 1}, 1},
		{"Index 10 vs 0", &Score{0, 0.0, 10}, &Score{0, 0.0, 0}, -1},
		{"Index 0 vs 10", &Score{0, 0.0, 0}, &Score{0, 0.0, 10}, 1},
	}
	for idx, test := range tests {
		actual := test.a.Cmp(test.b)
		if actual != test.expected {
			t.Errorf("%d: (%s) Expected %d, got %d", idx, test.title, test.expected, actual)
		}
	}
}

func Test_cmpInt(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{0, 0, 0},
		{1, 1, 0},
		{100, 100, 0},
		{-1, -1, 0},
		{-100, -100, 0},
		{0, 1, -1},
		{0, 10, -1},
		{1, 2, -1},
		{10, 100, -1},
		{1, 0, 1},
		{10, 0, 1},
		{2, 1, 1},
		{100, 10, 1},
	}
	for idx, test := range tests {
		actual := cmpInt(test.a, test.b)
		if actual != test.expected {
			t.Errorf("%d: expected %d, got %d", idx, test.expected, actual)
		}
	}
}

func Test_cmpFloat(t *testing.T) {
	tests := []struct {
		a, b     float64
		expected int
	}{
		{0.0, 0.0, 0},
		{1.0, 1.0, 0},
		{100.0, 100.0, 0},
		{-1.0, -1.0, 0},
		{-100.0, -100.0, 0},
		{0.0, 1.0, -1},
		{0.0, 10.0, -1},
		{1.0, 2.0, -1},
		{10.0, 100.0, -1},
		{1.0, 0.0, 1},
		{10.0, 0.0, 1},
		{2.0, 1.0, 1},
		{100.0, 10.0, 1},
	}
	for idx, test := range tests {
		actual := cmpFloat(test.a, test.b)
		if actual != test.expected {
			t.Errorf("%d: expected %d, got %d", idx, test.expected, actual)
		}
	}
}

func TestMatch(t *testing.T) {
	tests := []struct {
		title, a, b string
		expected    *Score
	}{
		{"text/plain does not match audio/basic", "text/plain", "audio/basic", nil},
		{"text/plain does not match text/html", "text/plain", "text/html", nil},
		{"Non-match ignores quality factor in a", "text/plain;q=0.8", "audio/basic", nil},
		{"Non-match ignores quality factor in b", "text/plain", "audio/basic;q=0.8", nil},
		{"Non-match ignores quality factor in a and in b", "text/plain;q=0.8", "audio/basic;q=0.8", nil},
		{"text/plain;version=8 does not match audio/basic;version=8", "text/plain;version=8", "audio/basic;version=8", nil},
		{"text/plain does not match audio/*", "text/plain", "audio/*", nil},
		{"text/* does not match audio/basic", "text/*", "audio/basic", nil},
		{"text/* does not match audio/*", "text/*", "audio/*", nil},
		{"Matching type and subtype is worth 110", "text/plain", "text/plain", &Score{110, 1.0, 0}},
		{"Match honors a quality factor", "text/plain;q=0.8", "text/plain", &Score{110, 0.8, 0}},
		{"Match honors b quality factor", "text/plain", "text/plain;q=0.8", &Score{110, 0.8, 0}},
		{"Match is product of both quality factors", "text/plain;q=0.8", "text/plain;q=0.2", &Score{110, 0.16, 0}},
		{"Quality factor equality is ignored in scoring", "text/plain;q=0.8", "text/plain;q=0.8", &Score{110, 0.64, 0}},
		{"Wildcard subtype match is worth 100", "text/*", "text/plain", &Score{100, 1.0, 0}},
		{"Wildcard match is worth 0", "*/*", "text/plain", &Score{0, 1.0, 0}},
		{"Wildcard subtype in b", "text/plain", "text/*", &Score{100, 1.0, 0}},
		{"Wildcard match in b", "text/plain", "*/*", &Score{0, 1.0, 0}},
		{"Parameter in a mismatch", "text/plain;version=1", "text/plain", &Score{110, 1.0, 0}},
		{"Parameter in b mismatch", "text/plain", "text/plain;version=1", &Score{110, 1.0, 0}},
		{"Paramater match is worth one", "text/plain;version=1", "text/plain;version=1", &Score{111, 1.0, 0}},
		{"Paramater match is worth one, subtype wildcard", "text/*;version=1", "text/plain;version=1", &Score{101, 1.0, 0}},
		{"Paramater match is worth one, wildcard", "*/*;version=1", "text/plain;version=1", &Score{1, 1.0, 0}},
	}
	for idx, test := range tests {
		a, err := mtrest.NewMediaType(test.a)
		if err != nil {
			t.Errorf("%d: %q", idx, err)
		}
		b, err := mtrest.NewMediaType(test.b)
		if err != nil {
			t.Errorf("%d: %q", idx, err)
		}
		actual := Match(a, b)
		if !assertScoreEqual(test.expected, actual) {
			t.Errorf("%d: (%s) expected %+v, got %+v", idx, test.title, test.expected, actual)
		}
	}
}

func TestBestMatch(t *testing.T) {
	basicChoices := []string{"application/json", "application/yaml", "application/xml", "text/plain"}
	qsChoices := []string{"application/json", "application/yaml; q=0.8", "application/xml; q=0.5", "text/plain; q=0.1"}
	tests := []struct {
		title, m string
		choices  []string
		expected *Score
	}{
		{"Type and subtype match", "text/plain", basicChoices, &Score{110, 1.0, 3}},
		{"Type match and subtype wildcard", "text/*", basicChoices, &Score{100, 1.0, 3}},
		{"*/* picks first match", "*/*", basicChoices, &Score{0, 1.0, 0}},
		{"Full match gets product of q factors", "text/plain; q=0.5", qsChoices, &Score{110, 0.05, 3}},
		{"Subtype wildcard gets product of q factors", "text/*; q=0.5", qsChoices, &Score{100, 0.05, 3}},
		{"*/* gets product of q factors", "*/*; q=0.5", qsChoices, &Score{0, 0.5, 0}},
		{"Wildcard subtype match in choices loses to exact match", "text/plain", []string{"text/*", "text/plain"}, &Score{110, 1.0, 1}},
		{"Wildcard subtype match in choices loses to exact match with quality factor", "text/plain", []string{"text/*", "text/plain;q=0.5"}, &Score{110, 0.5, 1}},
		{"Paramater scores higher than non-paramter match", "text/plain;version=1", []string{"text/plain;version=2", "text/plain;version=1;q=0.4"}, &Score{111, 0.4, 1}},
	}
	for idx, test := range tests {
		var choices = make([]*mtrest.MediaType, len(test.choices))
		for i, choice := range test.choices {
			c, err := mtrest.NewMediaType(choice)
			if err != nil {
				t.Fatalf("%d: (%s) %q", idx, test.title, err)
			}
			choices[i] = c
		}
		m, err := mtrest.NewMediaType(test.m)
		if err != nil {
			t.Fatalf("%d: (%s) %q", idx, test.title, err)
		}
		actual := BestMatch(m, choices)
		if !assertScoreEqual(actual, test.expected) {
			t.Errorf("%d: (%s) expected %+v, got %+v", idx, test.title, test.expected, actual)
		}
	}
}

func assertScoreEqual(a, b *Score) bool {
	if a != nil {
		return (a.Cmp(b) == 0)
	}
	return b == nil
}
