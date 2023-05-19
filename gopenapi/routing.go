package gopenapi

type Router struct {
	Actions []*Action
	tags    []string
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Tagged(tags ...string) *Router {
	return &Router{tags: tags, Actions: r.Actions}
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
func (r *Router) Post(method string, path string, action *Action) {
	r.Method("post", path, action)
}
func (r *Router) Put(method string, path string, action *Action) {
	r.Method("put", path, action)
}
func (r *Router) Patch(method string, path string, action *Action) {
	r.Method("patch", path, action)
}
func (r *Router) Delete(method string, path string, action *Action) {
	r.Method("delete", path, action)
}
