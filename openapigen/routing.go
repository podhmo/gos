package openapigen

import (
	"github.com/iancoleman/orderedmap"
	"github.com/podhmo/gos/pkg/maplib"
)

type rootRouter struct {
	Actions      []*Action
	DefaultError TypeBuilder
}

type Router struct {
	*rootRouter

	tags []string
}

func NewRouter(defaultError TypeBuilder) *Router {
	return &Router{rootRouter: &rootRouter{
		DefaultError: defaultError,
	}}
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

func (r *Router) ToSchemaWith(b *Builder, doc *orderedmap.OrderedMap) error {
	paths, _ := maplib.GetOrCreate(doc, "paths")

	useRef := !b.Config.DisableRefLinks
	for _, action := range r.Actions {
		if action.metadata.DefaultError == nil {
			action.DefaultError(r.DefaultError)
		}
		pathItem, _ := maplib.GetOrCreate(paths, action.metadata.Path)
		op := action.toSchema(b, useRef)
		pathItem.Set(action.metadata.Method, op)
	}

	if useRef {
		// currently supports components/schemas only

		components, _ := maplib.GetOrCreate(doc, "components")
		schemas, _ := maplib.GetOrCreate(components, "schemas")

		defs := b.Config.defs
		seen := map[int]bool{}
		toplevelDefinitionIsalwaysUseRef := false
		for _, typ := range defs {
			tm := typ.GetTypeMetadata()
			id := tm.id
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = true
			s := typ.toSchema(b, toplevelDefinitionIsalwaysUseRef)
			schemas.Set(tm.Name, s)
		}
	}
	return nil
}
