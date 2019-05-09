package types

type (
	Film struct {
		ID          string              `json:"id" bson:"_id"`
		Rank        int                 `json:"rank" bson:"rank"`
		URL         string              `json:"url" bson:"url"`
		Title       string              `json:"title" bson:"title"`
		Rate        string              `json:"rate" bson:"rate"`
		ReleaseDate int                 `json:"release_date" bson:"release_date"`
		Description string              `json:"description" bson:"description"`
		Credit      map[string][]string `json:"credit" bson:"credit"`
	}
)
