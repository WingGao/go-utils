package wsql

import "time"

type SqlModel struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt *time.Time
}
