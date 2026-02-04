package dto

type NavBarOption struct {
	DisplayName string
	Email       string
}

type NavItem struct {
	Title              string
	Icon               string
	NavURL             string
	DropdownIdentifier string
	DropdownActionURL  string
}
