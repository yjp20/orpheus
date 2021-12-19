package main

import (
	"fmt"
	"strconv"
	"testing"
)

var tests = []struct {
	links []string
	users []string
	ans   string
}{
	{
		[]string{
			"https://www.youtube.com/watch?v=Nz4ODnhQPAk",
			"https://www.youtube.com/watch?v=UOYHaKHnBsk",
			"https://www.youtube.com/watch?v=WW3OojhKdUM",
		},
		[]string{
			"mycho",
			"mycho",
			"mycho",
		},
		"0.  **Araki-Snow Drive** (0:00/4:26)\n1.  **osu!memories** (7:54)\n2.  **SakiZ - osu!memories 2** (7:44)",
	},
	{
		[]string{
			"https://www.youtube.com/watch?v=C-i4q-Xlnis",
			"https://www.youtube.com/watch?v=SUMLKweFAYk",
			"https://www.youtube.com/watch?v=pvRY3r-b0QI",
			"https://www.youtube.com/watch?v=ZEcqHA7dbwM",
			"https://www.youtube.com/watch?v=PJ0u5c9EF1E",
		},
		[]string{
			"mycho",
			"mycho",
			"mycho",
			"david",
			"david",
		},
		"0.  **Alex Kontorovich: Improving math | 3b1b podcast #1** (0:00/1:24:22)\n1.  **Fly Me To The Moon (2008 Remastered)** (2:27)\n2.  **Grover Washington Jr. feat. Bill Withers - Just The Two of Us [HQ]** (3:56)\n3.  **Steven Strogatz: In and out of love with math | 3b1b podcast #3** (1:54:07)\n4.  **Tai-Danae Bradley: Where math meets language | 3b1b Podcast #5** (1:36:05)",
	},
	{
		[]string{
			"https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			"https://www.youtube.com/watch?v=2f5Az8TqN0k",
			"https://www.youtube.com/watch?v=L_XJ_s5IsQc",
			"https://www.youtube.com/watch?v=zua831utwMM",
			"https://www.youtube.com/watch?v=2ch_p04MTHg",
			"https://www.youtube.com/watch?v=diPciEKB9-s",
		},
		[]string{
			"mycho",
			"yjp",
			"yjp",
			"theory",
			"mycho",
			"theory",
		},
		"0.  **Rick Astley - Never Gonna Give You Up (Official Music Video)** (0:00/3:32)\n1.  **Flintstones - Jacob Collier** (3:10)\n2.  **Aoi Chou - Sougetsu Eli** (3:35)\n3.  **Nippon Egao Hyakkei** (3:57)\n4.  **Snarky Puppy - Lingus (We Like It Here)** (10:43)\n5.  **JAE ON TWITCH ft. eaJ - OfflineTV Podcast #8** (1:07:28)",
	},
	{
		[]string{
			"https://www.youtube.com/playlist?list=PLGl3Zr2INfG0R0LXAWPCR_SZXGV16pu4i",
			"https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			"https://www.youtube.com/watch?v=A1AJEv50Ld4",
		},
		[]string{
			"mycho",
			"rlarkdfus",
			"rlarkdfus",
		},
		"0.  **Catch Me If You Can** (0:00/3:26)\n1.  **Rick Astley - Never Gonna Give You Up (Official Music Video)** (3:32)\n2.  **Celeste - Strange (Official Video)** (3:30)\n3.  **BUMP OF CHICKEN「Hello,world!」** (4:22)\n4.  **マシュマリー / feat.初音ミク** (3:32)\n5.  **Empire State Of Mind** (4:36)\n6.  **blood blockade battlefront ED 1 full** (4:08)",
	},
}

func TestAdd(t *testing.T) {
	for i, test := range tests {
		server := getServer(strconv.Itoa(i))
		server.Player = Player{Time: 0}
		for index, url := range test.links {
			fmt.Printf("YEP %s\n", test.users[index])
			fmt.Printf("%d %d\n", server.getQueueSum(test.users[index]), server.currentIndex())
			songs, _ := fetchSongsFromURL(url, true)
			item := server.Add(songs, test.users[index], false)
			fmt.Printf("item number %d Index %d\n", index, item[0].Index)
		}
		queue := PrintQueue(server)
		if queue != test.ans {
			t.Errorf("Current Queue:\n%s\nDesired Queue:\n%s", queue, test.ans)
		}
	}
}

func TestSkipTo(t *testing.T) {
	skipTests := []struct {
		toindex []int
		whenIndex []int
		ans string
	}{
		{
			[]int {
				2,
				0,
			},
			[]int {
				2,
				0,
			},
			"0.  **Araki-Snow Drive** (4:26)\n1.  **osu!memories** (7:54)\n2.  **SakiZ - osu!memories 2** (0:00/7:44)",
		},
		{
			[]int {
				2,
				0,
			},
			[]int {
				2,
				0,
			},
			"0.  **Alex Kontorovich: Improving math | 3b1b podcast #1** (1:24:22)\n1.  **Steven Strogatz: In and out of love with math | 3b1b podcast #3** (1:54:07)\n2.  **Tai-Danae Bradley: Where math meets language | 3b1b Podcast #5** (0:00/1:36:05)\n3.  **Fly Me To The Moon (2008 Remastered)** (2:27)\n4.  **Grover Washington Jr. feat. Bill Withers - Just The Two of Us [HQ]** (3:56)",
		},
		{
			[]int {
				2,
				1,
				4,
			},
			[]int {
				2,
				4,
				5,
			},
			"0.  **Rick Astley - Never Gonna Give You Up (Official Music Video)** (3:32)\n1.  **Aoi Chou - Sougetsu Eli** (3:35)\n2.  **Snarky Puppy - Lingus (We Like It Here)** (10:43)\n3.  **Flintstones - Jacob Collier** (3:10)\n4.  **Nippon Egao Hyakkei** (0:00/3:57)\n5.  **JAE ON TWITCH ft. eaJ - OfflineTV Podcast #8** (1:07:28)",
		},
		{
			[]int { -1 },
			[]int { -1 },
			"0.  **Catch Me If You Can** (0:00/3:26)\n1.  **Rick Astley - Never Gonna Give You Up (Official Music Video)** (3:32)\n2.  **Celeste - Strange (Official Video)** (3:30)\n3.  **BUMP OF CHICKEN「Hello,world!」** (4:22)\n4.  **マシュマリー / feat.初音ミク** (3:32)\n5.  **Empire State Of Mind** (4:36)\n6.  **blood blockade battlefront ED 1 full** (4:08)",
		},
	}

	for i, test := range tests {
		server := getServer("skip"+strconv.Itoa(i))
		server.Player = Player{ Time: 0 }
		ct := 0
		for index, url := range test.links {
			songs, _ := fetchSongsFromURL(url, true)
			item := server.Add(songs, test.users[index], false)
			fmt.Printf("item number %d Index %d\n", index, item[0].Index)
			if skipTests[i].whenIndex[ct] == index{
				server.SkipTo(skipTests[i].toindex[ct])
				ct += 1
			}
		}
		queue := PrintQueue(server)
		if queue != skipTests[i].ans{
			t.Errorf("Current Queue:\n%s\nDesired Queue:\n%s", queue, skipTests[i].ans)
		}
	}
}

func TestMove(t *testing.T){
	MoveTests := []struct {
		fromIndex []int
		toIndex []int
		ans string
	} {
		{
			[]int {
				1,
				0,
			},
			[]int {
				2,
				2,
			},
			"0.  **SakiZ - osu!memories 2** (7:44)\n1.  **osu!memories** (7:54)\n2.  **Araki-Snow Drive** (0:00/4:26)",
		},
		{
			[]int {
				2,
				2,
			},
			[]int {
				0,
				0,
			},
			"0.  **Fly Me To The Moon (2008 Remastered)** (2:27)\n1.  **Grover Washington Jr. feat. Bill Withers - Just The Two of Us [HQ]** (3:56)\n2.  **Alex Kontorovich: Improving math | 3b1b podcast #1** (0:00/1:24:22)\n3.  **Steven Strogatz: In and out of love with math | 3b1b podcast #3** (1:54:07)\n4.  **Tai-Danae Bradley: Where math meets language | 3b1b Podcast #5** (1:36:05)",
		},
		{
			[]int { 5 },
			[]int { 1 },
			"0.  **Rick Astley - Never Gonna Give You Up (Official Music Video)** (0:00/3:32)\n1.  **JAE ON TWITCH ft. eaJ - OfflineTV Podcast #8** (1:07:28)\n2.  **Flintstones - Jacob Collier** (3:10)\n3.  **Aoi Chou - Sougetsu Eli** (3:35)\n4.  **Nippon Egao Hyakkei** (3:57)\n5.  **Snarky Puppy - Lingus (We Like It Here)** (10:43)",
		},
		{
			[]int {},
			[]int {},
			"0.  **Catch Me If You Can** (0:00/3:26)\n1.  **Rick Astley - Never Gonna Give You Up (Official Music Video)** (3:32)\n2.  **Celeste - Strange (Official Video)** (3:30)\n3.  **BUMP OF CHICKEN「Hello,world!」** (4:22)\n4.  **マシュマリー / feat.初音ミク** (3:32)\n5.  **Empire State Of Mind** (4:36)\n6.  **blood blockade battlefront ED 1 full** (4:08)",
		},
	}
	for i, test := range tests {
		server := getServer("move"+strconv.Itoa(i))
		server.Player = Player{ Time: 0 }
		for index, url := range test.links {
			songs, _ := fetchSongsFromURL(url, true)
			item := server.Add(songs, test.users[index], false)
			fmt.Printf("item number %d Index %d\n", index, item[0].Index)
		}
		for index, from := range MoveTests[i].fromIndex {
			server.Move(from, MoveTests[i].toIndex[index])
		}
		queue := PrintQueue(server)
		if queue != MoveTests[i].ans{
			t.Errorf("Current Queue:\n%s\nDesired Queue:\n%s", queue, MoveTests[i].ans)
		}
	}
}

func TestRemove(t *testing.T) {
	RemoveTests := []struct {
		rindex []int
		whenIndex []int
		ans string
	} {
		{
			[]int {
				1,
				0,
			},
			[]int {
				1,
				2,
			},
			"0.  **SakiZ - osu!memories 2** (0:00/7:44)",
		},
		{
			[]int {
				0,
				0,
				2,
			},
			[]int {
				0,
				3,
				4,
			},
			"0.  **Fly Me To The Moon (2008 Remastered)** (0:00/2:27)\n1.  **Grover Washington Jr. feat. Bill Withers - Just The Two of Us [HQ]** (3:56)",
		},
		{
			[]int {
				0,
				4,
			},
			[]int {
				0,
				5,
			},
			"0.  **Aoi Chou - Sougetsu Eli** (0:00/3:35)\n1.  **Flintstones - Jacob Collier** (3:10)\n2.  **Nippon Egao Hyakkei** (3:57)\n3.  **Snarky Puppy - Lingus (We Like It Here)** (10:43)",
		},
		{
			[]int { -1 },
			[]int { -1 },
			"0.  **Catch Me If You Can** (0:00/3:26)\n1.  **Rick Astley - Never Gonna Give You Up (Official Music Video)** (3:32)\n2.  **Celeste - Strange (Official Video)** (3:30)\n3.  **BUMP OF CHICKEN「Hello,world!」** (4:22)\n4.  **マシュマリー / feat.初音ミク** (3:32)\n5.  **Empire State Of Mind** (4:36)\n6.  **blood blockade battlefront ED 1 full** (4:08)",
		},
	}
	for i, test := range tests {
		server := getServer("remove"+strconv.Itoa(i))
		server.Player = Player{ Time: 0 }
		ct := 0
		for index, url := range test.links {
			songs, _ := fetchSongsFromURL(url, true)
			item := server.Add(songs, test.users[index], false)
			fmt.Printf("item number %d Index %d\n", index, item[0].Index)
			if RemoveTests[i].whenIndex[ct] == index{
				server.Remove(RemoveTests[i].rindex[ct])
				ct += 1
			}
		}
		queue := PrintQueue(server)
		if queue != RemoveTests[i].ans{
			t.Errorf("Current Queue:\n%s\nDesired Queue:\n%s", queue, RemoveTests[i].ans)
		}
	}
}
