package handlers

// constructing the model from client requests
type SubjectRequests struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Category  string `json:"category,omitempty"`
	Priority  int    `json:"priority,omitempty"`
	Notes     string `json:"notes,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

type SubjectID string

type SubjectResponse []SubjectRequests
