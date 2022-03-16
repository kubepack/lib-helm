/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by set-gen. DO NOT EDIT.

package sets

import (
	"reflect"
	"sort"
)

// sets.Uint32 is a set of uint32s, implemented via map[uint32]struct{} for minimal memory consumption.
type Uint32 map[uint32]Empty

// NewUint32 creates a Uint32 from a list of values.
func NewUint32(items ...uint32) Uint32 {
	ss := make(Uint32, len(items))
	ss.Insert(items...)
	return ss
}

// Uint32KeySet creates a Uint32 from a keys of a map[uint32](? extends interface{}).
// If the value passed in is not actually a map, this will panic.
func Uint32KeySet(theMap interface{}) Uint32 {
	v := reflect.ValueOf(theMap)
	ret := Uint32{}

	for _, keyValue := range v.MapKeys() {
		ret.Insert(keyValue.Interface().(uint32))
	}
	return ret
}

// Insert adds items to the set.
func (s Uint32) Insert(items ...uint32) Uint32 {
	for _, item := range items {
		s[item] = Empty{}
	}
	return s
}

// Delete removes all items from the set.
func (s Uint32) Delete(items ...uint32) Uint32 {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

// Has returns true if and only if item is contained in the set.
func (s Uint32) Has(item uint32) bool {
	_, contained := s[item]
	return contained
}

// HasAll returns true if and only if all items are contained in the set.
func (s Uint32) HasAll(items ...uint32) bool {
	for _, item := range items {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

// HasAny returns true if any items are contained in the set.
func (s Uint32) HasAny(items ...uint32) bool {
	for _, item := range items {
		if s.Has(item) {
			return true
		}
	}
	return false
}

// Difference returns a set of objects that are not in s2
// For example:
// s1 = {a1, a2, a3}
// s2 = {a1, a2, a4, a5}
// s1.Difference(s2) = {a3}
// s2.Difference(s1) = {a4, a5}
func (s Uint32) Difference(s2 Uint32) Uint32 {
	result := NewUint32()
	for key := range s {
		if !s2.Has(key) {
			result.Insert(key)
		}
	}
	return result
}

// Union returns a new set which includes items in either s1 or s2.
// For example:
// s1 = {a1, a2}
// s2 = {a3, a4}
// s1.Union(s2) = {a1, a2, a3, a4}
// s2.Union(s1) = {a1, a2, a3, a4}
func (s1 Uint32) Union(s2 Uint32) Uint32 {
	result := NewUint32()
	for key := range s1 {
		result.Insert(key)
	}
	for key := range s2 {
		result.Insert(key)
	}
	return result
}

// Intersection returns a new set which includes the item in BOTH s1 and s2
// For example:
// s1 = {a1, a2}
// s2 = {a2, a3}
// s1.Intersection(s2) = {a2}
func (s1 Uint32) Intersection(s2 Uint32) Uint32 {
	var walk, other Uint32
	result := NewUint32()
	if s1.Len() < s2.Len() {
		walk = s1
		other = s2
	} else {
		walk = s2
		other = s1
	}
	for key := range walk {
		if other.Has(key) {
			result.Insert(key)
		}
	}
	return result
}

// IsSuperset returns true if and only if s1 is a superset of s2.
func (s1 Uint32) IsSuperset(s2 Uint32) bool {
	for item := range s2 {
		if !s1.Has(item) {
			return false
		}
	}
	return true
}

// Equal returns true if and only if s1 is equal (as a set) to s2.
// Two sets are equal if their membership is identical.
// (In practice, this means same elements, order doesn't matter)
func (s1 Uint32) Equal(s2 Uint32) bool {
	return len(s1) == len(s2) && s1.IsSuperset(s2)
}

type sortableSliceOfUint32 []uint32

func (s sortableSliceOfUint32) Len() int           { return len(s) }
func (s sortableSliceOfUint32) Less(i, j int) bool { return lessUint32(s[i], s[j]) }
func (s sortableSliceOfUint32) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// List returns the contents as a sorted uint32 slice.
func (s Uint32) List() []uint32 {
	res := make(sortableSliceOfUint32, 0, len(s))
	for key := range s {
		res = append(res, key)
	}
	sort.Sort(res)
	return []uint32(res)
}

// UnsortedList returns the slice with contents in random order.
func (s Uint32) UnsortedList() []uint32 {
	res := make([]uint32, 0, len(s))
	for key := range s {
		res = append(res, key)
	}
	return res
}

// Returns a single element from the set.
func (s Uint32) PopAny() (uint32, bool) {
	for key := range s {
		s.Delete(key)
		return key, true
	}
	var zeroValue uint32
	return zeroValue, false
}

// Len returns the size of the set.
func (s Uint32) Len() int {
	return len(s)
}

func lessUint32(lhs, rhs uint32) bool {
	return lhs < rhs
}
