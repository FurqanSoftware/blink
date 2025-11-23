package site

type Mark struct {
	Name    string `json:"name"`
	Path    string `json:"path,omitempty"`
	Marks   []Mark `json:"marks,omitempty"`
	SortKey string `json:"sortKey"`
}
