// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

import (
	"math"
	"strings"
	"testing"

	"golang.org/x/exp/constraints"
)

var raceEnabled bool

var equalIntTests = []struct {
	s1, s2 []int
	want   bool
}{
	{
		[]int{1},
		nil,
		false,
	},
	{
		[]int{},
		nil,
		true,
	},
	{
		[]int{1, 2, 3},
		[]int{1, 2, 3},
		true,
	},
	{
		[]int{1, 2, 3},
		[]int{1, 2, 3, 4},
		false,
	},
}

var equalFloatTests = []struct {
	s1, s2       []float64
	wantEqual    bool
	wantEqualNaN bool
}{
	{
		[]float64{1, 2},
		[]float64{1, 2},
		true,
		true,
	},
	{
		[]float64{1, 2, math.NaN()},
		[]float64{1, 2, math.NaN()},
		false,
		true,
	},
}

func TestEqual(t *testing.T) {
	for _, test := range equalIntTests {
		if got := Equal(test.s1, test.s2); got != test.want {
			t.Errorf("Equal(%v, %v) = %t, want %t", test.s1, test.s2, got, test.want)
		}
	}
	for _, test := range equalFloatTests {
		if got := Equal(test.s1, test.s2); got != test.wantEqual {
			t.Errorf("Equal(%v, %v) = %t, want %t", test.s1, test.s2, got, test.wantEqual)
		}
	}
}

// equal is simply ==.
func equal[T comparable](v1, v2 T) bool {
	return v1 == v2
}

// equalNaN is like == except that all NaNs are equal.
func equalNaN[T comparable](v1, v2 T) bool {
	isNaN := func(f T) bool { return f != f }
	return v1 == v2 || (isNaN(v1) && isNaN(v2))
}

// offByOne returns true if integers v1 and v2 differ by 1.
func offByOne[E constraints.Integer](v1, v2 E) bool {
	return v1 == v2+1 || v1 == v2-1
}

func TestEqualFunc(t *testing.T) {
	for _, test := range equalIntTests {
		if got := EqualFunc(test.s1, test.s2, equal[int]); got != test.want {
			t.Errorf("EqualFunc(%v, %v, equal[int]) = %t, want %t", test.s1, test.s2, got, test.want)
		}
	}
	for _, test := range equalFloatTests {
		if got := EqualFunc(test.s1, test.s2, equal[float64]); got != test.wantEqual {
			t.Errorf("Equal(%v, %v, equal[float64]) = %t, want %t", test.s1, test.s2, got, test.wantEqual)
		}
		if got := EqualFunc(test.s1, test.s2, equalNaN[float64]); got != test.wantEqualNaN {
			t.Errorf("Equal(%v, %v, equalNaN[float64]) = %t, want %t", test.s1, test.s2, got, test.wantEqualNaN)
		}
	}

	s1 := []int{1, 2, 3}
	s2 := []int{2, 3, 4}
	if EqualFunc(s1, s1, offByOne[int]) {
		t.Errorf("EqualFunc(%v, %v, offByOne) = true, want false", s1, s1)
	}
	if !EqualFunc(s1, s2, offByOne[int]) {
		t.Errorf("EqualFunc(%v, %v, offByOne) = false, want true", s1, s2)
	}

	s3 := []string{"a", "b", "c"}
	s4 := []string{"A", "B", "C"}
	if !EqualFunc(s3, s4, strings.EqualFold) {
		t.Errorf("EqualFunc(%v, %v, strings.EqualFold) = false, want true", s3, s4)
	}

	cmpIntString := func(v1 int, v2 string) bool {
		return string(rune(v1)-1+'a') == v2
	}
	if !EqualFunc(s1, s3, cmpIntString) {
		t.Errorf("EqualFunc(%v, %v, cmpIntString) = false, want true", s1, s3)
	}
}

var compareIntTests = []struct {
	s1, s2 []int
	want   int
}{
	{
		[]int{1, 2, 3},
		[]int{1, 2, 3, 4},
		-1,
	},
	{
		[]int{1, 2, 3, 4},
		[]int{1, 2, 3},
		+1,
	},
	{
		[]int{1, 2, 3},
		[]int{1, 4, 3},
		-1,
	},
	{
		[]int{1, 4, 3},
		[]int{1, 2, 3},
		+1,
	},
}

var compareFloatTests = []struct {
	s1, s2 []float64
	want   int
}{
	{
		[]float64{1, 2, math.NaN()},
		[]float64{1, 2, math.NaN()},
		0,
	},
	{
		[]float64{1, math.NaN(), 3},
		[]float64{1, math.NaN(), 4},
		-1,
	},
	{
		[]float64{1, math.NaN(), 3},
		[]float64{1, 2, 4},
		-1,
	},
	{
		[]float64{1, math.NaN(), 3},
		[]float64{1, 2, math.NaN()},
		0,
	},
	{
		[]float64{1, math.NaN(), 3, 4},
		[]float64{1, 2, math.NaN()},
		+1,
	},
}

func TestCompare(t *testing.T) {
	intWant := func(want bool) string {
		if want {
			return "0"
		}
		return "!= 0"
	}
	for _, test := range equalIntTests {
		if got := Compare(test.s1, test.s2); (got == 0) != test.want {
			t.Errorf("Compare(%v, %v) = %d, want %s", test.s1, test.s2, got, intWant(test.want))
		}
	}
	for _, test := range equalFloatTests {
		if got := Compare(test.s1, test.s2); (got == 0) != test.wantEqualNaN {
			t.Errorf("Compare(%v, %v) = %d, want %s", test.s1, test.s2, got, intWant(test.wantEqualNaN))
		}
	}

	for _, test := range compareIntTests {
		if got := Compare(test.s1, test.s2); got != test.want {
			t.Errorf("Compare(%v, %v) = %d, want %d", test.s1, test.s2, got, test.want)
		}
	}
	for _, test := range compareFloatTests {
		if got := Compare(test.s1, test.s2); got != test.want {
			t.Errorf("Compare(%v, %v) = %d, want %d", test.s1, test.s2, got, test.want)
		}
	}
}

func equalToCmp[T comparable](eq func(T, T) bool) func(T, T) int {
	return func(v1, v2 T) int {
		if eq(v1, v2) {
			return 0
		}
		return 1
	}
}

func cmp[T constraints.Ordered](v1, v2 T) int {
	if v1 < v2 {
		return -1
	} else if v1 > v2 {
		return 1
	} else {
		return 0
	}
}

func TestCompareFunc(t *testing.T) {
	intWant := func(want bool) string {
		if want {
			return "0"
		}
		return "!= 0"
	}
	for _, test := range equalIntTests {
		if got := CompareFunc(test.s1, test.s2, equalToCmp(equal[int])); (got == 0) != test.want {
			t.Errorf("CompareFunc(%v, %v, equalToCmp(equal[int])) = %d, want %s", test.s1, test.s2, got, intWant(test.want))
		}
	}
	for _, test := range equalFloatTests {
		if got := CompareFunc(test.s1, test.s2, equalToCmp(equal[float64])); (got == 0) != test.wantEqual {
			t.Errorf("CompareFunc(%v, %v, equalToCmp(equal[float64])) = %d, want %s", test.s1, test.s2, got, intWant(test.wantEqual))
		}
	}

	for _, test := range compareIntTests {
		if got := CompareFunc(test.s1, test.s2, cmp[int]); got != test.want {
			t.Errorf("CompareFunc(%v, %v, cmp[int]) = %d, want %d", test.s1, test.s2, got, test.want)
		}
	}
	for _, test := range compareFloatTests {
		if got := CompareFunc(test.s1, test.s2, cmp[float64]); got != test.want {
			t.Errorf("CompareFunc(%v, %v, cmp[float64]) = %d, want %d", test.s1, test.s2, got, test.want)
		}
	}

	s1 := []int{1, 2, 3}
	s2 := []int{2, 3, 4}
	if got := CompareFunc(s1, s2, equalToCmp(offByOne[int])); got != 0 {
		t.Errorf("CompareFunc(%v, %v, offByOne) = %d, want 0", s1, s2, got)
	}

	s3 := []string{"a", "b", "c"}
	s4 := []string{"A", "B", "C"}
	if got := CompareFunc(s3, s4, strings.Compare); got != 1 {
		t.Errorf("CompareFunc(%v, %v, strings.Compare) = %d, want 1", s3, s4, got)
	}

	compareLower := func(v1, v2 string) int {
		return strings.Compare(strings.ToLower(v1), strings.ToLower(v2))
	}
	if got := CompareFunc(s3, s4, compareLower); got != 0 {
		t.Errorf("CompareFunc(%v, %v, compareLower) = %d, want 0", s3, s4, got)
	}

	cmpIntString := func(v1 int, v2 string) int {
		return strings.Compare(string(rune(v1)-1+'a'), v2)
	}
	if got := CompareFunc(s1, s3, cmpIntString); got != 0 {
		t.Errorf("CompareFunc(%v, %v, cmpIntString) = %d, want 0", s1, s3, got)
	}
}

var indexTests = []struct {
	s    []int
	v    int
	want int
}{
	{
		nil,
		0,
		-1,
	},
	{
		[]int{},
		0,
		-1,
	},
	{
		[]int{1, 2, 3},
		2,
		1,
	},
	{
		[]int{1, 2, 2, 3},
		2,
		1,
	},
	{
		[]int{1, 2, 3, 2},
		2,
		1,
	},
}

func TestIndex(t *testing.T) {
	for _, test := range indexTests {
		if got := Index(test.s, test.v); got != test.want {
			t.Errorf("Index(%v, %v) = %d, want %d", test.s, test.v, got, test.want)
		}
	}
}

func equalToIndex[T any](f func(T, T) bool, v1 T) func(T) bool {
	return func(v2 T) bool {
		return f(v1, v2)
	}
}

func TestIndexFunc(t *testing.T) {
	for _, test := range indexTests {
		if got := IndexFunc(test.s, equalToIndex(equal[int], test.v)); got != test.want {
			t.Errorf("IndexFunc(%v, equalToIndex(equal[int], %v)) = %d, want %d", test.s, test.v, got, test.want)
		}
	}

	s1 := []string{"hi", "HI"}
	if got := IndexFunc(s1, equalToIndex(equal[string], "HI")); got != 1 {
		t.Errorf("IndexFunc(%v, equalToIndex(equal[string], %q)) = %d, want %d", s1, "HI", got, 1)
	}
	if got := IndexFunc(s1, equalToIndex(strings.EqualFold, "HI")); got != 0 {
		t.Errorf("IndexFunc(%v, equalToIndex(strings.EqualFold, %q)) = %d, want %d", s1, "HI", got, 0)
	}
}

func TestContains(t *testing.T) {
	for _, test := range indexTests {
		if got := Contains(test.s, test.v); got != (test.want != -1) {
			t.Errorf("Contains(%v, %v) = %t, want %t", test.s, test.v, got, test.want != -1)
		}
	}
}

var insertTests = []struct {
	s    []int
	i    int
	add  []int
	want []int
}{
	{
		[]int{1, 2, 3},
		0,
		[]int{4},
		[]int{4, 1, 2, 3},
	},
	{
		[]int{1, 2, 3},
		1,
		[]int{4},
		[]int{1, 4, 2, 3},
	},
	{
		[]int{1, 2, 3},
		3,
		[]int{4},
		[]int{1, 2, 3, 4},
	},
	{
		[]int{1, 2, 3},
		2,
		[]int{4, 5},
		[]int{1, 2, 4, 5, 3},
	},
}

func TestInsert(t *testing.T) {
	s := []int{1, 2, 3}
	if got := Insert(s, 0); !Equal(got, s) {
		t.Errorf("Insert(%v, 0) = %v, want %v", s, got, s)
	}
	for _, test := range insertTests {
		copy := Clone(test.s)
		if got := Insert(copy, test.i, test.add...); !Equal(got, test.want) {
			t.Errorf("Insert(%v, %d, %v...) = %v, want %v", test.s, test.i, test.add, got, test.want)
		}
	}
}

var deleteTests = []struct {
	s    []int
	i, j int
	want []int
}{
	{
		[]int{1, 2, 3},
		0,
		0,
		[]int{1, 2, 3},
	},
	{
		[]int{1, 2, 3},
		0,
		1,
		[]int{2, 3},
	},
	{
		[]int{1, 2, 3},
		3,
		3,
		[]int{1, 2, 3},
	},
	{
		[]int{1, 2, 3},
		0,
		2,
		[]int{3},
	},
	{
		[]int{1, 2, 3},
		0,
		3,
		[]int{},
	},
}

func TestDelete(t *testing.T) {
	for _, test := range deleteTests {
		copy := Clone(test.s)
		if got := Delete(copy, test.i, test.j); !Equal(got, test.want) {
			t.Errorf("Delete(%v, %d, %d) = %v, want %v", test.s, test.i, test.j, got, test.want)
		}
	}
}

func panics(f func()) (b bool) {
	defer func() {
		if x := recover(); x != nil {
			b = true
		}
	}()
	f()
	return false
}

func TestDeletePanics(t *testing.T) {
	for _, test := range []struct {
		name string
		s    []int
		i, j int
	}{
		{"with negative first index", []int{42}, -2, 1},
		{"with negative second index", []int{42}, 1, -1},
		{"with out-of-bounds first index", []int{42}, 2, 3},
		{"with out-of-bounds second index", []int{42}, 0, 2},
		{"with invalid i>j", []int{42}, 1, 0},
	} {
		if !panics(func() { Delete(test.s, test.i, test.j) }) {
			t.Errorf("Delete %s: got no panic, want panic", test.name)
		}
	}
}

func TestClone(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := Clone(s1)
	if !Equal(s1, s2) {
		t.Errorf("Clone(%v) = %v, want %v", s1, s2, s1)
	}
	s1[0] = 4
	want := []int{1, 2, 3}
	if !Equal(s2, want) {
		t.Errorf("Clone(%v) changed unexpectedly to %v", want, s2)
	}
	if got := Clone([]int(nil)); got != nil {
		t.Errorf("Clone(nil) = %#v, want nil", got)
	}
	if got := Clone(s1[:0]); got == nil || len(got) != 0 {
		t.Errorf("Clone(%v) = %#v, want %#v", s1[:0], got, s1[:0])
	}
}

var compactTests = []struct {
	name string
	s    []int
	want []int
}{
	{
		"nil",
		nil,
		nil,
	},
	{
		"one",
		[]int{1},
		[]int{1},
	},
	{
		"sorted",
		[]int{1, 2, 3},
		[]int{1, 2, 3},
	},
	{
		"1 item",
		[]int{1, 1, 2},
		[]int{1, 2},
	},
	{
		"unsorted",
		[]int{1, 2, 1},
		[]int{1, 2, 1},
	},
	{
		"many",
		[]int{1, 2, 2, 3, 3, 4},
		[]int{1, 2, 3, 4},
	},
}

func TestCompact(t *testing.T) {
	for _, test := range compactTests {
		copy := Clone(test.s)
		if got := Compact(copy); !Equal(got, test.want) {
			t.Errorf("Compact(%v) = %v, want %v", test.s, got, test.want)
		}
	}
}

func BenchmarkCompact(b *testing.B) {
	for _, c := range compactTests {
		b.Run(c.name, func(b *testing.B) {
			ss := make([]int, 0, 64)
			for k := 0; k < b.N; k++ {
				ss = ss[:0]
				ss = append(ss, c.s...)
				_ = Compact(ss)
			}
		})
	}

}

func TestCompactFunc(t *testing.T) {
	for _, test := range compactTests {
		copy := Clone(test.s)
		if got := CompactFunc(copy, equal[int]); !Equal(got, test.want) {
			t.Errorf("CompactFunc(%v, equal[int]) = %v, want %v", test.s, got, test.want)
		}
	}

	s1 := []string{"a", "a", "A", "B", "b"}
	copy := Clone(s1)
	want := []string{"a", "B"}
	if got := CompactFunc(copy, strings.EqualFold); !Equal(got, want) {
		t.Errorf("CompactFunc(%v, strings.EqualFold) = %v, want %v", s1, got, want)
	}
}

func TestGrow(t *testing.T) {
	s1 := []int{1, 2, 3}

	copy := Clone(s1)
	s2 := Grow(copy, 1000)
	if !Equal(s1, s2) {
		t.Errorf("Grow(%v) = %v, want %v", s1, s2, s1)
	}
	if cap(s2) < 1000+len(s1) {
		t.Errorf("after Grow(%v) cap = %d, want >= %d", s1, cap(s2), 1000+len(s1))
	}

	// Test mutation of elements between length and capacity.
	copy = Clone(s1)
	s3 := Grow(copy[:1], 2)[:3]
	if !Equal(s1, s3) {
		t.Errorf("Grow should not mutate elements between length and capacity")
	}
	s3 = Grow(copy[:1], 1000)[:3]
	if !Equal(s1, s3) {
		t.Errorf("Grow should not mutate elements between length and capacity")
	}

	// Test number of allocations.
	if n := testing.AllocsPerRun(100, func() { Grow(s2, cap(s2)-len(s2)) }); n != 0 {
		t.Errorf("Grow should not allocate when given sufficient capacity; allocated %v times", n)
	}
	if n := testing.AllocsPerRun(100, func() { Grow(s2, cap(s2)-len(s2)+1) }); n != 1 {
		errorf := t.Errorf
		if raceEnabled {
			errorf = t.Logf // this allocates multiple times in race detector mode
		}
		errorf("Grow should allocate once when given insufficient capacity; allocated %v times", n)
	}

	// Test for negative growth sizes.
	var gotPanic bool
	func() {
		defer func() { gotPanic = recover() != nil }()
		Grow(s1, -1)
	}()
	if !gotPanic {
		t.Errorf("Grow(-1) did not panic; expected a panic")
	}
}

func TestClip(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5, 6}[:3]
	orig := Clone(s1)
	if len(s1) != 3 {
		t.Errorf("len(%v) = %d, want 3", s1, len(s1))
	}
	if cap(s1) < 6 {
		t.Errorf("cap(%v[:3]) = %d, want >= 6", orig, cap(s1))
	}
	s2 := Clip(s1)
	if !Equal(s1, s2) {
		t.Errorf("Clip(%v) = %v, want %v", s1, s2, s1)
	}
	if cap(s2) != 3 {
		t.Errorf("cap(Clip(%v)) = %d, want 3", orig, cap(s2))
	}
}

// naiveReplace is a baseline implementation to the Replace function.
func naiveReplace[S ~[]E, E any](s S, i, j int, v ...E) S {
	s = Delete(s, i, j)
	s = Insert(s, i, v...)
	return s
}

func TestReplace(t *testing.T) {
	for _, test := range []struct {
		s, v []int
		i, j int
	}{
		{}, // all zero value
		{
			s: []int{1, 2, 3, 4},
			v: []int{5},
			i: 1,
			j: 2,
		},
		{
			s: []int{1, 2, 3, 4},
			v: []int{5, 6, 7, 8},
			i: 1,
			j: 2,
		},
		{
			s: func() []int {
				s := make([]int, 3, 20)
				s[0] = 0
				s[1] = 1
				s[2] = 2
				return s
			}(),
			v: []int{3, 4, 5, 6, 7},
			i: 0,
			j: 1,
		},
	} {
		ss, vv := Clone(test.s), Clone(test.v)
		want := naiveReplace(ss, test.i, test.j, vv...)
		got := Replace(test.s, test.i, test.j, test.v...)
		if !Equal(got, want) {
			t.Errorf("Replace(%v, %v, %v, %v) = %v, want %v", test.s, test.i, test.j, test.v, got, want)
		}
	}
}

func BenchmarkReplace(b *testing.B) {
	cases := []struct {
		name string
		s, v func() []int
		i, j int
	}{
		{
			name: "fast",
			s: func() []int {
				return make([]int, 100)
			},
			v: func() []int {
				return make([]int, 20)
			},
			i: 10,
			j: 40,
		},
		{
			name: "slow",
			s: func() []int {
				return make([]int, 100)
			},
			v: func() []int {
				return make([]int, 20)
			},
			i: 0,
			j: 2,
		},
	}

	for _, c := range cases {
		b.Run("naive-"+c.name, func(b *testing.B) {
			for k := 0; k < b.N; k++ {
				s := c.s()
				v := c.v()
				_ = naiveReplace(s, c.i, c.j, v...)
			}
		})
		b.Run("optimized-"+c.name, func(b *testing.B) {
			for k := 0; k < b.N; k++ {
				s := c.s()
				v := c.v()
				_ = Replace(s, c.i, c.j, v...)
			}
		})
	}

}
