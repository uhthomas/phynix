package pheonix

type Room map[*Client]bool

type Rooms map[string]Room
