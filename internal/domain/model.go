package domain

type Household struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
}

type User struct {
	ID          int64  `json:"id"`
	HouseholdID int64  `json:"householdId"`
	MemberID    *int64 `json:"memberId"`
	Name        string `json:"name"`
	Email       string `json:"email"`
}

type Member struct {
	HouseholdID int64  `json:"householdId"`
	ID          int64  `json:"id"`
	Name        string `json:"name"`
}
