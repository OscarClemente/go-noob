package models

type Review struct {
	ID        int    `json:"id"`
	Game      string `json:"game"`
	Title     string `json:"title"`
	Content   string `json:"description"`
	Rating    int    `json:"rating"`
	UserID    int    `json:"user"`
	CreatedAt string `json:"created_at"`
}

type ReviewList struct {
	Reviews []Review `json:"reviews"`
}
