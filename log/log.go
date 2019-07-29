package log

type Config struct {
	Family string
	Host   string

	//stdout
	Stdout bool

	// file
	Dir string
	//buffer size
	FileBufferSize int64
	//MaxLogFile
	MaxLogFile int
	// RotateSize
	RotateSize int64

	//log-agent
	//Agent *Agent
}
