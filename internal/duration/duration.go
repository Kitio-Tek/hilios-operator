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

// Package duration contains conversions between the int32 *Seconds fields
// HILIOS uses on its CRDs and the time.Duration values used by the controller
// runtime. Centralising the conversion avoids subtle off-by-one errors when
// individual fields are clamped or defaulted.
package duration

import "time"

// FromSeconds converts an int32 seconds value into a time.Duration. Negative
// inputs return zero.
func FromSeconds(s int32) time.Duration {
	if s <= 0 {
		return 0
	}
	return time.Duration(s) * time.Second
}

// FromSecondsOr returns FromSeconds(s) when s > 0 and the supplied default
// otherwise.
func FromSecondsOr(s int32, def time.Duration) time.Duration {
	d := FromSeconds(s)
	if d == 0 {
		return def
	}
	return d
}

// Clamp constrains d to the [min, max] interval. min must be <= max.
func Clamp(d, minD, maxD time.Duration) time.Duration {
	if d < minD {
		return minD
	}
	if d > maxD {
		return maxD
	}
	return d
}
