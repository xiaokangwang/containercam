module github.com/xiaokangwang/containercam

go 1.17

require (
	github.com/docker/distribution v2.7.1+incompatible
	github.com/google/go-containerregistry v0.8.0
	github.com/heroku/docker-registry-client v0.0.0-20211012143308-9463674c8930
	github.com/opencontainers/go-digest v1.0.0
)

require (
	github.com/docker/cli v20.10.12+incompatible // indirect
	github.com/docker/docker v20.10.12+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.6.4 // indirect
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/opencontainers/image-spec v1.0.2-0.20211117181255-693428a734f5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
)

replace github.com/heroku/docker-registry-client v0.0.0-20211012143308-9463674c8930 => github.com/xiaokangwang/ubiquitous-dollop v0.0.0-20220421171237-4bd11a8d62f1
