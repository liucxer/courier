package gengo

var registeredGenerators = map[string]Generator{}

func GetRegisteredGenerators(names ...string) (generators []Generator) {
	if len(names) == 0 {
		for name := range registeredGenerators {
			generators = append(generators, registeredGenerators[name])
		}
		return
	}

	for _, name := range names {
		if _, ok := registeredGenerators[name]; ok {
			generators = append(generators, registeredGenerators[name])
		}
	}

	return
}

func Register(g Generator) {
	registeredGenerators[g.Name()] = g
}
