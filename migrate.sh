#!/bin/bash
go get -u github.com/go-i2p/crypto@0d868ff0e313a787d83890a42cd17a68b212f4ef
go get -u github.com/go-i2p/common@1812f37de0377c85e5fe9beb854d78b6fc161051
go get -u github.com/go-i2p/logger@b7cf9a3377d987790cf451385b07f28993ed8c9a
find ./ -name '*.go' -exec sed -i 's|sirupsen/logrus|go-i2p/logger|g' {} \;
find ./ -name '*.go' -exec sed -i 's|logrus|logger|g' {} \;