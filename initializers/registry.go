package initializers

import (
	"presistentQueue/services"
	"presistentQueue/factories"
)

// This is a static registry, initialized at boot
type Registry struct {
	Messaging services.MessagingInterface
}

// Add more dependencies to the registry here
func GetRegistry() *Registry {
	registry := Registry{
		Messaging: services.InitMessaging(factories.DbConn("192.168.239.129", "root", "getalife", 3306, "queue", 100), 10),
	}
	return &registry
}
