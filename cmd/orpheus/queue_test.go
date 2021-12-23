package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

var tests = []struct {
	songs [][]*Song
	users []string
	policies []addPolicy
	ans   string
}{
	{
		[][]*Song {
			[]*Song{
				&Song{
					Name: "Araki-Snow Drive",
					Length: first("4m26s"),
				},
			},
			[]*Song{
				&Song{
					Name: "osu!memories",
					Length: first("7m54s"),
				},
			},
			[]*Song{
				&Song{
					Name: "SakiZ - osu!memories 2",
					Length: first("7m44s"),
				},
			},
		},
		[]string{
			"mycho",
			"mycho",
			"mycho",
		},
		[]addPolicy{
			Smart,
			Smart,
			Smart,
		},
		"0.  **Araki-Snow Drive** (0:00/4:26)\n1.  **osu!memories** (7:54)\n2.  **SakiZ - osu!memories 2** (7:44)",
	},
	{
		[][]*Song{
			[]*Song{
				&Song{
					Name: "Alex Kontorovich: Improving math | 3b1b podcast #1",
					Length: first("1h24m22s"),
				},
			},
			[]*Song{
				&Song{
					Name: "Steven Strogatz: In and out of love with math | 3b1b podcast #3",
					Length: first("1h54m7s"),
				},
			},
			[]*Song{
				&Song{
					Name: "Tai-Danae Bradley: Where math meets language | 3b1b Podcast #5",
					Length: first("1h36m5s"),
				},
			},
			[]*Song{
				&Song{
					Name: "Fly Me To The Moon (2008 Remastered)",
					Length: first("2m27s"),
				},
			},
			[]*Song{
				&Song{
					Name: "Grover Washington Jr. feat. Bill Withers - Just The Two of Us [HQ]",
					Length: first("3m56s"),
				},
			},
		},
		[]string{
			"mycho",
			"mycho",
			"mycho",
			"david",
			"david",
		},
		[]addPolicy{
			Smart,
			Smart,
			Smart,
			Smart,
			Smart,
		},
		"0.  **Alex Kontorovich: Improving math | 3b1b podcast #1** (0:00/1:24:22)\n1.  **Fly Me To The Moon (2008 Remastered)** (2:27)\n2.  **Grover Washington Jr. feat. Bill Withers - Just The Two of Us [HQ]** (3:56)\n3.  **Steven Strogatz: In and out of love with math | 3b1b podcast #3** (1:54:07)\n4.  **Tai-Danae Bradley: Where math meets language | 3b1b Podcast #5** (1:36:05)",
	},
	{
		[][]*Song{
			[]*Song{
				&Song{
					Name: "Rick Astley - Never Gonna Give You Up (Official Music Video)",
					Length: first("3m32s"),
				},
			},
			[]*Song{
				&Song{
					Name: "Aoi Chou - Sougetsu Eli",
					Length: first("3m35s"),
				},
			},
			[]*Song{
				&Song{
					Name: "Snarky Puppy - Lingus (We Like It Here)",
					Length: first("10m43s"),
				},
			},
			[]*Song{
				&Song{
					Name: "Flintstones - Jacob Collier",
					Length: first("3m10s"),
				},
			},
			[]*Song{
				&Song{
					Name: "Nippon Egao Hyakkei",
					Length: first("3m57s"),
				},
			},
			[]*Song{
				&Song{
					Name: "JAE ON TWITCH ft. eaJ - OfflineTV Podcast #8",
					Length: first("1h7m28s"),
				},
			},
		},
		[]string{
			"mycho",
			"yjp",
			"yjp",
			"theory",
			"mycho",
			"theory",
		},
		[]addPolicy{
			Smart,
			Smart,
			Smart,
			Smart,
			Smart,
			Smart,
		},
		"0.  **Rick Astley - Never Gonna Give You Up (Official Music Video)** (0:00/3:32)\n1.  **Flintstones - Jacob Collier** (3:10)\n2.  **Aoi Chou - Sougetsu Eli** (3:35)\n3.  **Nippon Egao Hyakkei** (3:57)\n4.  **Snarky Puppy - Lingus (We Like It Here)** (10:43)\n5.  **JAE ON TWITCH ft. eaJ - OfflineTV Podcast #8** (1:07:28)",
	},
	{
		[][]*Song{
			[]*Song{
				&Song{
					Name: "Catch Me If You Can",
					Length: first("3m26s"),
				},
				&Song{
					Name: "BUMP OF CHICKEN「Hello,world!」",
					Length: first("4m22s"),
				},
				&Song{
					Name: "マシュマリー / feat.初音ミク",
					Length: first("3m32s"),
				},
				&Song{
					Name: "Empire State Of Mind",
					Length: first("4m36s"),
				},
				&Song{
					Name: "blood blockade battlefront ED 1 full",
					Length: first("4m8s"),
				},
			},
			[]*Song{
				&Song{
					Name: "Rick Astley - Never Gonna Give You Up (Official Music Video)",
					Length: first("3m32s"),
				},
			},
			[]*Song{
				&Song{
					Name: "Celeste - Strange (Official Video)",
					Length: first("3m30s"),
				},
			},
		},
		[]string{
			"mycho",
			"rlarkdfus",
			"rlarkdfus",
		},
		[]addPolicy{
			Smart,
			Smart,
			Smart,
		},
		"0.  **Catch Me If You Can** (0:00/3:26)\n1.  **Rick Astley - Never Gonna Give You Up (Official Music Video)** (3:32)\n2.  **Celeste - Strange (Official Video)** (3:30)\n3.  **BUMP OF CHICKEN「Hello,world!」** (4:22)\n4.  **マシュマリー / feat.初音ミク** (3:32)\n5.  **Empire State Of Mind** (4:36)\n6.  **blood blockade battlefront ED 1 full** (4:08)",
	},
	{
		[][]*Song{
			[]*Song{
				&Song{
					Name: "Catch Me If You Can",
					Length: first("3m26s"),
				},
				&Song{
					Name: "BUMP OF CHICKEN「Hello,world!」",
					Length: first("4m22s"),
				},
				&Song{
					Name: "マシュマリー / feat.初音ミク",
					Length: first("3m32s"),
				},
			},
			[]*Song{
				&Song{
					Name: "Empire State of Mind",
					Length: first("4m36s"),
				},
			},
			[]*Song{
				&Song{
					Name: "blood blockade battlefront ED 1 full",
					Length: first("4m8s"),
				},
			},
			[]*Song{
				&Song{
					Name: "Rick Astley - Never Gonna Give You Up (Official Music Video)",
					Length: first("3m32s"),
				},
			},
			[]*Song{
				&Song{
					Name: "Fly Me To The Moon (2008 Remastered)",
					Length: first("2m27s"),
				},
			},
		},
		[]string{
			"mycho",
			"rlarkdfus",
			"rlarkdfus",
			"mycho",
			"rlarkdfus",
		},
		[]addPolicy{
			Smart,
			Smart,
			Next,
			Now,
			Last,
		},
		"0.  **Catch Me If You Can** (3:26)\n1.  **Rick Astley - Never Gonna Give You Up (Official Music Video)** (0:00/3:32)\n2.  **blood blockade battlefront ED 1 full** (4:08)\n3.  **BUMP OF CHICKEN「Hello,world!」** (4:22)\n4.  **Empire State of Mind** (4:36)\n5.  **マシュマリー / feat.初音ミク** (3:32)\n6.  **Fly Me To The Moon (2008 Remastered)** (2:27)",
	},
}

func TestAdd(t *testing.T) {
	for i, test := range tests {
		server := getServer(strconv.Itoa(i))
		server.Player = Player{Time: 0}
		for index, song := range test.songs {
			fmt.Printf("YEP %s\n", test.users[index])
			fmt.Printf("%d %d\n", server.getQueueSum(test.users[index]), server.currentRank())
			item := server.Add(song, test.users[index], false, test.policies[index])
			fmt.Printf("item number %d Index %d\n", index, item[0].Rank)
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
		{
			[]int { -1 },
			[]int { -1 },
			tests[4].ans,
		},
	}

	for i, test := range tests {
		server := getServer("skip"+strconv.Itoa(i))
		server.Player = Player{ Time: 0 }
		ct := 0
		for index, song := range test.songs {
			item := server.Add(song, test.users[index], false, test.policies[index])
			fmt.Printf("item number %d Index %d\n", index, item[0].Rank)
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
		{
			[]int {},
			[]int {},
			tests[4].ans,
		},
	}
	for i, test := range tests {
		server := getServer("move"+strconv.Itoa(i))
		server.Player = Player{ Time: 0 }
		for index, song := range test.songs {
			item := server.Add(song, test.users[index], false, test.policies[index])
			fmt.Printf("item number %d Index %d\n", index, item[0].Rank)
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
		{
			[]int { -1 },
			[]int { -1 },
			tests[4].ans,
		},
	}
	for i, test := range tests {
		server := getServer("remove"+strconv.Itoa(i))
		server.Player = Player{ Time: 0 }
		ct := 0
		for index, song := range test.songs {
			item := server.Add(song, test.users[index], false, test.policies[index])
			fmt.Printf("item number %d Index %d\n", index, item[0].Rank)
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

func first(s string) time.Duration {
	temp, _ := time.ParseDuration(s)
	return temp
}
