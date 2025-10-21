#!/bin/bash

find / -name '*.go' -exec sed -i 's|sirupsen/logrus|go-i2p/logger|g' {} \;
find / -name '*.go' -exec sed -i 's|logrus|logger|g' {} \;