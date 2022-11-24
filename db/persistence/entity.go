package persistence

type CycleState int32

const (
	Init CycleState = iota
	Active
	Destroy
)

type UpdateEntity struct {
	Ver    int
	Entity any
}
