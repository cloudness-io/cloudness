package enum

// TokenType represents the type of the JWT token.
type TokenType string

const (
	// TokenTypeSession is the token returned during user login or signup.
	TokenTypeSession TokenType = "session"

	// TokenTypePAT is a personal access token.
	TokenTypePAT TokenType = "pat"

	// TokenTypeSAT is a service account access token.
	TokenTypeSAT TokenType = "sat"
)
