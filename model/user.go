package model

/**
 * User(Player) Object
 *
 * For Managing Game User
 *
 * Id: Auto Increase From 1 / 0: ADMIN
 * Salt: Auto Injected - Random Key For Hashing Pass
 * Role: Player Role {admin: super user, basic: basic user}
 * Name: User(Player) Display Name
 * Username: User(Player) ID For LOG-IN
 * Password: User(Player) Pass For LOG-IN
 */
type User struct {
	Id       int
	Salt     string
	Role     string
	Name     string
	Username string
	Password string
	Friend   []int
}
