package main

type Server interface {
	Start()
	Close()
}