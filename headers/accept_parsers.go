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
	"strings"

	"github.com/wfscheper/mtrest"
	"github.com/wfscheper/mtrest/internal/fitness"
)

// Accepts is a set of media types accepted by a client.
type Accepts []*mtrest.MediaType

// NewAccepts returns a Accepts list constructed from s, a comman-separated list of media types.
func NewAccepts(s string) (Accepts, error) {
	var accepts Accepts

	parts := strings.Split(s, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		mt, err := mtrest.NewMediaType(part)
		if err != nil {
			return nil, err
		}
		accepts = append(accepts, mt)
	}
	return accepts, nil
}

// BestMatch returns the MediaType in offers that best matches the MediaTypes in an Accepts header.
func (a Accepts) BestMatch(offers []*mtrest.MediaType) (m *mtrest.MediaType) {
	var best *fitness.Score
	for _, accept := range a {
		score := fitness.BestMatch(accept, offers)
		if score != nil && score.Cmp(best) > 0 {
			best = score
		}
	}
	if best != nil {
		m = offers[best.Index]
	}
	return
}
