package Collect_Stats

import (
	"bufio"
	"encoding/json"
	"fantasy_football/Structs"
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
func Collect_Player_Data(linked_list *Structs.Linked_List, refresh_data bool) {
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
