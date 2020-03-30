package dare

const (
	// EncryptorEngineVersionMin ...
	EncryptorEngineVersionMin = 1
	// EncryptorEngineVersionMax ...
	EncryptorEngineVersionMax = 2
)

// EncryptorEngineState ...
type EncryptorEngineState int32

const (
	// EncryptorEngineRunning ...
	EncryptorEngineRunning EncryptorEngineState = iota
	// EncryptorEngineShutdown ...
	EncryptorEngineShutdown
)

// String ...
func (s EncryptorEngineState) String() string {
	switch s {
	case EncryptorEngineRunning:
		return "Encryptor-engine-running"
	case EncryptorEngineShutdown:
		return "Encryptor-engine-shutdown"
	default:
		return "Encryptor-engine-unknown"
	}
}
