// Copyright 2017 Prospect One https://prospectone.io/. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestInitLatencyCmd(t *testing.T) {
	testCases := map[string]struct {
		args   []string
		gotexp func() (interface{}, interface{})
	}{
		"limit": {[]string{"--limit", "23"}, func() (interface{}, interface{}) { return latencyLimit, 23 }},
	}
	parent := &cobra.Command{}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			latencyCmd.ResetFlags()
			initLatencyCmd(parent)
			if err := latencyCmd.ParseFlags(tc.args); err != nil {
				t.Fatalf("exepected nil; got %v", err)
			}
			flags := dnsPerfCmd.Flags()
			f := flags.Lookup(name)
			if f == nil {
				t.Fatal("expected flag; got nil")
			}
			if got, exp := tc.gotexp(); got != exp {
				t.Fatalf("expected %v; got %v", exp, got)
			}
		})
	}
}

func TestRunLatency(t *testing.T) {
	testCases := map[string]struct {
		from    string
		nodeIDs []int
		exp     string
	}{
		"Location": {"From here", []int{}, `{"target":"example.com","location":"From here","limit":12}`},
		"NodeID":   {"", []int{123}, `{"target":"example.com","nodes":"123","limit":12}`},
	}
	// We're only interested in the first HTTP call, e.g., the one to get the test ID
	// to validate our parameters got passed properly.
	tr := &recordingTransport{}
	c, err := newTestPerfopsClient(tr)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			runLatency(c, "example.com", tc.from, tc.nodeIDs, 12)
			if got, exp := tr.req.URL.Path, "/run/latency"; got != exp {
				t.Fatalf("expected %v; got %v", exp, got)
			}
			if got, exp := reqBody(tr.req), tc.exp; got != exp {
				t.Fatalf("expected %v; got %v", exp, got)
			}
		})
	}
}
