package users

import "time"

// UserModuleOptions contains configuration options for the user module.
type UserModuleOptions struct {
	EnableCache     bool
	CacheExpiration time.Duration
}

