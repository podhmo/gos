package prototype

type Route struct {
	Method string
	Path   string
	Action *ActionType
}

type Router struct {
	routes []*Route
}

func (r *Router) Method(method string, path string, action *ActionType) {
	r.routes = append(r.routes, &Route{Method: method, Path: path, Action: action})
}
