package dto

type NavBarOption struct {
	DisplayName string
	Email       string
}

type DropdownIdentifier string

const (
	DropdownIdentifierNone        DropdownIdentifier = "none"
	DropdownIdentifierTeam        DropdownIdentifier = "team"
	DropdownIdentifierProject     DropdownIdentifier = "project"
	DropdownIdentifierEnvironment DropdownIdentifier = "environment"
)

type NavItem struct {
	Title              string
	Icon               string
	NavURL             string
	DropdownIdentifier DropdownIdentifier
	DropdownActionURL  string
}
