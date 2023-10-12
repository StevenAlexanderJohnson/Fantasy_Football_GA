package main

import (
	"bufio"
	"encoding/json"
	GAF "fantasy_football/GeneticAlgorithmFunctions"
	Structs "fantasy_football/Structs"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
)

/*
This function will collect the player data from the text file if it exists, otherwise it will run the python script to collect the data.
Once the data is collected it will be read into a linked list of player data structs.

Parameters:

	linked_list: A pointer to the linked list that will be used to store the player data.
	refresh_data: A boolean that determines if you want to refresh the player data.

Returns:

	None
*/
func collect_player_data(linked_list *Structs.Linked_List, refresh_data bool) {
	if _, err := os.Stat("./Saved_Data"); os.IsNotExist(err) {
		os.Mkdir("./Saved_Data", 0777)
	}
	// If the file does not exists or you request a refresh, collect the player information
	if _, err := os.Stat("./Saved_Data/PLAYER_STATS_OUTPUT.txt"); err != nil || refresh_data {
		fmt.Println("Collecting player data, this may take a few minutes...")
		c := exec.Command("py", "./Collect_Stats/main.py")
		_, err := c.Output()
		if err != nil {
			panic(err)
		}
		fmt.Println("Player data has been collected...")
	}
	fmt.Println("Reading player data...")
	// Open the player stats file and defer its closure
	player_stats_file, err := os.OpenFile("./Saved_Data/PLAYER_STATS_OUTPUT.txt", os.O_RDONLY, fs.FileMode(os.O_RDONLY))
	if err != nil {
		panic(err)
	}
	defer player_stats_file.Close()
	// Create a scanner to read the file, read line by line and parse each player json into the structs
	reader := bufio.NewScanner(player_stats_file)
	player_data := Structs.Player_Data{}
	for reader.Scan() {
		json.Unmarshal([]byte(reader.Text()), &player_data)
		linked_list.Insert(player_data)
	}
}

func main() {
	// Parse flags from command line
	refresh_data := flag.Bool("refresh_data", false, "Determines if you want to refresh player data.")
	generation_count := flag.Int("generation_count", 10000, "Number of generations that the GA will run for.")
	population_size := flag.Int("population_size", 1000, "Size of the population to be used by the GA.")
	flag.Parse()
	var player_list Structs.Linked_List
	// Collect player data from either the text file or have the python script collect the data
	collect_player_data(&player_list, *refresh_data)
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
