package main

import (
    "os"
)

type Song struct {
    Name string
    URL string
    Length float64
    File *os.File
}

type QueueItem struct {
    Song Song
    QueuedBy string
    Index float64
}



