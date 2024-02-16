package server

import (
	"context"
	"flag"
	"fmt"
	html_template "html/template"
	io_fs "io/fs"

	"github.com/aaronland/go-http-server/handler"
	"github.com/mitchellh/copystructure"
	"github.com/sfomuseum/go-flags/flagset"
	sfom_funcs "github.com/sfomuseum/go-template/funcs"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd/static"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd/templates/html"
)

type RunOptions struct {
	ServerURI         string                              `json:"server_uri"`
	SpelunkerURI      string                              `json:"spelunker_uri"`
	AuthenticatorURI  string                              `json:"authenticator_uri"`
	URIs              *httpd.URIs                         `json:"uris"`
	HTMLTemplates     []io_fs.FS                          `json:"templates,omitemtpy"`
	HTMLTemplateFuncs html_template.FuncMap               `json:"template_funcs,omitempty"`
	StaticAssets      io_fs.FS                            `json:"static_assets,omitempty"`
	CustomHandlers    map[string]handler.RouteHandlerFunc `json:"custom_handlers,omitempty"`
}

func (o *RunOptions) Clone() (*RunOptions, error) {

	v, err := copystructure.Copy(o)

	if err != nil {
		return nil, fmt.Errorf("Failed to create local run options, %w", err)
	}

	new_opts := v.(*RunOptions)

	// Things that aren't handled by copystructure
	// TBD...

	new_opts.HTMLTemplates = o.HTMLTemplates
	new_opts.HTMLTemplateFuncs = o.HTMLTemplateFuncs
	new_opts.StaticAssets = o.StaticAssets

	return new_opts, nil
}

func RunOptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "SPELUNKER")

	if err != nil {
		return nil, fmt.Errorf("Failed to assign flags from environment variables, %w", err)
	}

	uris_table = httpd.DefaultURIs()

	t_funcs := html_template.FuncMap{
		"IsAvailable": sfom_funcs.IsAvailable,
		// "Add":              sfom_funcs.Add,
		"JoinPath": sfom_funcs.JoinPath,
		// "QRCodeB64":        sfom_funcs.QRCodeB64,
		// "QRCodeDataURI":    sfom_funcs.QRCodeDataURI,
		// "IsEven":           sfom_funcs.IsEven,
		// "IsOdd":            sfom_funcs.IsOdd,
		"FormatStringTime": sfom_funcs.FormatStringTime,
		"FormatUnixTime":   sfom_funcs.FormatUnixTime,
		"GjsonGet":         sfom_funcs.GjsonGet,
	}

	opts := &RunOptions{
		ServerURI:         server_uri,
		AuthenticatorURI:  authenticator_uri,
		SpelunkerURI:      spelunker_uri,
		URIs:              uris_table,
		HTMLTemplates:     []io_fs.FS{html.FS},
		HTMLTemplateFuncs: t_funcs,
		StaticAssets:      static.FS,
	}

	return opts, nil
}
