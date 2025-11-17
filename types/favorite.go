package types

type Favorite struct {
	UserID        int64 `json:"-" db:"favorite_user_id"`
	ApplicationID int64 `json:"-" db:"favorite_application_id"`

	Created int64 `json:"created" db:"favorite_created"`
}

type FavoriteDTO struct {
	ProjectUID      int64  `json:"project_uid"`
	ProjectName     string `json:"project_name"`
	EnvironmentUID  int64  `json:"environment_uid"`
	EnvironmentName string `json:"environment_name"`
	ApplicationUID  int64  `json:"application_uid"`
	AppName         string `json:"app_name"`
	AppDomain       string `json:"app_domain"`
}
