package embeded

import "embed"

//go:embed views/*
var ViewFs embed.FS
