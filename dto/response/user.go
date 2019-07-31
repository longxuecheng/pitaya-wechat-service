package response

type User struct {
	AvatarURL string `json:"avatarUrl"`
	NickName  string `json:"nickName"`
	UserType  int    `json:"userType"`
}
