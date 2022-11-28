package terminfindung

func VoteForDay(terminID *string, userAdmin *string, name *string, votedDay *string, voted int8) {
	termin, err := GetTerminFromShared(userAdmin, terminID)
	if err != nil {
		return
	}
	if voted > 0 {
		termin.Persons[*name].Votes[*votedDay] = true
	}
	if voted < 0 {
		termin.Persons[*name].Votes[*votedDay] = false
	}

}
