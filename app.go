package main

import (
	"context"
)

// App struct represents the application
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetAppInfo returns basic information about the application
func (a *App) GetAppInfo() map[string]string {
	return map[string]string{
		"name":    "Property Management System",
		"version": "0.1.0",
		"status":  "Initial Setup - No Features Implemented",
	}
}
