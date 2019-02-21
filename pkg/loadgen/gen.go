package loadgen

type LoadGenerator interface {
	Start()
	Status() string
}
