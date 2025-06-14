package interceptor

// TODO: Structure it better with user friendly messages
var errors = map[string]string{
	"BPB001": "Bad Request",
	"BPB002": "Invalid Email",
	"BPB003": "Invalid Password",
	"BPB004": "Password cannot be same as email",
	"BPB005": "Unable to register user",
	"BPB006": "Unable to authenticate user",
	"BPB007": "User not found",
	"BPB008": "Invalid username or password",
	"BPB009": "Unable to generate token",
	"BPB010": "Refresh token is required",
	"BPB011": "Unable to process refresh token",
	"BPB012": "Invalid or expired refresh token",
	"BPB500": "Internal Server Error",
}
