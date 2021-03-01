package app

import (
	"github.com/harriklein/pBE/pBEServer/utils"
)

const (
	envAppBindAddress = "BE_APP_BIND_ADDRESS"

	defAppBindAddress = ":9090"
)

var (
	cfgAppBindAddress = utils.EnvStr(envAppBindAddress, defAppBindAddress)
)
