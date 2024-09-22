package main

func Verify_User(username string, password string) bool {
	if Retrieve_User(username).Password == password {
		return true
	} else {
		return false
	}
}

func Get_User() {

}
