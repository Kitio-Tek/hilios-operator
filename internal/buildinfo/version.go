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

// Package buildinfo exposes build-time identity strings: version, git commit,
// build date. Values are populated via -ldflags at build time. Tests and the
// /version log line use the same source so they cannot drift.
package buildinfo

// Variables populated through -ldflags at build time. Defaults are "dev"
// values so unbuilt binaries (for example, those produced by `go run`) print
// something readable.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// Info returns a one-line summary suitable for logging.
func Info() string {
	return "hilios-operator " + Version + " (commit " + Commit + ", built " + Date + ")"
}
