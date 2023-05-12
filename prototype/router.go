package prototype

import (
	"strings"

	"github.com/iancoleman/orderedmap"
)

type Route struct {
	Method string
	Path   string
	Action *ActionType

	Tag string
}

type Router struct {
	b      *Builder
	routes []*Route
}

func (r *Router) ToSchema() (*orderedmap.OrderedMap, error) {
	d := orderedmap.New()
	paths := orderedmap.New()
	d.Set("paths", paths)

	for _, route := range r.routes {
		var path *orderedmap.OrderedMap
		if v, ok := d.Get(route.Path); ok {
			path = v.(*orderedmap.OrderedMap)
		} else {
			path = orderedmap.New()
			paths.Set(route.Path, path)
		}
		path.Set(strings.ToLower(route.Method), route.Action.toSchema(r.b))
	}

	// TODO: tree shaking
	return ToSchemaWith(r.b, d)
}

func (r *Router) Method(method string, path string, action *ActionType) {
	r.methodWithTag(method, path, action, "")
}
func (r *Router) methodWithTag(method string, path string, action *ActionType, tag string) {
	r.routes = append(r.routes, &Route{Method: method, Path: path, Action: action, Tag: tag})
}
func (r *Router) Group(name string, fn func(*Grouped)) {
	g := &Grouped{r: r, Name: name}
	fn(g)
}

func NewRouter(b *Builder) *Router {
	return &Router{b: b}
}

type Grouped struct {
	r    *Router
	Name string
}

func (g *Grouped) Method(method string, path string, action *ActionType) {
	g.r.methodWithTag(method, path, action, g.Name)
}
