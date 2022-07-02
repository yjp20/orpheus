package queue

import (
	"fmt"
	"testing"
	"time"

	"github.com/yjp20/orpheus/pkg/music"
)

var (
	snowDrive         = &music.Song{Name: "Araki-Snow Drive", Length: duration("4m26s")}
	osuMemories       = &music.Song{Name: "osu!memories", Length: duration("7m54s")}
	osuMemories2      = &music.Song{Name: "SakiZ - osu!memories 2", Length: duration("7m44s")}
	mathPodcast1      = &music.Song{Name: "Alex Kontorovich: Improving math | 3b1b podcast #1", Length: duration("1h24m22s")}
	mathPodcast3      = &music.Song{Name: "Steven Strogatz: In and out of love with math | 3b1b podcast #3", Length: duration("1h54m7s")}
	mathPodcast5      = &music.Song{Name: "Tai-Danae Bradley: Where math meets language | 3b1b Podcast #5", Length: duration("1h36m5s")}
	flyMeToTheMoon    = &music.Song{Name: "Fly Me To The Moon (2008 Remastered)", Length: duration("2m27s")}
	justTheTwoOfUs    = &music.Song{Name: "Grover Washington Jr. feat. Bill Withers - Just The Two of Us [HQ]", Length: duration("3m56s")}
	rickRoll          = &music.Song{Name: "Rick Astley - Never Gonna Give You Up (Official Music Video)", Length: duration("3m32s")}
	sougetsuEli       = &music.Song{Name: "Aoi Chou - Sougetsu Eli", Length: duration("3m35s")}
	lingus            = &music.Song{Name: "Snarky Puppy - Lingus (We Like It Here)", Length: duration("10m43s")}
	flintstones       = &music.Song{Name: "Flintstones - Jacob Collier", Length: duration("3m10s")}
	nipponEgaoHyakkei = &music.Song{Name: "Nippon Egao Hyakkei", Length: duration("3m57s")}
	offlinePodcast    = &music.Song{Name: "JAE ON TWITCH ft. eaJ - OfflineTV Podcast #8", Length: duration("1h7m28s")}
	catchMeIfYouCan   = &music.Song{Name: "Catch Me If You Can", Length: duration("3m26s")}
	helloWorld        = &music.Song{Name: "BUMP OF CHICKEN「Hello,world!」", Length: duration("4m22s")}
	marshmary         = &music.Song{Name: "マシュマリー / feat.初音ミク", Length: duration("3m32s")}
	empireStateOfMind = &music.Song{Name: "Empire State Of Mind", Length: duration("4m36s")}
	kekkaiSensen      = &music.Song{Name: "blood blockade battlefront ED 1 full", Length: duration("4m8s")}
	strange           = &music.Song{Name: "Celeste - Strange (Official Video)", Length: duration("3m30s")}
)

type addEvent struct {
	songs  []*music.Song
	user   string
	policy AddPolicy
}

type moveEvent struct {
	songs  []*music.Song
	user   string
	policy AddPolicy
}

type removeEvent struct {
	songs  []*music.Song
	user   string
	policy AddPolicy
}

type skipEvent struct {
	skip int
}

type test struct {
	events []interface{}
	ans    []*music.Song
}

var tests = []test{
	{
		[]interface{}{
			addEvent{[]*music.Song{snowDrive}, "mycho", Smart},
			addEvent{[]*music.Song{osuMemories}, "mycho", Smart},
			addEvent{[]*music.Song{osuMemories2}, "mycho", Smart},
		},
		[]*music.Song{
			snowDrive,
			osuMemories,
			osuMemories2,
		},
	},
	{
		[]interface{}{
			addEvent{[]*music.Song{mathPodcast1}, "mycho", Smart},
			addEvent{[]*music.Song{mathPodcast3}, "mycho", Smart},
			addEvent{[]*music.Song{mathPodcast5}, "mycho", Smart},
			addEvent{[]*music.Song{flyMeToTheMoon}, "david", Smart},
			addEvent{[]*music.Song{justTheTwoOfUs}, "david", Smart},
		},
		[]*music.Song{
			mathPodcast1,
			mathPodcast3,
			flyMeToTheMoon,
			mathPodcast5,
			justTheTwoOfUs,
		},
	},
	{
		[]interface{}{
			addEvent{[]*music.Song{rickRoll}, "mycho", Smart},
			addEvent{[]*music.Song{sougetsuEli}, "yjp", Smart},
			addEvent{[]*music.Song{lingus}, "yjp", Smart},
			addEvent{[]*music.Song{flintstones}, "theory", Smart},
			addEvent{[]*music.Song{nipponEgaoHyakkei}, "mycho", Smart},
			addEvent{[]*music.Song{offlinePodcast}, "theory", Smart},
		},
		[]*music.Song{
			rickRoll,
			sougetsuEli,
			flintstones,
			nipponEgaoHyakkei,
			lingus,
			offlinePodcast,
		},
	},
	{
		[]interface{}{
			addEvent{[]*music.Song{catchMeIfYouCan, helloWorld, marshmary, empireStateOfMind, kekkaiSensen}, "mycho", Smart},
			addEvent{[]*music.Song{rickRoll}, "rlarkdfus", Smart},
			addEvent{[]*music.Song{strange}, "rlarkdfus", Smart},
		},
		[]*music.Song{
			catchMeIfYouCan,
			helloWorld,
			rickRoll,
			marshmary,
			strange,
			empireStateOfMind,
			kekkaiSensen,
		},
	},
	{
		[]interface{}{
			addEvent{[]*music.Song{catchMeIfYouCan, helloWorld, marshmary}, "mycho", Smart},
			addEvent{[]*music.Song{empireStateOfMind}, "rlarkdfus", Smart},
			addEvent{[]*music.Song{kekkaiSensen}, "rlarkdfus", Next},
			addEvent{[]*music.Song{rickRoll}, "mycho", Now},
			addEvent{[]*music.Song{flyMeToTheMoon}, "rlarkdfus", Last},
		},
		[]*music.Song{
			catchMeIfYouCan,
			rickRoll,
			kekkaiSensen,
			helloWorld,
			empireStateOfMind,
			marshmary,
			flyMeToTheMoon,
		},
	},
	{
		[]interface{}{
			addEvent{[]*music.Song{empireStateOfMind}, "mycho", Smart},
			addEvent{[]*music.Song{rickRoll}, "mycho", Smart},
			skipEvent{1},
			addEvent{[]*music.Song{catchMeIfYouCan, helloWorld, marshmary}, "mycho", Smart},
		},
		[]*music.Song{
			empireStateOfMind,
			rickRoll,
			catchMeIfYouCan,
			helloWorld,
			marshmary,
		},
	},
}

func TestQueue(t *testing.T) {
	for idx, test := range tests {
		t.Run(fmt.Sprintf("test %d", idx), func(t *testing.T) {
			queue := NewQueue()
			for i, event := range test.events {
				t.Logf("%d", i)
				switch e := event.(type) {
				case addEvent:
					queue.Add(e.songs, e.user, false, e.policy)
				case skipEvent:
					queue.SkipTo((queue.Index + e.skip) % len(queue.List))
				}
				t.Logf("%s", queue.debugQueue())
			}
			if len(queue.List) != len(test.ans) {
				t.Errorf("Expected queue of length '%d', got '%d'", len(test.ans), len(queue.List))
			}
			for i, item := range queue.List {
				if item.Song != test.ans[i] {
					t.Errorf("Expected song '%s' at index '%d', got '%s'", test.ans[i].Name, i, item.Song.Name)
				}
			}
		})
	}
}

func duration(s string) time.Duration {
	temp, _ := time.ParseDuration(s)
	return temp
}

func (queue *Queue) debugQueue() string {
	s := ""
	for i, item := range queue.List {
		if i == queue.Index {
			s = s + "* " + item.Song.Name + "\n"
		} else {
			s = s + "  " + item.Song.Name + "\n"
		}
	}
	return s
}
