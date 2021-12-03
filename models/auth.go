package models

// Login struct
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginInfo struct {
	ID       string `json:"id" bson:"_id"`
	UserType string `json:"user_type" bson:"user_type"`
	FullName string `json:"full_name" bson:"full_name"`
	Password string `json:"password" bson:"password"`
	Login    string `json:"login" bson:"login"`
	Soato    string `json:"soato" bson:"soato"`
}

type LoginExistsRequest struct {
	Login string `json:"login"`
}
type LoginExistsResponse struct {
	Exist bool `json:"exist"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}
