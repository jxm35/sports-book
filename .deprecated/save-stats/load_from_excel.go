package save_stats

// deprecated
//
//func loadTeams(file string, includeIds bool) []*model.Team {
//	var teams []*model.Team
//	f, err := excelize.OpenFile(file)
//	if err != nil {
//		fmt.Println(err)
//		return nil
//	}
//	rows, err := f.GetRows("Teams")
//	for i, row := range rows {
//		if i != 0 {
//			name := row[1]
//			if includeIds {
//				id, _ := strconv.Atoi(row[0])
//				teams = append(teams, &model.Team{
//					Name: name,
//					ID:   int32(id),
//				})
//			} else {
//				teams = append(teams, &model.Team{Name: name})
//			}
//		}
//	}
//	return teams
//}
//
//func loadAppearances(file string, db *gorm.DB) []*model.PlayedIn {
//	var apps []*model.PlayedIn
//	f, err := excelize.OpenFile(file)
//	if err != nil {
//		fmt.Println(err)
//		return nil
//	}
//	matches := loadMatches(file, db, true)
//	players := loadPlayers(file, true)
//	teams := loadTeams(file, true)
//	rows, err := f.GetRows("Appearances2")
//	for i, row := range rows {
//		if i != 0 {
//			matchString := row[0]
//			match, err := lookupMatch(matchString, db, matches)
//			if err != nil {
//				continue
//			}
//			teamString := row[1]
//			team, err := lookupTeam(teamString, db, teams)
//			if err != nil {
//				continue
//			}
//			xgString := row[2]
//			xG, err := strconv.ParseFloat(xgString, 32)
//			if err != nil || xG == 0 {
//				xG = -1
//			}
//			xAString := row[3]
//			xA, err := strconv.ParseFloat(xAString, 32)
//			if err != nil || xA == 0 {
//				xA = -1
//			}
//			goalString := row[4]
//			goals, err := strconv.Atoi(goalString)
//			if err != nil {
//				continue
//			}
//			assistString := row[5]
//			assists, err := strconv.Atoi(assistString)
//			if err != nil {
//				continue
//			}
//			minutesString := row[6]
//			minutes, err := strconv.Atoi(minutesString)
//			if err != nil {
//				continue
//			}
//			playerString := row[7]
//			player, err := lookupPlayer(playerString, db, players)
//			if err != nil {
//				continue
//			}
//
//			app := model.PlayedIn{
//				Player:          player,
//				Match:           match,
//				Team:            team,
//				ExpectedGoals:   xG,
//				ActualGoals:     int32(goals),
//				Minutes:         int32(minutes),
//				ExpectedAssists: xA,
//				ActualAssists:   float64(assists),
//			}
//			apps = append(apps, &app)
//		}
//	}
//	return apps
//}
//
//func loadMatches(file string, db *gorm.DB, includeId bool) []*model.Match {
//	var matches []*model.Match
//	f, err := excelize.OpenFile(file)
//	if err != nil {
//		fmt.Println(err)
//		return nil
//	}
//
//	teams := loadTeams(file, true)
//	rows, err := f.GetRows("Matches")
//	for i, row := range rows {
//		if i != 0 {
//			dateString := row[1]
//			dateInt, err := strconv.Atoi(dateString)
//			if err != nil {
//				continue
//			}
//			date := time.Unix(int64(dateInt), 0)
//
//			homeTeamString := row[2]
//			homeTeam, err := lookupTeam(homeTeamString, db, teams)
//			if err != nil {
//				continue
//			}
//			awayTeamString := row[3]
//			awayTeam, err := lookupTeam(awayTeamString, db, teams)
//			if err != nil {
//				continue
//			}
//
//			homexGString := row[4]
//			homexG, err := strconv.ParseFloat(homexGString, 32)
//			if err != nil {
//				homexG = -1
//			}
//			awayxGString := row[5]
//			awayxG, err := strconv.ParseFloat(awayxGString, 32)
//			if err != nil {
//				awayxG = -1
//			}
//
//			homeGoalsString := row[6]
//			homeGoals, err := strconv.Atoi(homeGoalsString)
//			if err != nil {
//				continue
//			}
//			awayGoalsString := row[7]
//			if err != nil {
//				continue
//			}
//			awayGoals, err := strconv.Atoi(awayGoalsString)
//
//			homexA := float64(-1)
//			awayxA := float64(-1)
//			if len(row) >= 10 {
//				homexAString := row[8]
//				homexA, err = strconv.ParseFloat(homexAString, 32)
//				if err != nil {
//					homexA = -1
//				}
//				awayxAString := row[9]
//				awayxA, err = strconv.ParseFloat(awayxAString, 32)
//				if err != nil {
//					awayxA = -1
//				}
//			}
//			if !includeId {
//				matches = append(matches, &model.Match{
//					Date:                date,
//					HomeTeam:            homeTeam,
//					AwayTeam:            awayTeam,
//					Competition:         1, // prem
//					HomeExpectedGoals:   homexG,
//					AwayExpectedGoals:   awayxG,
//					HomeGoals:           int32(homeGoals),
//					AwayGoals:           int32(awayGoals),
//					HomeExpectedAssists: homexA,
//					AwayExpectedAssists: awayxA,
//				})
//			} else {
//				id, _ := strconv.Atoi(row[0])
//				matches = append(matches, &model.Match{
//					ID:                  int32(id),
//					Date:                date,
//					HomeTeam:            homeTeam,
//					AwayTeam:            awayTeam,
//					Competition:         1, // prem
//					HomeExpectedGoals:   homexG,
//					AwayExpectedGoals:   awayxG,
//					HomeGoals:           int32(homeGoals),
//					AwayGoals:           int32(awayGoals),
//					HomeExpectedAssists: homexA,
//					AwayExpectedAssists: awayxA,
//				})
//			}
//		}
//	}
//	return matches
//}
//
//func loadPlayers(file string, includeId bool) []*model.Player {
//	var players []*model.Player
//	f, err := excelize.OpenFile(file)
//	if err != nil {
//		fmt.Println(err)
//		return nil
//	}
//	rows, err := f.GetRows("Players")
//	for i, row := range rows {
//		if i != 0 {
//			name := row[1]
//			position := row[2]
//			if !includeId {
//				players = append(players, &model.Player{Name: name, Position: position})
//			} else {
//				id, _ := strconv.Atoi(row[0])
//				players = append(players, &model.Player{
//					Name:     name,
//					Position: position,
//					ID:       int32(id),
//				})
//			}
//
//		}
//	}
//	return players
//}
