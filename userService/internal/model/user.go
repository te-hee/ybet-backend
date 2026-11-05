package model

import "github.com/google/uuid"


type User struct{
	Id uuid.UUID;
	Username string;
	Password string; // password should be salted and encrypted
}

func NewUser(id uuid.UUID, username string, password string) (User){
	return User{Id: id, Username: username, Password: password} 
}
