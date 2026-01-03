package main

import (
	"expence_management/api"
	"log"
)

func initDB(h *api.Handler) {

	// add admin user to list of users

	err := h.Repo.Repo.AddAdminUser()
	if err != nil {
		log.Printf("cannot add admin user")
	}
}
