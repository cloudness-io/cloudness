package check

import (
	"fmt"
	"regexp"
)

const (
	minDisplayNameLength = 3
	maxDisplayNameLength = 25

	minIdentifierLength = 3
	MaxIdentifierLength = 25
	identifierRegex     = "^[a-z][a-z0-9-.]*[a-z0-9]$"
	gitRepoHttpsPattern = `^https:\/\/([\w.-]+)\/([\w-]+)\/([\w.-]+)\.git$`
	gitRepoSshPattern   = `^git@([\w.-]+):([\w-]+)\/([\w.-]+)\.git$`
	fqdnPattern         = `^(https?://)([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`
	ipV4Pattern         = `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	directoryPattern    = `^((/[a-zA-Z0-9-_]+)+|/)$`
	emailRegex          = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	//variables
	minVarKey   = 3
	maxVarKey   = 50
	varKeyRegex = `^[a-zA-Z][a-zA-Z0-9_]*[a-zA-Z0-9]$`

	minPasswordLength = 5
	maxPasswordLength = 64

	maxDescriptionLength = 1024
)

var (
	ErrIdentifierLength = &ValidationError{
		fmt.Sprintf(
			"Identifier has to be between %d and %d in length.",
			minIdentifierLength,
			MaxIdentifierLength,
		),
	}

	ErrDescriptionTooLong = &ValidationError{
		fmt.Sprintf("Description can be at most %d in length.", maxDescriptionLength),
	}

	ErrIdentifierRegex = &ValidationError{
		"Identifier can only contain lower case alphanumeric characters and dash, should start with alphabet and end with alphanumeric",
	}

	ErrInvalidCharacters = &ValidationError{"Input contains invalid characters."}

	ErrEmailInvalid = &ValidationError{
		"Invalid email address format",
	}

	ErrPasswordLen = &ValidationError{
		fmt.Sprintf("Password has to be within %d and %d characters", minPasswordLength, maxPasswordLength),
	}

	ErrDisplayNameLength = &ValidationError{
		fmt.Sprintf("Display name has to be between %d and %d in length.", minDisplayNameLength, maxDisplayNameLength),
	}

	ErrGitRepoUrl = &ValidationError{"Invalid Git repository url format"}

	ErrFqdn = &ValidationError{"Invalid FQDN format, should start with http:// or https:// and must not include any path or query parameters"}

	ErrIPV4 = &ValidationError{"Invalid IP v4 address format"}

	ErrDirectory = &ValidationError{"Invalid path, should shart with / and should be a valid directory path"}

	ErrVarKeyLength = &ValidationError{
		fmt.Sprintf("Variable key has to be between %d and %d in length.", minVarKey, maxVarKey),
	}
	ErrVarKeyRegex = &ValidationError{"Variable key should start with an alphabet and contain only alphanumeric and underscore"}
)

// ForControlCharacters ensures that there are no control characters in the provided string.
func ForControlCharacters(s string) error {
	for _, r := range s {
		if r < 32 || r == 127 {
			return ErrInvalidCharacters
		}
	}

	return nil
}

// Email checks the provided email and returns an error if it isn't valid.
func Email(email string) error {
	if ok, _ := regexp.Match(emailRegex, []byte(email)); !ok {
		return ErrEmailInvalid
	}

	return nil
}

// Password check the provided password and returns an error if it isn't valid.
func Password(password string) error {
	l := len(password)
	if l < minPasswordLength || l > maxPasswordLength {
		return ErrPasswordLen
	}

	return nil
}

// GitRepo checks the provided git repo url and returns an error if it isn't valid.
func GitRepo(repo string) error {
	httpsOk, _ := regexp.MatchString(gitRepoHttpsPattern, repo)
	// sshOk, _ := regexp.MatchString(gitRepoSshPattern, repo)
	if httpsOk {
		return ForControlCharacters(repo)
	}

	return ErrGitRepoUrl
}

// DisplayName checks the provided display name and returns an error if it isn't valid.
func DisplayName(displayName string) error {
	l := len(displayName)
	if l < minDisplayNameLength || l > maxDisplayNameLength {
		return ErrDisplayNameLength
	}

	return ForControlCharacters(displayName)
}

// Identifier checks the provided identifier and returns an error if it isn't valid.
func Identifier(identifier string) error {
	l := len(identifier)
	if l < minIdentifierLength || l > MaxIdentifierLength {
		return ErrIdentifierLength
	}

	if ok, _ := regexp.Match(identifierRegex, []byte(identifier)); !ok {
		return ErrIdentifierRegex
	}

	return nil
}

// Description checks the provided description and returns an error if it isn't valid.
func Description(description string) error {
	l := len(description)
	if l > maxDescriptionLength {
		return ErrDescriptionTooLong
	}

	return ForControlCharacters(description)
}

// FQDN checks the provided fqdn and returns an error if it isn't valid.
func FQDN(fqdn string) error {
	if ok, err := regexp.MatchString(fqdnPattern, fqdn); !ok || err != nil {
		return ErrFqdn
	}

	return ForControlCharacters(fqdn)
}

// IPV4 checks the provided ipv4 and returns an error if it isn't valid.
func IPV4(ipv4 string) error {
	if ok, _ := regexp.MatchString(ipV4Pattern, ipv4); !ok {
		return ErrIPV4
	}
	return ForControlCharacters(ipv4)
}

// Directory check the provided directory and returns an error if it isn't valid.
func Directory(dir string) error {
	if ok, _ := regexp.MatchString(directoryPattern, dir); !ok {
		return ErrDirectory
	}

	return ForControlCharacters(dir)
}

// VariableKey checks the provided variable key and returns an error if it isn't valid.
func VariableKey(key string) error {
	l := len(key)
	if l < minVarKey || l > maxVarKey {
		return ErrVarKeyLength
	}
	if ok, _ := regexp.MatchString(varKeyRegex, key); !ok {
		return ErrVarKeyRegex
	}
	return ForControlCharacters(key)
}
