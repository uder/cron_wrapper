package model

import (
	"gorm.io/gorm"
)

type Table interface {
	TableName() string
}

type CommandRecords struct {
	gorm.Model
	RecordType string
	RunId      uint64 `gorm:"uniqueIndex"`
	Command    string
	Timeout    int
}

func (CommandRecords) TableName() string {
	return "command_records"
}

func NewCommandRecord(recordType string, runId uint64, cmd string, timeout int) *CommandRecords {
	return &CommandRecords{
		RunId:      runId,
		RecordType: recordType,
		Command:    cmd,
		Timeout:    timeout,
	}
}

func Models() []Table {
	return []Table{
		&CommandRecords{},
	}
}
