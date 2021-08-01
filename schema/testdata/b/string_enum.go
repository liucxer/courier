package b

type PullPolicy string

const (
	PullAlways       PullPolicy = "Always"       // pull always
	PullNever        PullPolicy = "Never"        // never
	PullIfNotPresent PullPolicy = "IfNotPresent" // if not preset
)
