package service

func Login(username, password string) (string, error) {
	if username == "" || password == "" {
		return "", nil
	}
	return "demo-token", nil
}

func VerifyToken(token string) (bool, string) {
	if token == "demo-token" {
		return true, "student"
	}
	return false, ""
}
