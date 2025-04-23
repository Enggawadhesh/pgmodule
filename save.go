package savejsonfile

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(Savejsonfile{})
	httpcaddyfile.RegisterHandlerDirective("savejsonfile", parseCaddyfile)
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m Savejsonfile
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return m, err
}

type Savejsonfile struct {
	FilePath string `json:"file_path,omitempty"`
}

func (Savejsonfile) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.savejsonfile",
		New: func() caddy.Module { return new(Savejsonfile) },
	}
}

func (s Savejsonfile) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	file, err := os.Create(s.FilePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		return fmt.Errorf("failed to write request body to file: %w", err)
	}

	return next.ServeHTTP(w, r)
}

func (s *Savejsonfile) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if !d.NextArg() {
			return d.ArgErr()
		}
		s.FilePath = d.Val()
	}
	return nil
}

var (
	_ caddy.Module                = (*Savejsonfile)(nil)
	_ caddyhttp.MiddlewareHandler = (*Savejsonfile)(nil)
	_ caddyfile.Unmarshaler       = (*Savejsonfile)(nil)
)
