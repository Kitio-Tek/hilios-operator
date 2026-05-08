/*
Copyright 2026.

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

// Package safeint provides bounded conversions between Go integer types so
// that callers can satisfy gosec G115 ("integer overflow") without silencing
// the rule. The conversions clamp values to the destination type's limits
// rather than wrapping, which is the right behaviour when we copy a count
// (for example, len(items)) into a status int32 field.
package safeint

import "math"

// Int32 clamps i to the int32 range and returns it. Negative values map to 0 (counts are
// non-negative by definition), values larger than math.MaxInt32 are clamped to
// math.MaxInt32.
func Int32(i int) int32 {
	if i <= 0 {
		return 0
	}
	if i > math.MaxInt32 {
		return math.MaxInt32
	}
	return int32(i)
}

// Int32From64 clamps i (typically a duration in seconds) to the int32 range.
func Int32From64(i int64) int32 {
	if i <= 0 {
		return 0
	}
	if i > math.MaxInt32 {
		return math.MaxInt32
	}
	return int32(i)
}

// Int32EqualsInt reports whether a and b are equal, performing the comparison
// in int space so callers do not need to convert b to int32 (which would
// trigger gosec G115).
func Int32EqualsInt(a int32, b int) bool {
	return int(a) == b
}

// Int32Min returns the smaller of two int32 values. It is a tiny helper to
// avoid importing math purely for the comparison.
func Int32Min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

// Int32Max returns the larger of two int32 values.
func Int32Max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
