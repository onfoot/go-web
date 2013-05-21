package handler

import (
	"regexp"
	"net/url"
	"net/http"
	"reflect"
	"sort"
	"fmt"
	"log"
	"bytes"	
)

var h = New()

func New() *RegexpHandler {
	handler := new(RegexpHandler)
	http.Handle("/", handler)
	return handler
}

type route struct {
	url *regexp.Regexp
	handler reflect.Value
}

type Routes []*route

func (r *Routes) Len() int {
	return len(*r)
}

func (r *Routes) Swap(i, j int) {
	routes := *r
	routes[i], routes[j] = routes[j], routes[i]
}

func (r *Routes) Less(i, j int) bool {
	routes := *r
	return len(routes[i].url.String()) > len(routes[j].url.String())
}

func (r *Routes) String() string {
	buffer := bytes.NewBufferString("Current routes: ")

	for _, route := range *r {
		buffer.WriteString(route.url.String())
		buffer.WriteString("\n")
	}

	return buffer.String()
}

type RegexpHandler struct {
	routes Routes
}

func AddRoute(urlPattern string, routeHandler interface{}) {
	h.AddRoute(urlPattern, routeHandler)
}

func (handler *RegexpHandler) AddRoute(urlPattern string, routeHandler interface{}) {

	if funcval, noerror := routeHandler.(reflect.Value); noerror {
		newRoute := &route{regexp.MustCompile(urlPattern), funcval}
		handler.routes = append(handler.routes, newRoute)
	} else {
		funcval := reflect.ValueOf(routeHandler)
		newRoute := &route{regexp.MustCompile(urlPattern), funcval}
		handler.routes = append(handler.routes, newRoute)
	}

	sort.Sort(&handler.routes)
	log.Printf("%v", handler.routes.String())
}

func (handler *RegexpHandler) ServeHTTP(rw http.ResponseWriter, request *http.Request) {

	uri, err := url.QueryUnescape(request.URL.RequestURI())

	if err != nil {
		http.Error(rw, "Malformed URI", http.StatusBadRequest)
		return
	}

	for _, route := range handler.routes {
		matches := route.url.FindStringSubmatch(uri)
		if matches != nil {
			defer func() {
				if err := recover(); err != nil {
					rw.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(rw, "Route handler panicked with error %v", err)
				}
			}()
			
			if route.handler.Type().NumIn() != len(matches ) - 1 + 2 {
				panic("Wrong count of arguments in handler function of " + route.url.String())
			}
			
			var arguments []reflect.Value
			
			arguments = append(arguments, reflect.ValueOf(rw))
			arguments = append(arguments, reflect.ValueOf(request))
			
			for _, match := range matches[1:] {
				arguments = append(arguments, reflect.ValueOf(match))
			}
			
			rw.Header().Set("Content-Type", "text/html")
			route.handler.Call(arguments)
			break
		}
	}
}
