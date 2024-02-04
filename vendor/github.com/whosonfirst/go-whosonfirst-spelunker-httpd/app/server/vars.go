package server

import (
	"sync"

	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
)

var sp spelunker.Spelunker

var uris_table *httpd.URIs

var setupCommonOnce sync.Once
var setupCommonError error
