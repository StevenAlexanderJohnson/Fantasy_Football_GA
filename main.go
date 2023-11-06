package main

import (
	"fantasy_football/Collect_Stats"
	GAF "fantasy_football/GeneticAlgorithmFunctions"
	"fantasy_football/Structs"
	"flag"
	"fmt"
)

func main() {
	// Parse flags from command line
	refresh_data := flag.Bool("refresh_data", false, "Determines if you want to refresh player data.")
	generation_count := flag.Int("generation_count", 10000, "Number of generations that the GA will run for.")
	population_size := flag.Int("population_size", 1000, "Size of the population to be used by the GA.")
	flag.Parse()
	var player_list Structs.Linked_List
	// Collect player data from either the text file or have the python script collect the data
	Collect_Stats.Collect_Player_Data(&player_list, *refresh_data)
	// Generate the teams
	team_channel := make(chan []Structs.Player_Data)
	number_of_teams := *population_size
	for i := 0; i < number_of_teams; i++ {
		go GAF.Generate_Team(&player_list, team_channel)
	}
	var team_list []Structs.Team
	for i := 0; i < number_of_teams; i++ {
		fmt.Printf("Generated Team %d/%d\r", i+1, number_of_teams)
		team_list = append(team_list, Structs.Team{Players: <-team_channel, Score: 0})
	}
	fmt.Printf("\n")

	// Loop until generation cap is reached
	for i := 0; i < *generation_count; i++ {
		// Score the teams then sort
		GAF.Score_teams(team_list)
		// Get the new population from the crossover function
		team_list = GAF.Crossover(team_list)
		// Mutate some of the population
		GAF.Mutate(team_list, &player_list)
		// Print progress report
		fmt.Printf("Generating your team... %v%v\r", fmt.Sprintf("%d", int(float64(i)/float64(*generation_count)*100)), "%")
	}

	fmt.Println("\nTeam has been generated. Your Players are:")
	for _, player := range team_list[0].Players {
		fmt.Println(player.FullName, " from the ", player.Team)
	}
	fmt.Println("With a score of: ", team_list[0].Score)
}
