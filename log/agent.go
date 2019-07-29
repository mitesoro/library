package log

type AgentConfig struct {
	TaskID string
	Buffer int
	Proto  string `dsn:"network"`
	Addr   string `dsn:"address"`
	Chan   string `dsn:"query.chan"`
	//Timeout
}
