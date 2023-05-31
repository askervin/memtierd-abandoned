// Copyright 2023 Intel Corporation. All Rights Reserved.
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

package main

import (
	"bufio"
	_ "bytes"
	"fmt"
	"os"
	"strings"
	_ "sync"
	"testing"

	"github.com/intel/memtierd/pkg/memtier"
)

func FuzzPrompt(f *testing.F) {
	pidwatcherCommonArgs := " -listener log -config-dump -poll -start -stop -dump"
	trackerCommonArgs := " -config-dump -start -stop -dump"
	policyCommonArgs := " -config-dump -start -stop -dump"
	testcases := []string{
		"help",
		"nop",
		"pidwatcher -ls",
		"pidwatcher -create pidlist -config {\"Pids\":[42,4242]}" + pidwatcherCommonArgs,
		"pidwatcher -create cgroups -config {\"IntervalMs\":10000,\"Cgroups\":[\"/sys/fs/cgroup/memtierd-test\"]}" + pidwatcherCommonArgs,
		"pidwatcher -create proc -config {\"IntervalMs\":10000}" + pidwatcherCommonArgs,
		"pidwatcher -create pidlist -config {\"Pids\":[42,4242]}" + pidwatcherCommonArgs,
		"pidwatcher -create filter -config {}" + pidwatcherCommonArgs,
		"policy -ls",
		"policy -create age -config {\"IntervalMs\":10000}" + policyCommonArgs,
		"tracker -ls",
		"tracker -create damon" + trackerCommonArgs,
		"tracker -create idlepage" + trackerCommonArgs,
		"tracker -create softdirty" + trackerCommonArgs,
		"q",
		"mover -pages-to 1",
		"stats",
	}

	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, input string) {
		// var promptInBuf bytes.Buffer
		// var promptOutBuf bytes.Buffer
		// var wg sync.WaitGroup
		fmt.Printf("input: %q\n", input)
		if strings.Contains(input, "|") {
			// Do not fuzz inputs with pipes, as it would
			// execute fuzzed strings in shell.
			return
		}
		// promptIn := bufio.NewReadWriter(
		// 	bufio.NewReader(&promptInBuf),
		// 	bufio.NewWriter(&promptInBuf))
		// promptOut := bufio.NewWriter(&promptOutBuf)
		prompt := memtier.NewPrompt("memtierd-fuzzed> ", bufio.NewReader(strings.NewReader("")), bufio.NewWriter(os.Stderr))
		prompt.SetEcho(true)

		// promptIn.WriteString(input)
		// if len(input) > 0 && input[len(input)-1] != '\n' {
		//	promptIn.WriteString("\n")
		// }
		// promptIn.Flush()
		// wg.Add(1)
		// go func() {
		// 	defer wg.Done()
		// 	prompt.Interact()
		// }()
		// promptIn.WriteString("\nq\n")
		// promptIn.Flush()
		// wg.Wait()
		prompt.RunCmdString(input)
		// fmt.Printf("---response-begin---\n%s---response-end---\n", promptOutBuf.String())

	})
}
