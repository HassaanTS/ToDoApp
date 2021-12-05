package todos

type ToDo struct {
	Id      string `json:"_id,omitempty" bson:"_id,omitempty"`
	Title   string `json:"title" bson:"title"`
	Desc    string `json:"desc" bson:"desc"`
	DueDate string `json:"date" bson:"date"`
	Done    bool   `json:"done" bson:"done"`
}

func New() ToDo {
	return ToDo{}
}
