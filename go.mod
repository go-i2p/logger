module github.com/go-i2p/logger

go 1.26.3

require github.com/sirupsen/logrus v1.10.2

require (
	github.com/go-i2p/elgamal v0.1.69-0.20260908193841-095c65b90f3e
	github.com/go-i2p/su3 v0.1.69-0.20260908193243-800e506feb49
	golang.org/x/sys v0.48.0 // indirect
)

//replace github.com/go-i2p/go-i2p => /home/idk/go/src/github.com/go-i2p/go-i2p

//replace github.com/go-i2p/go-i2p => /home/idk/go/src/github.com/go-i2p/go-i2p

//replace github.com/go-i2p/go-i2p => /home/idk/go/src/github.com/go-i2p/go-i2p

//replace github.com/go-i2p/go-i2p => /home/idk/go/src/github.com/go-i2p/go-i2p

retract (
	v0.1.59999
	v0.1.5999
	v0.1.599
)
