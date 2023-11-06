package GAF

import (
	"fantasy_football/Collect_Stats"
	"fantasy_football/Structs"
	"fmt"
	"os"
	"testing"
)

func TestGenerate_Team(t *testing.T) {
	os.Chdir("..")
	var player_list Structs.Linked_List
	Collect_Stats.Collect_Player_Data(&player_list, false)

	tests := []struct {
		name     string
		athletes *Structs.Linked_List
		want     chan []Structs.Player_Data
		willFail bool
	}{
		{
			name:     "Valid Athletes List",
			athletes: &player_list,
			want:     make(chan []Structs.Player_Data),
			willFail: false,
		},
		{
			name: "Empty Athletes List",
			athletes: &Structs.Linked_List{
				Count: 0,
				Head:  nil,
				Tail:  nil,
			},
			want:     make(chan []Structs.Player_Data),
			willFail: true,
		},
		{
			name:     "Nil Athletes List",
			athletes: nil,
			want:     make(chan []Structs.Player_Data),
			willFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			go Generate_Team(test.athletes, test.want)
			output := <-test.want

			if len(output) == 0 && !test.willFail {
				t.Errorf("Generate_Team() did not return a valid team")
				t.Fail()
				return
			}
		})
	}
}

func TestScore_team(t *testing.T) {
	tests := []struct {
		name string
		team []Structs.Team
		want int32
	}{
		{
			name: "Passing Score >25",
			team: []Structs.Team{
				{
					Players: []Structs.Player_Data{
						{
							Stats: Structs.Player_Stats{
								Passing: 1,
							},
						},
					},
				},
			},
			want: 0,
		},
		{
			name: "Passing Score 25",
			team: []Structs.Team{
				{
					Players: []Structs.Player_Data{
						{
							Stats: Structs.Player_Stats{
								Passing: 25,
							},
						},
					},
				},
			},
			want: 1,
		},
		{
			name: "Rushing Score >10",
			team: []Structs.Team{
				{
					Players: []Structs.Player_Data{
						{
							Stats: Structs.Player_Stats{
								Rushing: 1,
							},
						},
					},
				},
			},
		},
		{
			name: "Rushing Score 10",
			team: []Structs.Team{
				{
					Players: []Structs.Player_Data{
						{
							Stats: Structs.Player_Stats{
								Rushing: 10,
							},
						},
					},
				},
			},
			want: 1,
		},
		{
			name: "Receiving Score >10",
			team: []Structs.Team{
				{
					Players: []Structs.Player_Data{
						{
							Stats: Structs.Player_Stats{
								Receiving: 1,
							},
						},
					},
				},
			},
		},
		{
			name: "Receiving Score 10",
			team: []Structs.Team{
				{
					Players: []Structs.Player_Data{
						{
							Stats: Structs.Player_Stats{
								Receiving: 10,
							},
						},
					},
				},
			},
			want: 1,
		},
		{
			name: "RushingTouchdowns Score",
			team: []Structs.Team{
				{
					Players: []Structs.Player_Data{
						{
							Stats: Structs.Player_Stats{
								RushingTouchdowns: 1,
							},
						},
					},
				},
			},
			want: 6,
		},
		{
			name: "ReceivingTouchdowns Score",
			team: []Structs.Team{
				{
					Players: []Structs.Player_Data{
						{
							Stats: Structs.Player_Stats{
								ReceivingTouchdowns: 1,
							},
						},
					},
				},
			},
			want: 6,
		},
		{
			name: "PassingTouchdowns Score",
			team: []Structs.Team{
				{
					Players: []Structs.Player_Data{
						{
							Stats: Structs.Player_Stats{
								PassingTouchdowns: 1,
							},
						},
					},
				},
			},
			want: 4,
		},
		{
			name: "FumblesLost Score",
			team: []Structs.Team{
				{
					Players: []Structs.Player_Data{
						{
							Stats: Structs.Player_Stats{
								FumblesLost: 1,
							},
						},
					},
				},
			},
			want: -2,
		},
		{
			name: "Interceptions Score",
			team: []Structs.Team{
				{
					Players: []Structs.Player_Data{
						{
							Stats: Structs.Player_Stats{
								Interceptions: 1,
							},
						},
					},
				},
			},
			want: -2,
		},
		{
			name: "ExtraPoints Score",
			team: []Structs.Team{
				{
					Players: []Structs.Player_Data{
						{
							Stats: Structs.Player_Stats{
								ExtraPoints: 1,
							},
						},
					},
				},
			},
			want: 1,
		},
		{
			name: "TwoPointRecConvs Score",
			team: []Structs.Team{
				{
					Players: []Structs.Player_Data{
						{
							Stats: Structs.Player_Stats{
								TwoPointRecConvs: 1,
							},
						},
					},
				},
			},
			want: 2,
		},
		{
			name: "TwoPointRushConvs Score",
			team: []Structs.Team{
				{
					Players: []Structs.Player_Data{
						{
							Stats: Structs.Player_Stats{
								TwoPointRushConvs: 1,
							},
						},
					},
				},
			},
			want: 2,
		},
		{
			name: "TwoPointPassConvs Score",
			team: []Structs.Team{
				{
					Players: []Structs.Player_Data{
						{
							Stats: Structs.Player_Stats{
								TwoPointPassConvs: 1,
							},
						},
					},
				},
			},
			want: 2,
		},
		{
			name: "FieldGoals Score",
			team: []Structs.Team{
				{
					Players: []Structs.Player_Data{
						{
							Stats: Structs.Player_Stats{
								FieldGoals: 1,
							},
						},
					},
				},
			},
			want: 1,
		},
		{
			name: "Sacks Score",
			team: []Structs.Team{
				{
					Players: []Structs.Player_Data{
						{
							Stats: Structs.Player_Stats{
								Sacks: 1,
							},
						},
					},
				},
			},
			want: 1,
		},
		{
			name: "Safeties Score",
			team: []Structs.Team{
				{
					Players: []Structs.Player_Data{
						{
							Stats: Structs.Player_Stats{
								Safeties: 1,
							},
						},
					},
				},
			},
			want: 2,
		},
		{
			name: "DefensiveTouchdowns Score",
			team: []Structs.Team{
				{
					Players: []Structs.Player_Data{
						{
							Stats: Structs.Player_Stats{
								DefensiveTouchdowns: 1,
							},
						},
					},
				},
			},
			want: 6,
		},
		{
			name: "BlockedKicks Score",
			team: []Structs.Team{
				{
					Players: []Structs.Player_Data{
						{
							Stats: Structs.Player_Stats{
								BlockedKicks: 1,
							},
						},
					},
				},
			},
			want: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			score_team(&test.team[0])
			fmt.Println(test.team[0].Score, test.want)
			if test.team[0].Score != test.want {
				t.Errorf("score_team() did not score the team")
				t.Fail()
				return
			}
		})
	}
}

func TestInsertion_Sort(t *testing.T) {
	tests := []struct {
		name string
		team []Structs.Team
		want []Structs.Team
	}{
		{
			name: "Sorted List",
			team: []Structs.Team{
				{
					Score: 1,
				},
				{
					Score: 2,
				},
				{
					Score: 3,
				},
			},
			want: []Structs.Team{
				{
					Score: 3,
				},
				{
					Score: 2,
				},
				{
					Score: 1,
				},
			},
		},
		{
			name: "Unsorted List",
			team: []Structs.Team{
				{
					Score: 3,
				},
				{
					Score: 1,
				},
				{
					Score: 2,
				},
			},
			want: []Structs.Team{
				{
					Score: 3,
				},
				{
					Score: 2,
				},
				{
					Score: 1,
				},
			},
		},
		{
			name: "Empty List",
			team: []Structs.Team{},
			want: []Structs.Team{},
		},
		{
			name: "Single Element List",
			team: []Structs.Team{
				{
					Score: 1,
				},
			},
			want: []Structs.Team{
				{
					Score: 1,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			insertion_sort(test.team)
			for i := 0; i < len(test.team); i++ {
				if test.team[i].Score != test.want[i].Score {
					t.Errorf("insertion_sort() did not sort the list")
					t.Fail()
					return
				}
			}
		})
	}
}

/*
	No point in testing Score_teams because it is just a wrapper for score_team and insertion_sort
*/
