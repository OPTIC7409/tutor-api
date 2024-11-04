package models

type DashboardData struct {
	ID               int               `json:"id"`
	Name             string            `json:"name"`
	Email            string            `json:"email"`
	UserType         string            `json:"userType"`
	Chats            []Chat            `json:"chats"`
	Requests         []Request         `json:"requests,omitempty"`
	UpcomingSessions []UpcomingSession `json:"upcomingSessions,omitempty"`
	Stats            *TutorStats       `json:"stats,omitempty"`
}

type Request struct {
	ID      int    `json:"id"`
	Student string `json:"student"`
	Subject string `json:"subject"`
	Budget  string `json:"budget"`
	Avatar  string `json:"avatar"`
}

type UpcomingSession struct {
	Tutor    string `json:"tutor"`
	Subject  string `json:"subject"`
	Datetime string `json:"datetime"`
}

type TutorStats struct {
	ActiveStudents    int     `json:"activeStudents"`
	UpcomingSessions  int     `json:"upcomingSessions"`
	EarningsThisMonth float64 `json:"earningsThisMonth"`
}
