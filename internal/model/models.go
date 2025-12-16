package model

import "time"

type User struct {
	ID        string
	Email     string
	Name      string
	GoogleSub string
	CreatedAt time.Time
}

type Poll struct {
	ID        string
	Question  string
	CreatorID string
	Status    string
	CreatedAt time.Time
	ClosedAt  *time.Time
}

type Option struct {
	ID        string
	PollID    string
	Text      string
	VoteCount int
}
