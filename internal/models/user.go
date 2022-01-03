package models

type User struct {
	Nickname string `json:"nickname" valid:"type(string),minstringlength(1)"`
	Fullname string `json:"fullname" valid:"type(string),minstringlength(1)"`
	About    string `json:"about" valid:"type(string),minstringlength(0)"`
	Email    string `json:"email" valid:"email"`
}

type UserUpdate struct {
	Fullname string `json:"fullname" valid:"type(string),minstringlength(1)"`
	About    string `json:"about" valid:"type(string),minstringlength(0)"`
	Email    string `json:"email" valid:"email"`
}

func (ua *UserUpdate) ToUser(nickname string) *User {
	return &User{
		Nickname: nickname,
		Fullname: ua.Fullname,
		About:    ua.About,
		Email:    ua.Email,
	}
}
