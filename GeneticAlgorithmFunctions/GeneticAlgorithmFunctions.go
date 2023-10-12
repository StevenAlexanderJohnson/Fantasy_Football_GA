package GAF

import (
	Structs "fantasy_football/Structs"
	"math/rand"
)

/*
Generate a team to compete in the fantasy football league

Parameters:

	athletes: The list of athletes to choose from.
	return_channel: The channel to return the team to.

Returns:

	An array of players that make up the team.
*/
func Generate_Team(athletes *Structs.Linked_List, return_channel chan []Structs.Player_Data) {
	var random_Players []int
	// Generate a team of 16 Players, they will choose who to bench once the team is generated.
	for i := 0; i < 16; i++ {
		// Select a random player from the list of athletes
		random_athlete := rand.Intn(athletes.Count)
		// Check if the athlete has already been selected
		for j := 0; j < len(random_Players); j++ {
			if random_athlete == random_Players[j] {
				random_athlete = rand.Intn(athletes.Count)
				j = 0
			}
		}
		// Add random athlete to the team
		random_Players = append(random_Players, random_athlete)
	}
	// Grab the player data for the team and return it
	var generated_Players []Structs.Player_Data
	for _, value := range random_Players {
		generated_Players = append(generated_Players, athletes.Select(value))
	}
	return_channel <- generated_Players
	// return generated_Players
}

/*
Score a team based on the players' stats.

Parameters:

	team: The team to be scored.

Returns:

	Nothing
*/
func score_team(team *Structs.Team) {
	team.Score = 0
	for i := 0; i < len(team.Players); i++ {
		team.Score += int32(team.Players[i].Stats.Passing) / 25
		team.Score += int32(team.Players[i].Stats.Rushing) / 10
		team.Score += int32(team.Players[i].Stats.Receiving) / 10
		team.Score += int32(team.Players[i].Stats.RushingTouchdowns) * 6
		team.Score += int32(team.Players[i].Stats.ReceivingTouchdowns) * 6
		team.Score += int32(team.Players[i].Stats.PassingTouchdowns) * 4
		team.Score -= int32(team.Players[i].Stats.FumblesLost) * 2
		team.Score -= int32(team.Players[i].Stats.Interceptions) * 2
		team.Score += int32(team.Players[i].Stats.ExtraPoints)
		team.Score += int32(team.Players[i].Stats.TwoPointRecConvs) * 2
		team.Score += int32(team.Players[i].Stats.TwoPointRushConvs) * 2
		team.Score += int32(team.Players[i].Stats.TwoPointPassConvs) * 2
		team.Score += int32(team.Players[i].Stats.FieldGoals)
		team.Score += int32(team.Players[i].Stats.Sacks)
		team.Score += int32(team.Players[i].Stats.Safeties) * 2
		team.Score += int32(team.Players[i].Stats.DefensiveTouchdowns) * 6
		team.Score += int32(team.Players[i].Stats.BlockedKicks) * 2
	}
}

// Simple insertion sort bsed on the Score
func insertion_sort(team_list []Structs.Team) {
	for i := 1; i < len(team_list); i++ {
		for j := i; j > 0; j-- {
			if team_list[j-1].Score < team_list[j].Score {
				team_list[j-1], team_list[j] = team_list[j], team_list[j-1]
			}
		}
	}
}

/*
Score all of the teams in the population

Parameters:

	teams: The teams to be scored.

Returns:

	Nothing
*/
func Score_teams(teams []Structs.Team) {
	for i := 0; i < len(teams); i++ {
		go score_team(&teams[i])
	}
	insertion_sort(teams)
}

/*
Takes the population and simulates trading players between teams

Parameters:

	teams: The teams to crossover.

Returns:

	The new list of teams.
*/
func Crossover(teams []Structs.Team) []Structs.Team {
	// Keep the top half of the previous population, the rest are killed off
	output := teams[:len(teams)/2]
	// Split teams into fitness
	most_fit := teams[:int(float64(len(teams))*0.5)]
	well_fit := teams[int(float64(len(teams))*0.5):int(float64(len(teams))*0.8)]
	worst_fit := teams[int(float64(len(teams))*0.8):]
	// Regnerate the population with the offspring
	var team1 Structs.Team
	var team2 Structs.Team
	team2_select := false
	for len(output) < len(teams) {
		for i := 0; i < 2; i++ {
			// Generate a random number to determine who will be the parent
			random_team := rand.Intn(10)
			if random_team < 5 {
				if !team2_select {
					team1 = most_fit[rand.Intn(len(most_fit))]
					team2_select = true
					continue
				}
				team2 = most_fit[rand.Intn(len(most_fit))]
			} else if random_team < 8 {
				if !team2_select {
					team1 = well_fit[rand.Intn(len(well_fit))]
					team2_select = true
					continue
				}
				team2 = well_fit[rand.Intn(len(well_fit))]
			} else {
				if !team2_select {
					team1 = worst_fit[rand.Intn(len(worst_fit))]
					team2_select = true
					continue
				}
				team2 = worst_fit[rand.Intn(len(worst_fit))]
			}
		}

		// Mix the genetics of the two parents and add them to the team
		start_index := rand.Intn(len(team1.Players) - int(float64(len(team1.Players))*.5))
		end_index := rand.Intn(len(team1.Players)-start_index) + start_index
		var new_team []Structs.Player_Data
		new_team = append(new_team, team1.Players[start_index:end_index]...)
		for _, team2_value := range team2.Players {
			is_in := false
			for _, player := range new_team {
				if team2_value.FullName == player.FullName {
					is_in = true
					break
				}
			}
			if !is_in {
				new_team = append(new_team, team2_value)
				if len(new_team) == len(team1.Players) {
					break
				}
			}
		}
		output = append(output, Structs.Team{Players: new_team, Score: 0})
	}
	return output
}

/*
Mutates a member of the population by swapping a random player with a random athlete

Parameters:

	team: The team to be mutated.
	athletes: The list of athletes to choose from.

Returns:

	Nothing
*/
func Mutate(team []Structs.Team, athletes *Structs.Linked_List) {
	// Mutate 5% of the population
	mutate_number := int(float64(len(team)) * .05)
	for i := 0; i < mutate_number; i++ {
		// Get a random player on a random team and replace him with a random athlete
		random_team := rand.Intn(len(team))
		team_member := rand.Intn(len(team[random_team].Players))
		new_team_member := rand.Intn(athletes.Count)
		team[random_team].Players[team_member] = athletes.Select(new_team_member)
	}
}
