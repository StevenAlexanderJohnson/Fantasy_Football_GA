package structs

import "fmt"

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
	Players []Player_Data
	Score   int32
}

// Node to be used in the linked list
type Node struct {
	Value Player_Data
	Next  *Node
}

// Linked list that will hold every active player in the NFL
type Linked_List struct {
	Head  *Node
	Tail  *Node
	Count int
}

// Method to insert a node into the linked list
func (l *Linked_List) Insert(insertValue Player_Data) {
	newNode := &Node{insertValue, nil}
	if l.Head == nil {
		l.Head = newNode
		l.Tail = newNode
	} else {
		l.Tail.Next = newNode
		l.Tail = newNode
	}
	l.Count++
}

func (l *Linked_List) Select(index int) Player_Data {
	temp := l.Head
	for i := 0; i < index; i++ {
		temp = temp.Next
	}
	return temp.Value
}

// Method to print out the linked list
func (l *Linked_List) Print() {
	current_node := l.Head
	for current_node != nil {
		// fmt.Println(current_node.value.FullName)
		fmt.Printf("Name: %v, Age: %d, Receiving: %f\n", current_node.Value.FullName, current_node.Value.Age, current_node.Value.Stats.Receiving)
		current_node = current_node.Next
	}
}
