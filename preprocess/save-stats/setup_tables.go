package save_stats

//
//func SetupTables() {
//	a := query.PlayedIn
//	m := query.Match
//	ctx := context.Background()
//	// update xg to not have -1 values
//	info, err := a.WithContext(ctx).
//		Where(a.ExpectedAssists.Eq(-1), a.ExpectedGoals.Eq(-1)).
//		Update(a.ExpectedAssists, nil)
//	fmt.Println(info.Error, err)
//
//	info, err = a.WithContext(ctx).
//		Where(a.ExpectedAssists.IsNull(), a.ExpectedGoals.Eq(-1)).
//		Update(a.ExpectedGoals, nil)
//	fmt.Println(info.Error, err)
//
//	info, err = m.WithContext(ctx).
//		Where(
//			m.HomeExpectedGoals.Eq(-1),
//			m.AwayExpectedGoals.Eq(-1),
//			m.HomeExpectedAssists.Eq(-1),
//			m.AwayExpectedAssists.Eq(-1)).
//		Update(m.HomeExpectedGoals, nil)
//	fmt.Println(info.Error, err)
//
//	info, err = m.WithContext(ctx).
//		Where(
//			m.HomeExpectedGoals.IsNull(),
//			m.AwayExpectedGoals.Eq(-1),
//			m.HomeExpectedAssists.Eq(-1),
//			m.AwayExpectedAssists.Eq(-1)).
//		Update(m.AwayExpectedGoals, nil)
//	fmt.Println(info.Error, err)
//
//	info, err = m.WithContext(ctx).
//		Where(
//			m.HomeExpectedGoals.IsNull(),
//			m.AwayExpectedGoals.IsNull(),
//			m.HomeExpectedAssists.Eq(-1),
//			m.AwayExpectedAssists.Eq(-1)).
//		Update(m.HomeExpectedAssists, nil)
//	fmt.Println(info.Error, err)
//
//	info, err = m.WithContext(ctx).
//		Where(
//			m.HomeExpectedGoals.IsNull(),
//			m.AwayExpectedGoals.IsNull(),
//			m.HomeExpectedAssists.IsNull(),
//			m.AwayExpectedAssists.Eq(-1)).
//		Update(m.AwayExpectedAssists, nil)
//	fmt.Println(info.Error, err)
//}
