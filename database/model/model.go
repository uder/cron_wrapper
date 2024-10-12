package model

import (
	"database/sql"
	"gorm.io/gorm"
)

type Table interface {
	TableName() string
}

type CommandRecords struct {
	gorm.Model
	RecordType string `gorm:"uniqueIndex:idx_runId_type"`
	RunId      string `gorm:"uniqueIndex:idx_runId_type"`
	Command    string
	ExitCode   sql.NullInt16   `gorm:"default:null"`
	Duration   sql.NullFloat64 `gorm:"default:null"`
}

func (CommandRecords) TableName() string {
	return "command_records"
}

func NewCommandBeginRecord(recordType string, runId string, cmd string) *CommandRecords {
	return &CommandRecords{
		RunId:      runId,
		RecordType: recordType,
		Command:    cmd,
	}
}

func NewCommandEndRecord(recordType string, runId string, cmd string, exitCode int, duration float64) *CommandRecords {
	return &CommandRecords{
		RunId:      runId,
		RecordType: recordType,
		Command:    cmd,
		ExitCode:   sql.NullInt16{Int16: int16(exitCode), Valid: true},
		Duration:   sql.NullFloat64{Float64: duration, Valid: true},
	}
}

func Models() []Table {
	return []Table{
		&CommandRecords{},
	}
}
