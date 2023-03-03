package main
import(
	"fmt"
	"time"
	"gorm.io/driver/postgres"
  	"gorm.io/gorm"
)
// add database
func create_database(){
	dsn := "host=localhost user=postgres password=postgres dbname=test port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			panic("failed to get underlying sql.DB")
		}
		sqlDB.Close()
	}()
	// migrate the model
	db.AutoMigrate(&SearchResult{})
	
	res, err := BingScrape("Reserve bank Of India", "com", nil, 1, 15, 30)
	if err == nil {
		for _, r := range res {
			// create a new SearchResult with a generated ID and the current time
			sr := SearchResult{
				ResultRank:  r.ResultRank,
				ResultURL:   r.ResultURL,
				ResultTitle: r.ResultTitle,
				ResultDesc:  r.ResultDesc,
				CreatedAt:   time.Now(),
			}
			db.Create(&sr)
			// insert the SearchResult into the database
			if err := db.Create(&sr).Error; err != nil {
				fmt.Println("failed to insert search result:", err)
			}
			fmt.Println(sr)
		}
	} else {
		fmt.Println(err)
	}
	
}