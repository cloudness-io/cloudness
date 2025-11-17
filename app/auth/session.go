package auth

import "github.com/cloudness-io/cloudness/types"

// Session contains information of the authenticated principal and auth related metadata.
type Session struct {
	// Principal is the authenticated principal.
	Principal types.Principal

	// Metadata contains auth related information (access grants, tokenId, sshKeyId, ...)
	Metadata Metadata
}
