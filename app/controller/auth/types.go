package auth

type ChangePasswordSettings struct {
	EnableRegistration  bool `json:"enable_registration,string,omitempty"`
	EnablePasswordLogin bool `json:"enable,string,omitempty"`
}

type DemoUserSettings struct {
	DemoUserEnabled bool `json:"demo_user_enabled,string,omitempty"`
}
