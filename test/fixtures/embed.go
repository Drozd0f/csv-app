package fixtures

import "embed"

//go:embed *.csv
var Fixtures embed.FS
