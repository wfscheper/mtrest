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
	"math"

	"github.com/wfscheper/mtrest"
)

type Score struct {
	Value int
	Q     float64
	Index int
}

// Returns an empty Score.
func New() *Score {
	return &Score{}
}

// Cmp returns 1 if s is greater than other, 0 if they are equal, and -1 if s is less than other.
func (s *Score) Cmp(other *Score) (d int) {
	// something is always greater than nothing
	if other == nil {
		return 1
	}
	// compare scores
	d = cmpInt(s.Value, other.Value)
	if d == 0 {
		// compare quality factors
		d = cmpFloat(s.Q, other.Q)
		if d == 0 {
			d = cmpInt(other.Index, s.Index)
		}
	}
	return d
}

// SetIndex updates the index of s and returns it.
func (s *Score) SetIndex(i int) *Score {
	s.Index = i
	return s
}

// BestMatch returns m's highest score against the MediaTypes in choices. If m does not match any MediaType in choices, BestMatch returns nil.
func BestMatch(m *mtrest.MediaType, choices []*mtrest.MediaType) *Score {
	var best *Score
	for idx, choice := range choices {
		score := Match(m, choice)
		if score != nil && score.Cmp(best) > 0 {
			best = score.SetIndex(idx)
		}
	}
	return best
}

// Match returns a Score representing how closely the MediaTYpes a and b match. If either a or b are nil, or there is no match between them, then Match returns nil.
func Match(a, b *mtrest.MediaType) (score *Score) {
	if !(a == nil || b == nil) && (a.Type == b.Type || a.Type == "*" || b.Type == "*") && (a.SubType == b.SubType || a.SubType == "*" || b.SubType == "*") {
		score = New()
		if a.Type == b.Type {
			score.Value += 100
		}
		if a.SubType == b.SubType {
			score.Value += 10
		}
		for k, v := range a.Params {
			if k != "q" && b.Params[k] == v {
				score.Value += 1
			}
		}
		score.Q = toFixed(a.Q*b.Q, 3)
	}
	return
}

func cmpInt(a, b int) int {
	d := a - b
	switch {
	case d > 0:
		return 1
	case d < 0:
		return -1
	default:
		return 0
	}
}

func cmpFloat(a, b float64) int {
	d := a - b
	switch {
	case d > 0:
		return 1
	case d < 0:
		return -1
	default:
		return 0
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
