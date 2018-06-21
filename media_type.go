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
	"fmt"
	"mime"
	"strconv"
	"strings"
)

type MediaType struct {
	Type     string
	SubType  string
	Params   map[string]string
	Q        float64
	Unparsed string
}

func NewMediaType(s string) (*MediaType, error) {
	mt, p, err := mime.ParseMediaType(s)
	if err != nil {
		return nil, err
	}
	m := &MediaType{Params: p, Q: 1.0, Unparsed: s}
	i := strings.Index(mt, "/")
	if i == -1 {
		m.Type = mt
	} else {
		m.Type, m.SubType = mt[:i], mt[i+1:]
	}
	if v, ok := p["q"]; ok {
		m.Q, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("Error parsing quality factor: '%s'", v)
		}
	}
	return m, nil
}

// Encoding returns the media type's encoding, or 'json' if no encoding is
// specified. The encoding is specified by adding +<encoding> to the subtype.
func (m MediaType) Encoding() string {
	if i := strings.LastIndex(m.SubType, "+"); i >= 0 {
		return m.SubType[i+1:]
	}
	return m.SubType
}

func (m MediaType) String() string {
	return mime.FormatMediaType(m.Type+"/"+m.SubType, m.Params)
}
