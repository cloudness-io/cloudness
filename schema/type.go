package schema

import _ "embed"

//go:embed application.json
var Application []byte

//go:embed template.json
var Template []byte
