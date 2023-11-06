package Structs

import (
	"testing"
)

func TestInsert(t *testing.T) {
	tests := []struct {
		name        string
		linked_List *Linked_List
		insertValue []Player_Data
		Count       int
	}{
		{
			name: "Insert into empty list",
			linked_List: &Linked_List{
				Head:  nil,
				Tail:  nil,
				Count: 0,
			},
			insertValue: []Player_Data{
				{
					FullName: "Test",
				},
			},
			Count: 1,
		},
		{
			name: "Insert into list with one value",
			linked_List: &Linked_List{
				Head:  nil,
				Tail:  nil,
				Count: 0,
			},
			insertValue: []Player_Data{
				{
					FullName: "Test2",
				},
				{
					FullName: "Test3",
				},
			},
			Count: 2,
		},
		{
			name: "Insert into list with multiple values",
			linked_List: &Linked_List{
				Head:  nil,
				Tail:  nil,
				Count: 0,
			},
			insertValue: []Player_Data{
				{
					FullName: "Test3",
				},
				{
					FullName: "Test4",
				},
				{
					FullName: "Test5",
				},
			},
			Count: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, value := range test.insertValue {
				test.linked_List.Insert(value)
			}

			if test.linked_List.Count != test.Count {
				t.Errorf("Insert() did not insert the correct number of values")
				t.Fail()
				return
			}
		})
	}
}
