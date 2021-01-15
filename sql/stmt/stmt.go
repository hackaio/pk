package stmt

//name
//username
//email
//password
//created

const (
	DELETE = "DELETE FROM accounts WHERE name = ? AND username = ?;"
	ADD    = "INSERT INTO accounts(name,username,email,password, created) VALUES(?, ?, ?, ?, ?);"
	GET    = "SELECT * FROM accounts WHERE name = ? AND username = ?;"
	LIST   = "SELECT * FROM accounts;"
	UPDATE = ""
)
