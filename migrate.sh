#!/bin/bash
go get -u github.com/go-i2p/crypto@v0.0.1
go get -u github.com/go-i2p/common@v0.0.1
go get -u github.com/go-i2p/logger@v0.0.1
go get -u github.com/go-i2p/su3@v0.0.1
go get -u github.com/go-i2p/go-noise@v0.0.1
find ./ -name '*.go' -exec sed -i 's|sirupsen/logrus|go-i2p/logger|g' {} \;
find ./ -name '*.go' -exec sed -i 's|logrus|logger|g' {} \;
go mod tidy