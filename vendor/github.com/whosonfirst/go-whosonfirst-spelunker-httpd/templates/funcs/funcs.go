package funcs

import (
	"github.com/whosonfirst/go-whosonfirst-sources"
)

func NameForSource(prefix string) string {

	src, err := sources.GetSourceByName(prefix)

	if err != nil {
		return prefix
	}

	return src.Fullname
}
