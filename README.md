# logger
--
    import "github.com/go-i2p/logger"

Originally written by @hkh4n, modified by @eyedeekay.

Package logger provides a logger for go-i2p.
It is basically a wrapper around logrus that implements "fast-fail" behavior when an environment variable is set.
Right now it only wraps the parts of logrus that we use but if we need more, we can add them.

## Verbosity ##
Logging can be enabled and configured using the `DEBUG_I2P` environment variable. By default, logging is disabled.

There are three available log levels:

- Debug
```shell
export DEBUG_I2P=debug
```
- Warn
```shell
export DEBUG_I2P=warn
```
- Error
```shell
export DEBUG_I2P=error
```

If DEBUG_I2P is set to an unrecognized variable, it will fall back to "debug".

## Fast-Fail mode ##

Fast-Fail mode can be activated by setting `WARNFAIL_I2P` to any non-empty value. When set, every warning or error is Fatal.
It is unsafe for production use, and intended only for debugging and testing purposes.

```shell
export WARNFAIL_I2P=true
```

If `WARNFAIL_I2P` is set and `DEBUG_I2P` is unset, `DEBUG_I2P` will be set to `debug`.

## Usage

#### func  GetGoI2PLogger

```go
func GetGoI2PLogger() *logrus.Logger
```
GetGoI2PLogger returns the initialized logger

#### func  InitializeGoI2PLogger

```go
func InitializeGoI2PLogger()
```
