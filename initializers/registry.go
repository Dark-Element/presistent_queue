package initializers

import (
	"persistentQueue/services"
	"os"
	"fmt"
)

// This is a static registry, initialized at boot
type Registry struct {
	Messaging services.MessagingInterface
}

// Add more dependencies to the registry here
func GetRegistry(sigs chan os.Signal, done chan bool) *Registry {
	registry := Registry{
		Messaging: services.InitMessaging(),
	}
	go func(){
		<-sigs
		registry.Close(done)
	}()
	return &registry
}


func (r *Registry) Close(done chan bool){
	fmt.Println("Closing Registry")
	defer fmt.Println("Closed Registry")
	r.Messaging.Close()
	done <- true
}