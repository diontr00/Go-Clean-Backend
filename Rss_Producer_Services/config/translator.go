package config

import (
	"embed"
	"khanhanhtr/sample2/translator"
)

//go:embed trans_file/*.toml
var trans_folder embed.FS

func newTranslator() (*translator.UniversalTrans, error) {
	return translator.NewUtTrans(trans_folder, "trans_file")
}
