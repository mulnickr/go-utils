package serve

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

// AllMethods is an array of all methods which will be used if no methods are specified.
var AllMethods = []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace}

type route struct {
	method   string
	segments []string
	handler  Handler
}

func (r *route) match(c context.Context, rq *http.Request, us []string) (context.Context, bool) {
	if len(us) != len(r.segments) {
		return c, false
	}

	for i, rs := range r.segments {
		if i > len(us)-1 {
			return c, false
		}

		if key, ok := strings.CutPrefix(rs, ":"); ok {
			uv, err := url.QueryUnescape(us[i])
			if err != nil || uv == "" {
				return c, false
			}

			rq.SetPathValue(key, uv)
			continue
		}

		if us[i] != rs {
			return c, false
		}
	}

	return c, true
}
