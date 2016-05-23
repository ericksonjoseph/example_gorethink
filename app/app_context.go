package main

import r "github.com/dancannon/gorethink"

/**
 * Holds stuff needed throughout the App
 */
type AppContext struct {
	DB *r.Session
}
