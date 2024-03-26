package www

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jtacoma/uritemplates"
)

const protomaps_api string = "https://api.protomaps.com/tiles/v3/{z}/{x}/{y}.mvt?key={key}"

type TilesAPIHandlerOptions struct {
	ProtomapsApiKey string
}

func TilesAPIHandler(opts *TilesAPIHandlerOptions) (http.Handler, error) {

	t, err := uritemplates.Parse(protomaps_api)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse protomaps URI template, %w", err)
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		// ctx := req.Context()

		logger := slog.Default()
		logger = logger.With("request", req.URL)

		z := req.PathValue("z")
		x := req.PathValue("x")
		y := req.PathValue("y")

		values := map[string]interface{}{
			"z":   z,
			"x":   x,
			"y":   y,
			"key": opts.ProtomapsApiKey,
		}

		tile_url, err := t.Expand(values)

		if err != nil {
			slog.Error("Failed to expand template", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(rsp, req, tile_url, http.StatusSeeOther)
		return
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
