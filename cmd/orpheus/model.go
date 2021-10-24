package main

type Song struct {
    Name string
    URL string
    Length int
    File string
}

type QueueItem struct {
    Song Song
    QueuedBy string
    Index float64
}



