package tasks

import (
	"net/http"

	utils "github.com/nelsliu9121/wgwheretogoserver/utils"
)

var client = &http.Client{}
var fb = utils.Firebase()
