package object

func NewNamespace() *Namespace {
	s := make(map[string]Object)
	return &Namespace{binds: s, parent: nil}
}

func NewNestedNamespace(ns *Namespace) *Namespace {
	nestedNS := NewNamespace()
	nestedNS.parent = ns
	return nestedNS
}

type Namespace struct {
	binds  map[string]Object
	parent *Namespace
}

func (ns *Namespace) Get(name string) (Object, bool) {
	obj, ok := ns.binds[name]
	if !ok && ns.parent != nil {
		obj, ok = ns.parent.Get(name)
	}
	return obj, ok
}

func (ns *Namespace) Set(name string, val Object) Object {
	ns.binds[name] = val
	return val
}
