package response

// Error is the uniform error envelope returned for every failure.
type Error struct {
	Message string `json:"error"`
}
