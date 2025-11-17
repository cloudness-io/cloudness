package vinstance

import (
	"github.com/cloudness-io/cloudness/app/web/views/components/icons"
	"github.com/cloudness-io/cloudness/app/web/views/shared"
)

const (
	InstanceNavSettings string = "Settings"
	InstanceNavAuth     string = "Auth"
	InstanceNavRegistry string = "Registry"
	InstanceNavServer   string = "Server"
)

func getInstanceNav() []*shared.PageNavItem {
	return []*shared.PageNavItem{
		{
			Name:      InstanceNavSettings,
			Icon:      icons.LogoIcon,
			ActionUrl: "/settings",
		},
		{
			Name:      InstanceNavAuth,
			Icon:      icons.AuthIcon,
			ActionUrl: "/settings/auth",
		},
		{
			Name:      InstanceNavRegistry,
			Icon:      icons.SourceRegistryIcon,
			ActionUrl: "/settings/registry",
		},
		{
			Name:      InstanceNavServer,
			Icon:      icons.ServerIcon,
			ActionUrl: "/settings/server",
		},
	}
}
