package domain

import "time"

type CmdType string

type CmdStatus string

const (
	// CmdTypeShell shell command
	CmdTypeShell CmdType = "SHELL"

	// CmdTypeKill kill command
	CmdTypeKill CmdType = "SHELL"

	// CmdTypeClose close command
	CmdTypeClose CmdType = "SHELL"
)

const (
	CmdStatusPending CmdStatus = "PENDING"

	CmdStatusRunning CmdStatus = "RUNNING"

	CmdStatusSuccess CmdStatus = "SUCCESS"

	CmdStatusSkipped CmdStatus = "SKIPPED"

	CmdStatusException CmdStatus = "EXCEPTION"

	CmdStatusKilled CmdStatus = "KILLED"

	CmdStatusTimeout CmdStatus = "TIMEOUT"
)

const (
	// CmdExitCodeUnknow default exit code
	CmdExitCodeUnknow = -1

	// CmdExitCodeTimeOut exit code for timeout
	CmdExitCodeTimeOut = -100

	// CmdExitCodeKilled exit code for killed
	CmdExitCodeKilled = -1

	// CmdExitCodeSuccess exit code for command executed successfuly
	CmdExitCodeSuccess = 0
)

type Cmd struct {
	ID           string `json:"id"`
	AllowFailure bool   `json:"allowFailure"`
	Plugin       string `json:"plugin"`
}

func (cmd *Cmd) HasPlugin() bool {
	return cmd.Plugin != ""
}

type CmdIn struct {
	Cmd
	Type       CmdType   `json:"type"`
	Scripts    []string  `json:"scripts"`
	WorkDir    string    `json:"workDir"`
	Timeout    int64     `json:"timeout"`
	Inputs     Variables `json:"inputs"`
	EnvFilters []string  `json:"envFilters"`
}

type ExecutedCmd struct {
	Cmd
	ProcessId int       `json:"processId"`
	Status    CmdStatus `json:"status"`
	Code      int       `json:"code"`
	Output    Variables `json:"output"`
	StartAt   time.Time `json:"startAt"`
	FinishAt  time.Time `json:"finishAt"`
	Error     string    `json:"error"`
	LogSize   int64     `json:"logSize"`
}
