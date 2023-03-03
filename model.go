package main
import(
	"time"
)
// Add structs
type SearchResult struct {
    ResultRank  int    `gorm:"primary_key"`
    ResultURL   string `gorm:"not null"`
    ResultTitle string `gorm:"not null"`
    ResultDesc  string `gorm:"not null"`
    CreatedAt   time.Time `gorm:"not null"`
}