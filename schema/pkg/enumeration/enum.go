package enumeration

type WithConstValues interface {
	ConstValues() []IntStringerEnum
}

type IntStringerEnum interface {
	WithConstValues

	TypeName() string
	Int() int
	String() string
	Label() string
}

// DriverValueOffset
// sql value maybe have offset from const value in go
type DriverValueOffset interface {
	Offset() int
}
