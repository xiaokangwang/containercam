package containercam

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/opencontainers/go-digest"
	"io"
)

func NewStorageBackend(
	registry *registry.Registry,
	registryName string) *Storage {
	return &Storage{registryName: registryName, registry: registry, keepaliveName: identity}
}

func NewStorageBackendWithCustomNameFunc(
	registry *registry.Registry,
	registryName string,
	keepaliveName func(sha256value string) string) *Storage {
	return &Storage{registryName: registryName, registry: registry, keepaliveName: keepaliveName}
}

type Storage struct {
	registry      *registry.Registry
	registryName  string
	keepaliveName func(sha256value string) string
}

func (s *Storage) DownloadByHash(sha256Value string) ([]byte, error) {
	data, err := s.registry.DownloadBlob(s.registryName, digest.Digest("sha256:"+sha256Value))
	if err != nil {
		return nil, err
	}
	hashedData, err := io.ReadAll(data)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(hashedData)
	if hex.EncodeToString(hash[:]) != sha256Value {
		return nil, errors.New("invalid content from registry")
	}
	return hashedData, nil
}

func (s *Storage) UploadByHash(sha256Value string, hashedData []byte) error {
	hash := sha256.Sum256(hashedData)
	if hex.EncodeToString(hash[:]) != sha256Value {
		return errors.New("invalid content from registry")
	}
	dataDigest := digest.FromBytes(hashedData)
	cached, err := s.registry.HasBlob(s.registryName, dataDigest)
	if err != nil {
		return err
	}
	if !cached {
		err = s.registry.UploadBlob(s.registryName, dataDigest, bytes.NewReader(hashedData))
		if err != nil {
			return err
		}
	}
	keepAliveManifest := schema2.Manifest{
		Versioned: schema2.SchemaVersion,
		Config:    distribution.Descriptor{},
		Layers: []distribution.Descriptor{
			{
				MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip+encrypted",
				Size:      int64(len(hashedData)),
				Digest:    dataDigest,
			},
		},
	}
	keepAliveManifestBlob, err := schema2.FromStruct(keepAliveManifest)
	if err != nil {
		return err
	}
	err = s.registry.PutManifest(s.registryName, s.keepaliveName(sha256Value), keepAliveManifestBlob)
	return err
}

func identity(s string) string {
	return s
}
