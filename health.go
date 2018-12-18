package gounity

// Health defines Unity corresponding `health` type.
type Health struct {
	Value          int      `json:"value"`
	DescriptionIds []string `json:"descriptionIds"`
	Descriptions   []string `json:"descriptions"`
}
