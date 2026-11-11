package slinx

type SyncRequest struct {
	User []ReportUser `json:"user"`
}

type ReportUser struct {
	ID       int      `json:"id"`
	Upload   int64    `json:"u"`
	Download int64    `json:"d"`
	IP       []string `json:"ip"`
}

type SyncResponse struct {
	User []UserResponse `json:"user"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	UUID     string `json:"uuid"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
