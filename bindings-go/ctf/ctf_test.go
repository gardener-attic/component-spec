// Copyright 2020 Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ctf_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gardener/component-spec/bindings-go/ctf"
	"github.com/mandelsoft/vfs/pkg/memoryfs"
	"github.com/mandelsoft/vfs/pkg/vfs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CTF Archive", func() {

	Context("build", func() {
		It("should build a ctf archive", func() {
			ctx := context.Background()
			defer ctx.Done()
			ca, err := ctf.OpenComponentArchive(nil, "./testdata/component-01")
			Expect(err).ToNot(HaveOccurred())
			Expect(ca.ComponentDescriptor.Resources).To(HaveLen(1))

			file, err := ioutil.TempFile("", ".ctf")
			Expect(err).ToNot(HaveOccurred())

			name := file.Name()
			file.Close()
			os.Remove(name)
			defer os.Remove(name)

			fmt.Println("file: " + name)
			c, err := ctf.OpenCTF(nil, name, ctf.CTF_TAR)
			Expect(err).ToNot(HaveOccurred())
			err = c.AddComponentArchive(ca, ctf.ArchiveFormatTar)
			Expect(err).ToNot(HaveOccurred())
			err = c.Write()
			Expect(err).ToNot(HaveOccurred())
			c.Close()
			checkCTF(nil, name)
		})

		It("should build a ctf archive and write it to a filesystem", func() {
			ctx := context.Background()
			defer ctx.Done()
			ca, err := ctf.OpenComponentArchive(nil, "./testdata/component-01")
			Expect(err).ToNot(HaveOccurred())
			Expect(ca.ComponentDescriptor.Resources).To(HaveLen(1))

			fs := memoryfs.New()

			c, err := ctf.OpenCTF(fs, "archive", ctf.CTF_TAR)
			Expect(err).ToNot(HaveOccurred())
			err = c.AddComponentArchive(ca, ctf.ArchiveFormatTarGzip)
			Expect(err).ToNot(HaveOccurred())
			err = c.Write()
			Expect(err).ToNot(HaveOccurred())
			Expect(c.WriteToFilesystem(fs, "folder")).To(Succeed())
			c.Close()
			fis, err := vfs.ReadDir(fs, "folder")
			Expect(err).To(Succeed())
			fmt.Printf("found %d entries\n", len(fis))
			Expect(len(fis)).To(Equal(1))
			checkCTF(fs, "folder")
			//checkCTF(fs, "archive")
		})

		It("should build a ctf folder", func() {
			ctx := context.Background()
			defer ctx.Done()
			ca, err := ctf.OpenComponentArchive(nil, "./testdata/component-01")
			Expect(err).ToNot(HaveOccurred())
			Expect(ca.ComponentDescriptor.Resources).To(HaveLen(1))

			file, err := ioutil.TempFile("", ".ctf")
			Expect(err).ToNot(HaveOccurred())

			name := file.Name()
			tar := "name" + ".tar"
			file.Close()
			os.Remove(name)
			defer os.Remove(name)
			defer os.Remove(tar)

			fmt.Println("file: " + name)
			c, err := ctf.OpenCTF(nil, name, ctf.CTF_DIR)
			Expect(err).ToNot(HaveOccurred())
			err = c.AddComponentArchive(ca, ctf.ArchiveFormatTar)
			Expect(err).ToNot(HaveOccurred())
			err = c.WriteToArchive(nil, tar, ctf.ArchiveFormatTarGzip)
			Expect(err).ToNot(HaveOccurred())
			c.Close()
			checkCTF(nil, name)
			checkCTF(nil, tar)
		})
	})
})

func checkCTF(fs vfs.FileSystem, name string) {
	c, err := ctf.OpenCTF(fs, name, ctf.CTF_OPEN)
	Expect(err).ToNot(HaveOccurred())
	defer c.Close()
	cnt := 0
	c.Walk(func(ca *ctf.ComponentArchive) error {
		fmt.Println(ca.Digest())
		cnt++
		return nil
	})
	Expect(cnt).To(Equal(1))
}
