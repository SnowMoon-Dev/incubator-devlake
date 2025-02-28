/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package parser

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	git "github.com/libgit2/git2go/v33"
	ssh2 "golang.org/x/crypto/ssh"
)

const DefaultUser = "git"

func cloneOverSSH(url, dir, passphrase string, pk []byte) error {
	key, err := ssh.NewPublicKeys(DefaultUser, pk, passphrase)
	if err != nil {
		return err
	}
	key.HostKeyCallbackHelper = ssh.HostKeyCallbackHelper{
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh2.PublicKey) error {
			return nil
		},
	}
	_, err = gogit.PlainClone(dir, true, &gogit.CloneOptions{
		URL:  url,
		Auth: key,
	})
	if err != nil {
		return err
	}
	return nil
}

func (l *LibGit2) CloneOverHTTP(repoId, url, user, password, proxy string) error {
	cloneOptions := &git.CloneOptions{Bare: true}
	if proxy != "" {
		cloneOptions.FetchOptions.ProxyOptions.Type = git.ProxyTypeSpecified
		cloneOptions.FetchOptions.ProxyOptions.Url = proxy
	}
	if user != "" {
		auth := fmt.Sprintf("Authorization: Basic %s", base64.StdEncoding.EncodeToString([]byte(user+":"+password)))
		cloneOptions.FetchOptions.Headers = []string{auth}
	}
	dir, err := ioutil.TempDir("", "gitextractor")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	repo, err := git.Clone(url, dir, cloneOptions)
	if err != nil {
		return err
	}
	return l.run(repo, repoId)
}

func (l *LibGit2) CloneOverSSH(repoId, url, privateKey, passphrase string) error {
	dir, err := ioutil.TempDir("", "gitextractor")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	pk, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return err
	}
	err = cloneOverSSH(url, dir, passphrase, pk)
	if err != nil {
		return err
	}
	return l.LocalRepo(dir, repoId)
}
