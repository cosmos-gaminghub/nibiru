package services

import (
	"github.com/cosmos-gaminghub/nibiru/pkg/xfilepath"
)

var (
	// StarportConfPath returns the Starport Configuration directory
	StarportConfPath = xfilepath.JoinFromHome(xfilepath.Path(".starport"))
)
