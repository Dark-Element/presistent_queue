package initializers

import (
	"../services"
)

// This is a static registry, initialized at boot
type Registry struct {
	Messaging services.MessagingInterface
}

// Add more dependencies to the registry here
func GetRegistry() *Registry {
	registry := Registry{
		Messaging: services.InitMessaging(),
	}
	return &registry
}
