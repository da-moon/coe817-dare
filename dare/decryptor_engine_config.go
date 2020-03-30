package dare

const (
	// DecryptorEngineVersionMin ...
	DecryptorEngineVersionMin = 1
	// DecryptorEngineVersionMax ...
	DecryptorEngineVersionMax = 2
)

// DecryptorEngineState ...
type DecryptorEngineState int32

const (
	// DecryptorEngineRunning ...
	DecryptorEngineRunning DecryptorEngineState = iota
	// DecryptorEngineShutdown ...
	DecryptorEngineShutdown
)

// String ...
func (s DecryptorEngineState) String() string {
	switch s {
	case DecryptorEngineRunning:
		return "Decryptor-engine-running"
	case DecryptorEngineShutdown:
		return "Decryptor-engine-shutdown"
	default:
		return "Decryptor-engine-unknown"
	}
}
