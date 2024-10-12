package database

import (
	"main/database/model"
	"time"
)

func GetUnfinishedRecords(db DB, commandLine string, timeout int) []model.CommandRecords {
	var unfinishedRecords []model.CommandRecords
	subQuery := db.Conn().
		Table("command_records").
		Select("run_id").
		Where("record_type IN ?", []string{"INFO", "ERROR"})
	db.Conn().
		Where("command = ?", commandLine).
		Where("record_type = ?", "BEGIN").
		Where("run_id NOT IN (?)", subQuery).
		Where("created_at >= ?", time.Now().Add(-2*time.Duration(timeout)*time.Second)).
		Find(&unfinishedRecords)
	return unfinishedRecords
}
