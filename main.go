package main

import (
	"githubEvents/api"
	"githubEvents/collector"
	"githubEvents/shared"
	_ "githubEvents/shared"
)

func main() {
	if shared.Role == shared.ApiRoleName {
		api.Main()
	} else {
		collector.Main()
	}
}
