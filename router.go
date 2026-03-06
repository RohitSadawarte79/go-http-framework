package main

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/RohitSadawarte79/go-http-framework/internal/domain"
)

type Router struct {
	routes []route
}

func NewRouter() *Router {
	return &Router{
		routes: make([]route, 0),
	}
}

type route struct {
	method   string
	pattern  string
	segments []segment
	handler  http.HandlerFunc
}

type segment struct {
	value   string
	isParam bool
}

func conflictsWithRoute(a, b []segment) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].isParam && b[i].isParam {
			continue
		}

		if a[i].isParam != b[i].isParam {
			return false
		}

		if a[i].value != b[i].value {
			return false
		}
	}

	return true
}

func countParams(segments []segment) int {
	var count int
	for _, segment := range segments {
		if segment.isParam {
			count++
		}
	}

	return count
}

func isParamValid(segments []segment) bool {
	for _, segment := range segments {
		if segment.isParam {
			if segment.value == "" {
				return false
			}
		}
	}

	return true
}

func (rt *Router) HandleFunc(method, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if method == "" {
		panic("router: method must not be empty")
	}
	if pattern == "" {
		panic("router: pattern must not be empty")
	}
	if handler == nil {
		panic("router: handler must not be empty")
	}

	pattern = strings.TrimRight(pattern, "/")
	segments := parsePath(pattern)
	if !isParamValid(segments) {
		panic("router: param name must not be empty in pattern registered it carrefully")
	}
	CurrRoute := route{
		method:   method,
		pattern:  pattern,
		segments: segments,
		handler:  http.HandlerFunc(handler),
	}

	for _, registeredRoutes := range rt.routes {
		if registeredRoutes.method != CurrRoute.method {
			continue
		}
		segmentsRegistered := registeredRoutes.segments
		if conflictsWithRoute(segmentsRegistered, segments) {
			panic("Conflicts between routes: " + registeredRoutes.pattern + " and " + CurrRoute.pattern)
		}
	}

	rt.routes = append(rt.routes, CurrRoute)

	slices.SortFunc(rt.routes, func(a, b route) int {
		return countParams(a.segments) - countParams(b.segments)
	})
}

/*
for "/user/:id" parsePath will return list of segements will be [{value:"user", isParam: false}, {value:"id", isParam: true}]
for "/user/:name" parsePath will return list of segements will be [{value:"user", isParam: false}, {value:"name", isParam: true}]
after this both will be registered and if the request arrives like
hostname/user/42 then it will be parsed by parsePath and then segments will be [{value:"user", isParam: false}, {value"42" , isParam: false}]
after that the loop will search in routes then it will found this path id path was registered first so it will math and then the method will match
and then handler assosiated with it will be called whenever this form of request comes the match will always be of :id param because it was registered first.
*/

func URLParam(r *http.Request, name string) string {
	params, ok := r.Context().Value(domain.ParamKey).(map[string]string)

	if !ok {
		return ""
	}

	return params[name]
}

func parsePath(pattern string) []segment {
	// first split the string with /
	// iterate the split string while creating segment and checking if
	// param or not pursh into the segment slice
	// return the segment slice

	res := make([]segment, 0)

	patternSplit := strings.SplitSeq(pattern, "/")

	for segStr := range patternSplit {
		if segStr == "" {
			continue
		}

		segment := segment{
			value:   segStr,
			isParam: false,
		}

		if segStr[:1] == ":" {
			segment.value = segStr[1:]
			segment.isParam = true
		}

		res = append(res, segment)
	}

	return res
}

func matchSegments(registered []segment, requested []segment) bool {
	if len(registered) != len(requested) {
		return false
	}

	for i, regSeg := range registered {
		if regSeg.isParam {
			continue
		}

		if regSeg.value != requested[i].value {
			return false
		}
	}

	return true
}

func parseParams(registered []segment, request []segment) map[string]string {
	params := make(map[string]string)

	for i, segment := range registered {
		if segment.isParam {
			params[segment.value] = request[i].value
		}
	}

	return params
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requested := parsePath(r.URL.Path)
	methodNotAllowed := false

	for _, route := range rt.routes {
		if matchSegments(route.segments, requested) {
			if route.method == r.Method {
				params := parseParams(route.segments, requested)

				ctx := context.WithValue(r.Context(), domain.ParamKey, params)

				r = r.WithContext(ctx)

				route.handler.ServeHTTP(w, r)
				return
			}

			methodNotAllowed = true
		}
	}

	if methodNotAllowed {
		http.Error(w, "405 Method Not allowed", http.StatusMethodNotAllowed)
	} else {
		http.Error(w, "404: Page Not found", http.StatusNotFound)
	}
}
