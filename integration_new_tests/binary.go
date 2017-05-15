// Copyright (C) 2016-Present Pivotal Software, Inc. All rights reserved.
// This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

package integration_new_tests

import (
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

type Binary struct {
	srcPath    string
	binaryPath string
}

func NewBinary(srcPath string) *Binary {
	return &Binary{srcPath: srcPath, binaryPath: ""}
}

func (b *Binary) Path() string {
	if b.binaryPath == "" {
		var err error
		b.binaryPath, err = gexec.Build(b.srcPath)
		Expect(err).NotTo(HaveOccurred())
	}
	return b.binaryPath
}
