package model

/**
 * Game Model Object
 *
 * TODO: Add Game Mode Option
 */
type Game struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	MinPlayer int    `json:"minPlayer"`
}
