package main

import (
    "fmt"
)

/*func tester(w http.ResponseWriter, r *http.Request){
    s := fetchMusicFromURL("https://www.youtube.com/watch?v=_mMyPJSx8RU")
    fmt.Printf("DONE")
    fmt.Fprintf(w, "%+v\n", s)
}*/


func main() {
    s := fetchMusicFromURL("https://www.youtube.com/watch?v=_mMyPJSx8RU")
    fmt.Printf("%+v\n", s)
}
