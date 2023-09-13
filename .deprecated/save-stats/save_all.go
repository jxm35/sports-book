package save_stats

// deprecated
//
//func SaveAll() {
//	gormDb, err := gorm.Open(mysql.Open("root:password@tcp(127.0.0.1:3306)/sports-book?charset=utf8mb4&parseTime=True&loc=Local"))
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	files := []string{
//		"stats/prem21-22-gw1-26.xlsx",
//		"stats/prem21-22-gw26-38.xlsx",
//		"stats/prem22-23-gw1-29.xlsx",
//		"stats/prem22-23gw29-38.xlsx",
//		"stats/prem20-21-gw1-9.xlsx",
//		"stats/prem20-21-gw9-38.xlsx",
//	}
//	for _, file := range files {
//		//teams := loadTeams(file, false)
//		//saveTeams(gormDb, teams)
//		//
//		//matches := loadMatches(file, gormDb, false)
//		//saveMatches(gormDb, matches)
//		//
//		//players := loadPlayers(file, false)
//		//savePlayers(gormDb, players)
//
//		//apps := loadAppearances(file, gormDb)
//		//saveAppearances(gormDb, apps)
//	}
//}
