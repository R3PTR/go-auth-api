package sites

type Workspace struct {
	Id   string `bson:"_id,omitempty"`
	Name string `json:"name"`
}

type Site struct {
	Id   string `bson:"_id,omitempty"`
	Name string `json:"name"`
}
