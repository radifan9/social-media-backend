package models

// Success
// {
//   "success": true,
//   "status": 200,
//   "data": { "id": 1, "name": "John" }
// }

// Error
// {
//   "success": false,
//   "status": 400,
//   "message": "Invalid email format"
// }

type Response struct {
	Success bool   `json:"success"`
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type SuccessResponse struct {
	Success bool `json:"success" example:"true"`
	Status  int  `json:"status" example:"200"`
	Data    any  `json:"data"`
}

type SuccessLoginResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Status  int    `json:"status" example:"500"`
	Error   string `json:"error" example:"error message"`
}
