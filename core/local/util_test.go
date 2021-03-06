/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package local

import (
	"testing"
	"time"

	"github.com/loadimpact/k6/lib"
	"github.com/stretchr/testify/assert"
	null "gopkg.in/guregu/null.v3"
)

func TestSumStages(t *testing.T) {
	testdata := map[string]struct {
		Time   lib.NullDuration
		Stages []lib.Stage
	}{
		"Blank":    {lib.NullDuration{}, []lib.Stage{}},
		"Infinite": {lib.NullDuration{}, []lib.Stage{{}}},
		"Limit": {
			lib.NullDurationFrom(10 * time.Second),
			[]lib.Stage{
				{Duration: lib.NullDurationFrom(5 * time.Second)},
				{Duration: lib.NullDurationFrom(5 * time.Second)},
			},
		},
		"InfiniteTail": {
			lib.NullDuration{Duration: lib.Duration(10 * time.Second), Valid: false},
			[]lib.Stage{
				{Duration: lib.NullDurationFrom(5 * time.Second)},
				{Duration: lib.NullDurationFrom(5 * time.Second)},
				{},
			},
		},
	}
	for name, data := range testdata {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, data.Time, SumStages(data.Stages))
		})
	}
}

func TestProcessStages(t *testing.T) {
	type checkpoint struct {
		D    time.Duration
		Keep bool
		VUs  null.Int
	}
	testdata := map[string]struct {
		Stages      []lib.Stage
		Checkpoints []checkpoint
	}{
		"none": {
			[]lib.Stage{},
			[]checkpoint{
				{0 * time.Second, false, null.NewInt(0, false)},
				{10 * time.Second, false, null.NewInt(0, false)},
				{24 * time.Hour, false, null.NewInt(0, false)},
			},
		},
		"one": {
			[]lib.Stage{
				{Duration: lib.NullDurationFrom(10 * time.Second)},
			},
			[]checkpoint{
				{0 * time.Second, true, null.NewInt(0, false)},
				{1 * time.Second, true, null.NewInt(0, false)},
				{10 * time.Second, true, null.NewInt(0, false)},
				{11 * time.Second, false, null.NewInt(0, false)},
			},
		},
		"one/targeted": {
			[]lib.Stage{
				{Duration: lib.NullDurationFrom(10 * time.Second), Target: null.IntFrom(100)},
			},
			[]checkpoint{
				{0 * time.Second, true, null.IntFrom(0)},
				{1 * time.Second, true, null.IntFrom(10)},
				{2 * time.Second, true, null.IntFrom(20)},
				{3 * time.Second, true, null.IntFrom(30)},
				{4 * time.Second, true, null.IntFrom(40)},
				{5 * time.Second, true, null.IntFrom(50)},
				{6 * time.Second, true, null.IntFrom(60)},
				{7 * time.Second, true, null.IntFrom(70)},
				{8 * time.Second, true, null.IntFrom(80)},
				{9 * time.Second, true, null.IntFrom(90)},
				{10 * time.Second, true, null.IntFrom(100)},
				{11 * time.Second, false, null.IntFrom(100)},
			},
		},
		"two": {
			[]lib.Stage{
				{Duration: lib.NullDurationFrom(5 * time.Second)},
				{Duration: lib.NullDurationFrom(5 * time.Second)},
			},
			[]checkpoint{
				{0 * time.Second, true, null.NewInt(0, false)},
				{1 * time.Second, true, null.NewInt(0, false)},
				{11 * time.Second, false, null.NewInt(0, false)},
			},
		},
		"two/targeted": {
			[]lib.Stage{
				{Duration: lib.NullDurationFrom(5 * time.Second), Target: null.IntFrom(100)},
				{Duration: lib.NullDurationFrom(5 * time.Second), Target: null.IntFrom(0)},
			},
			[]checkpoint{
				{0 * time.Second, true, null.IntFrom(0)},
				{1 * time.Second, true, null.IntFrom(20)},
				{2 * time.Second, true, null.IntFrom(40)},
				{3 * time.Second, true, null.IntFrom(60)},
				{4 * time.Second, true, null.IntFrom(80)},
				{5 * time.Second, true, null.IntFrom(100)},
				{6 * time.Second, true, null.IntFrom(80)},
				{7 * time.Second, true, null.IntFrom(60)},
				{8 * time.Second, true, null.IntFrom(40)},
				{9 * time.Second, true, null.IntFrom(20)},
				{10 * time.Second, true, null.IntFrom(0)},
				{11 * time.Second, false, null.IntFrom(0)},
			},
		},
		"three": {
			[]lib.Stage{
				{Duration: lib.NullDurationFrom(5 * time.Second)},
				{Duration: lib.NullDurationFrom(10 * time.Second)},
				{Duration: lib.NullDurationFrom(15 * time.Second)},
			},
			[]checkpoint{
				{0 * time.Second, true, null.NewInt(0, false)},
				{1 * time.Second, true, null.NewInt(0, false)},
				{15 * time.Second, true, null.NewInt(0, false)},
				{30 * time.Second, true, null.NewInt(0, false)},
				{31 * time.Second, false, null.NewInt(0, false)},
			},
		},
		"three/targeted": {
			[]lib.Stage{
				{Duration: lib.NullDurationFrom(5 * time.Second), Target: null.IntFrom(50)},
				{Duration: lib.NullDurationFrom(5 * time.Second), Target: null.IntFrom(100)},
				{Duration: lib.NullDurationFrom(5 * time.Second), Target: null.IntFrom(0)},
			},
			[]checkpoint{
				{0 * time.Second, true, null.IntFrom(0)},
				{1 * time.Second, true, null.IntFrom(10)},
				{2 * time.Second, true, null.IntFrom(20)},
				{3 * time.Second, true, null.IntFrom(30)},
				{4 * time.Second, true, null.IntFrom(40)},
				{5 * time.Second, true, null.IntFrom(50)},
				{6 * time.Second, true, null.IntFrom(60)},
				{7 * time.Second, true, null.IntFrom(70)},
				{8 * time.Second, true, null.IntFrom(80)},
				{9 * time.Second, true, null.IntFrom(90)},
				{10 * time.Second, true, null.IntFrom(100)},
				{11 * time.Second, true, null.IntFrom(80)},
				{12 * time.Second, true, null.IntFrom(60)},
				{13 * time.Second, true, null.IntFrom(40)},
				{14 * time.Second, true, null.IntFrom(20)},
				{15 * time.Second, true, null.IntFrom(0)},
				{16 * time.Second, false, null.IntFrom(0)},
			},
		},
		"mix": {
			[]lib.Stage{
				{Duration: lib.NullDurationFrom(5 * time.Second), Target: null.IntFrom(20)},
				{Duration: lib.NullDurationFrom(5 * time.Second), Target: null.IntFrom(10)},
				{Duration: lib.NullDurationFrom(2 * time.Second)},
				{Duration: lib.NullDurationFrom(5 * time.Second), Target: null.IntFrom(20)},
				{Duration: lib.NullDurationFrom(2 * time.Second)},
				{Duration: lib.NullDurationFrom(5 * time.Second), Target: null.IntFrom(10)},
			},
			[]checkpoint{
				{0 * time.Second, true, null.IntFrom(0)},

				{1 * time.Second, true, null.IntFrom(4)},
				{2 * time.Second, true, null.IntFrom(8)},
				{3 * time.Second, true, null.IntFrom(12)},
				{4 * time.Second, true, null.IntFrom(16)},
				{5 * time.Second, true, null.IntFrom(20)},

				{6 * time.Second, true, null.IntFrom(18)},
				{7 * time.Second, true, null.IntFrom(16)},
				{8 * time.Second, true, null.IntFrom(14)},
				{9 * time.Second, true, null.IntFrom(12)},
				{10 * time.Second, true, null.IntFrom(10)},

				{11 * time.Second, true, null.IntFrom(10)},
				{12 * time.Second, true, null.IntFrom(10)},

				{13 * time.Second, true, null.IntFrom(12)},
				{14 * time.Second, true, null.IntFrom(14)},
				{15 * time.Second, true, null.IntFrom(16)},
				{16 * time.Second, true, null.IntFrom(18)},
				{17 * time.Second, true, null.IntFrom(20)},

				{18 * time.Second, true, null.IntFrom(20)},
				{19 * time.Second, true, null.IntFrom(20)},

				{20 * time.Second, true, null.IntFrom(18)},
				{21 * time.Second, true, null.IntFrom(16)},
				{22 * time.Second, true, null.IntFrom(14)},
				{23 * time.Second, true, null.IntFrom(12)},
				{24 * time.Second, true, null.IntFrom(10)},
			},
		},
		"infinite": {
			[]lib.Stage{{}},
			[]checkpoint{
				{0 * time.Second, true, null.NewInt(0, false)},
				{1 * time.Minute, true, null.NewInt(0, false)},
				{1 * time.Hour, true, null.NewInt(0, false)},
				{24 * time.Hour, true, null.NewInt(0, false)},
				{365 * 24 * time.Hour, true, null.NewInt(0, false)},
			},
		},
	}
	for name, data := range testdata {
		t.Run(name, func(t *testing.T) {
			for _, ckp := range data.Checkpoints {
				t.Run(ckp.D.String(), func(t *testing.T) {
					vus, keepRunning := ProcessStages(data.Stages, ckp.D)
					assert.Equal(t, ckp.VUs, vus)
					assert.Equal(t, ckp.Keep, keepRunning)
				})
			}
		})
	}
}
