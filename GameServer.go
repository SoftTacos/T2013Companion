package main

type GameServer struct {
	Players []PlayerClient
	DM      DMClient
}

type Client interface {
	ClientMethod()
}

type PlayerClient struct {
}

type DMClient struct {
}
