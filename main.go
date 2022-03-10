package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

// Structs to hold the player information
type Player_Data struct {
	FullName string       `json:"fullName"`
	Weight   float64      `json:"weight"`
	Height   float64      `json:"height"`
	Age      int          `json:"age"`
	Years    int          `json:"years"`
	Position string       `json:"position"`
	Stats    Player_Stats `json:"stats"`
	Team     string       `json:"team"`
}

// Struct to hold the stats for each player
type Player_Stats struct {
	ReceivingTouchdowns  int     `json:"receivingTouchdowns"`
	PuntReturnTouchdowns int     `json:"puntReturnTouchdowns"`
	FumbleTouchdown      int     `json:"fumbleTouchdowns"`
	RushingTouchdowns    int     `json:"rushingTouchdowns"`
	PassingTouchdowns    int     `json:"passingTouchdowns"`
	TwoPointRecConvs     float64 `json:"twoPointRecConvs"`
	TwoPointPassConvs    float64 `json:"twoPointPassConvs"`
	TwoPointRushConvs    float64 `json:"twoPointRushConvs"`
	Rushing              float64 `json:"rushing"`
	Receiving            float64 `json:"receiving"`
	Passing              float64 `json:"passing"`
	FumblesLost          float64 `json:"fumblesLost"`
	Interceptions        float64 `json:"interceptions"`
	ExtraPoints          float64 `json:"extraPoints"`
	FieldGoals           float64 `json:"fieldGoals"`
	Sacks                float64 `json:"sacks"`
	Safeties             float64 `json:"safeties"`
	DefensiveTouchdowns  float64 `json:"defensiveTouchdowns"`
	BlockedKicks         float64 `json:"blockedKicks"`
}

type Team struct {
	players []Player_Data
	score   int32
}

// Node to be used in the linked list
type Node struct {
	value Player_Data
	next  *Node
}

// Linked list that will hold every active player in the NFL
type Linked_List struct {
	head  *Node
	tail  *Node
	count int
}

// Method to insert a node into the linked list
func (l *Linked_List) Insert(insertValue Player_Data) {
	newNode := &Node{insertValue, nil}
	if l.head == nil {
		l.head = newNode
		l.tail = newNode
	} else {
		l.tail.next = newNode
		l.tail = newNode
	}
	l.count++
}

func (l *Linked_List) Select(index int) Player_Data {
	temp := l.head
	for i := 0; i < index; i++ {
		temp = temp.next
	}
	return temp.value
}

// Method to print out the linked list
func (l *Linked_List) Print() {
	current_node := l.head
	for current_node != nil {
		// fmt.Println(current_node.value.FullName)
		fmt.Printf("Name: %v, Age: %d, Receiving: %f\n", current_node.value.FullName, current_node.value.Age, current_node.value.Stats.Receiving)
		current_node = current_node.next
	}
}

func collect_player_data(linked_list *Linked_List, refresh_data bool) {
	if _, err := os.Stat("./Saved_Data"); os.IsNotExist(err) {
		os.Mkdir("./Saved_Data", 777)
	}
	// If the file does not exists or you request a refresh, collect the player information
	if _, err := os.Stat("./Saved_Data/PLAYER_STATS_OUTPUT.txt"); err != nil || refresh_data == true {
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
	player_data := Player_Data{}
	for reader.Scan() {
		json.Unmarshal([]byte(reader.Text()), &player_data)
		linked_list.Insert(player_data)
	}
}

func generate_team(athletes Linked_List, return_channel chan []Player_Data) {
	var random_players []int
	// Generate a team of 16 players, they will choose who to bench once the team is generated.
	for i := 0; i < 16; i++ {
		// Select a random player from the list of athletes
		random_athlete := rand.Intn(athletes.count)
		// Check if the athlete has already been selected
		for j := 0; j < len(random_players); j++ {
			if random_athlete == random_players[j] {
				random_athlete = rand.Intn(athletes.count)
				j = 0
			}
		}
		// Add random athlete to the team
		random_players = append(random_players, random_athlete)
	}
	// Grab the player data for the team and return it
	var generated_players []Player_Data
	for _, value := range random_players {
		generated_players = append(generated_players, athletes.Select(value))
	}
	return_channel <- generated_players
	// return generated_players
}

// Simple insertion sort bsed on the score
func insertion_sort(team_list []Team) {
	for i := 1; i < len(team_list); i++ {
		for j := i; j > 0; j-- {
			if team_list[j-1].score < team_list[j].score {
				team_list[j-1], team_list[j] = team_list[j], team_list[j-1]
			}
		}
	}
}

func score_team(team *Team, channel chan bool) {
	for i := 0; i < len(team.players); i++ {
		team.score += int32(team.players[i].Stats.Passing) / 25
		team.score += int32(team.players[i].Stats.Rushing) / 10
		team.score += int32(team.players[i].Stats.Receiving) / 10
		team.score += int32(team.players[i].Stats.RushingTouchdowns) * 6
		team.score += int32(team.players[i].Stats.ReceivingTouchdowns) * 6
		team.score += int32(team.players[i].Stats.PassingTouchdowns) * 4
		team.score -= int32(team.players[i].Stats.FumblesLost) * 2
		team.score -= int32(team.players[i].Stats.Interceptions) * 2
		team.score += int32(team.players[i].Stats.ExtraPoints)
		team.score += int32(team.players[i].Stats.TwoPointRecConvs) * 2
		team.score += int32(team.players[i].Stats.TwoPointRushConvs) * 2
		team.score += int32(team.players[i].Stats.TwoPointPassConvs) * 2
		team.score += int32(team.players[i].Stats.FieldGoals)
		team.score += int32(team.players[i].Stats.Sacks)
		team.score += int32(team.players[i].Stats.Safeties) * 2
		team.score += int32(team.players[i].Stats.DefensiveTouchdowns) * 6
		team.score += int32(team.players[i].Stats.BlockedKicks) * 2
	}
	channel <- true
}

func score_teams(teams []Team) {
	finish_chan := make(chan bool)
	for i := 0; i < len(teams); i++ {
		go score_team(&teams[i], finish_chan)
	}
	for i := 0; i < len(teams); i++ {
		<-finish_chan
	}
	insertion_sort(teams)
}

func crossover(teams []Team, athletes Linked_List) []Team {
	// Keep the top half of the previous population, the rest are killed off
	output := teams[:len(teams)/2]
	// Split teams into fitness
	most_fit := teams[:int(float64(len(teams))*0.5)]
	well_fit := teams[int(float64(len(teams))*0.5):int(float64(len(teams))*0.8)]
	worst_fit := teams[int(float64(len(teams))*0.8):]
	// Regnerate the population with the offspring
	var team1 Team
	var team2 Team
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
		start_index := rand.Intn(len(team1.players) - int(float64(len(team1.players))*.5))
		end_index := rand.Intn(len(team1.players)-start_index) + start_index
		var new_team []Player_Data
		new_team = append(new_team, team1.players[start_index:end_index]...)
		for _, team2_value := range team2.players {
			is_in := false
			for _, player := range new_team {
				if team2_value.FullName == player.FullName {
					is_in = true
					break
				}
			}
			if !is_in {
				new_team = append(new_team, team2_value)
				if len(new_team) == len(team1.players) {
					break
				}
			}
		}
		output = append(output, Team{new_team, 0})
	}
	return output
}

func mutate(team []Team, athletes Linked_List) {
	// Mutate 5% of the population
	mutate_number := int(float64(len(team)) * .05)
	for i := 0; i < mutate_number; i++ {
		// Get a random player on a random team and replace him with a random athlete
		random_team := rand.Intn(len(team))
		team_member := rand.Intn(len(team[random_team].players))
		new_team_member := rand.Intn(athletes.count)
		team[random_team].players[team_member] = athletes.Select(new_team_member)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	// Parse flags from command line
	refresh_data := flag.Bool("refresh_data", false, "Determines if you want to refresh player data.")
	generation_count := flag.Int("generation_count", 10000, "Number of generations that the GA will run for.")
	population_size := flag.Int("population_size", 1000, "Size of the population to be used by the GA.")
	flag.Parse()
	var player_list Linked_List
	// Collect player data from either the text file or have the python script collect the data
	collect_player_data(&player_list, *refresh_data)
	// Generate the teams
	team_channel := make(chan []Player_Data)
	number_of_teams := *population_size
	for i := 0; i < number_of_teams; i++ {
		go generate_team(player_list, team_channel)
	}
	var team_list []Team
	for i := 0; i < number_of_teams; i++ {
		fmt.Printf("Generated Team %d/%d\r", i+1, number_of_teams)
		team_list = append(team_list, Team{<-team_channel, 0})
	}
	fmt.Printf("\n")

	// Loop until generation cap is reached
	for i := 0; i < *generation_count; i++ {
		// Score the teams then sort
		score_teams(team_list)
		// Get the new population from the crossover function
		team_list = crossover(team_list, player_list)
		// Mutate some of the population
		mutate(team_list, player_list)
		// Print progress report
		fmt.Printf("Generating your team... %v%v\r", fmt.Sprintf("%d", int(float64(i)/float64(*generation_count)*100)), "%")
	}

	fmt.Println("\nTeam has been generated. Your players are:")
	for _, player := range team_list[0].players {
		fmt.Println(player.FullName, " from the ", player.Team)
	}
}
