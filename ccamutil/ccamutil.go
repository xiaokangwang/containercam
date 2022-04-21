package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/xiaokangwang/containercam"
	"net/http"
	"os"
)

func main() {
	username := flag.String("username", "", "")
	password := flag.String("password", "", "")
	repoName := flag.String("repo", "", "")
	url := flag.String("url", "", "")

	action := flag.String("action", "", "")
	hash := flag.String("hash", "", "")
	data := flag.String("data", "", "")

	tagName := flag.String("tag", "", "")

	flag.Parse()
	authvalue := &authn.Basic{Username: *username, Password: *password}

	repo, err := name.NewRepository(*repoName, name.WithDefaultRegistry(*url))
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	// Construct an http.Client that is authorized to pull from gcr.io/google-containers/pause.
	scopes := []string{repo.Scope(transport.PushScope)}
	t, err := transport.New(repo.Registry, authvalue, http.DefaultTransport, scopes)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	client := &http.Client{Transport: t}

	reg := &registry.Registry{URL: "https://" + *url, Client: client}
	reg.Logf = func(format string, args ...interface{}) {
		fmt.Printf(format+"\n", args...)
	}

	err = reg.Ping()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	var storage *containercam.Storage
	if *tagName == "" {
		storage = containercam.NewStorageBackend(reg, *repoName)
	} else {
		storage = containercam.NewStorageBackendWithCustomNameFunc(reg, *repoName, func(sha256value string) string {
			return *tagName
		})
	}

	switch *action {
	case "upload":
		dataContent, err := os.ReadFile(*data)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		hashValue := sha256.Sum256(dataContent)
		hashstr := hex.EncodeToString(hashValue[:])
		fmt.Println(hashstr)
		err = storage.UploadByHash(hashstr, dataContent)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	case "download":
		dataContent, err := storage.DownloadByHash(*hash)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		err = os.WriteFile(*data, dataContent, 0600)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
}
