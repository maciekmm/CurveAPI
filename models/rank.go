package models

type Rank struct {
	Rank  int `json:"rank" bson:"rank"`
	Bonus int `json:"bonus" bson:"bonus"`
	IsNew int `json:"isNew" bson:"isNew"`
	//TODO: Find a way to deserialize as string and serialize as int
	Pos   int `json:"pos,string" bson:"pos"`
}
