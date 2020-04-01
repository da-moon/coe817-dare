package dare

const (
	// DareEngineVersionMin ...
	DareEngineVersionMin = 1
	// DareEngineVersionMax ...
	DareEngineVersionMax = 2
)

// DareEngineState ...
type DareEngineState int32

const (
	// DareEngineRunning ...
	DareEngineRunning DareEngineState = iota
	// DareEngineShutdown ...
	DareEngineShutdown
)

// String ...
func (s DareEngineState) String() string {
	switch s {
	case DareEngineRunning:
		return "dare-engine-running"
	case DareEngineShutdown:
		return "dare-engine-shutdown"
	default:
		return "dare-engine-unknown"
	}
}
