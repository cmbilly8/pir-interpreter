package object

func NewNamespace() *Namespace {
	s := make(map[string]Object)
	return &Namespace{store: s}
}

type Namespace struct {
	store map[string]Object
}

func (ns *Namespace) Get(name string) (Object, bool) {
	obj, ok := ns.store[name]
	return obj, ok
}
func (ns *Namespace) Set(name string, val Object) Object {
	ns.store[name] = val
	return val
}
