package openapigen

import (
	"github.com/iancoleman/orderedmap"
)

type rootRouter struct {
	Actions []*Action
}

type Router struct {
	*rootRouter
	tags []string
}

func NewRouter() *Router {
	return &Router{rootRouter: &rootRouter{}}
}

func (r *Router) Tagged(tags ...string) *Router {
	return &Router{tags: tags, rootRouter: r.rootRouter}
}

func (r *Router) Method(method string, path string, action *Action) {
	action = action.Method(method).Path(path)
	if r.tags != nil {
		action = action.Tags(append(action.metadata.Tags, r.tags...))
	}
	r.Actions = append(r.Actions, action)
}
func (r *Router) Get(path string, action *Action) {
	r.Method("get", path, action)
}
func (r *Router) Post(path string, action *Action) {
	r.Method("post", path, action)
}
func (r *Router) Put(path string, action *Action) {
	r.Method("put", path, action)
}
func (r *Router) Patch(path string, action *Action) {
	r.Method("patch", path, action)
}
func (r *Router) Delete(path string, action *Action) {
	r.Method("delete", path, action)
}

// hmm
func (r *Router) ToSchemaWith(b *Builder, doc *orderedmap.OrderedMap) error {
	var paths *orderedmap.OrderedMap
	v, ok := doc.Get("paths")
	if !ok {
		paths = orderedmap.New()
		doc.Set("paths", paths)
	} else {
		paths = v.(*orderedmap.OrderedMap)
	}

	useRef := !b.Config.DisableRefLinks
	for _, action := range r.Actions {
		op := action.toSchema(b, useRef)
		var pathItem *orderedmap.OrderedMap
		v, ok := paths.Get(action.metadata.Path)
		if !ok {
			pathItem = orderedmap.New()
			paths.Set(action.metadata.Path, pathItem)
		} else {
			pathItem = v.(*orderedmap.OrderedMap)
		}
		pathItem.Set(action.metadata.Name, op)
	}
	return nil
}
