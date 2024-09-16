package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	c "github.com/chromedp/chromedp"
)

type TeamData struct {
	Logo *string `json:"logo,omitempty"`
	Name string  `json:"name,omitempty"`
}

type EventData struct {
	Logo *string `json:"logo,omitempty"`
	Name string  `json:"name,omitempty"`
}

type MatchData struct {
	HltvURL string      `json:"hltv_url"`
	Format  string      `json:"format"`
	Time    string      `json:"time"`
	Teams   [2]TeamData `json:"teams"`
	Event   EventData   `json:"event"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	options := append(c.DefaultExecAllocatorOptions[:],
		c.Flag("disable-web-security", true),
		c.Flag("disable-site-isolation-trials", true),
		c.Flag("headless", true),
	)

	allocCtx, cancel := c.NewExecAllocator(ctx, options...)
	defer cancel()

	ctx, cancel = c.NewContext(allocCtx)
	defer cancel()

	var upcomingMatchesSection []*cdp.Node
	err := c.Run(ctx,
		c.Navigate("https://www.hltv.org/matches"),
		c.WaitReady(`div.upcomingMatchesAll`),
		c.WaitReady(`div.matchInfo`),
		c.WaitReady(`div.matchTime`),
		c.WaitReady(`div.matchTeamName`),
		c.WaitReady(`div.matchEventName`),
		c.Nodes(`div.upcomingMatchesAll div.upcomingMatchesSection`, &upcomingMatchesSection),
	)

	if err != nil {
		log.Fatal(err)
	}

	matchesData := make(map[string][]MatchData)

	for _, section := range upcomingMatchesSection {
		matchDayHeadline := getFirstChildByClass(section, "matchDayHeadline")
		if matchDayHeadline == nil {
			continue
		}

		var matchDayHeadlineText string
		err := c.Run(ctx,
			c.Text(matchDayHeadline.FullXPath(), &matchDayHeadlineText, c.NodeVisible),
		)

		if err != nil {
			log.Fatal(err)
		}

		matches := getChildrenByClass(section, "upcomingMatch")
		if len(matches) == 0 {
			continue
		}

		for _, matchElement := range matches {
			match := matchElement.Children[0]
			if match == nil {
				continue
			}

			hltvUrl := match.AttributeValue("href")
			matchTime := deepGetFirstChildByClass(match, "matchTime")
			matchMeta := deepGetFirstChildByClass(match, "matchMeta")

			matchEventName := deepGetFirstChildByClass(match, "matchEventName")
			if matchEventName == nil {
				continue
			}

			matchTeams := deepGetFirstChildByClass(match, "matchTeams")
			if matchTeams == nil {
				continue
			}

			matchTeam1 := getFirstChildByClass(matchTeams, "team1")
			matchTeam2 := getFirstChildByClass(matchTeams, "team2")
			if matchTeam1 == nil || matchTeam2 == nil {
				continue
			}

			matchTeam1Name := getFirstChildByClass(matchTeam1, "matchTeamName")
			matchTeam2Name := getFirstChildByClass(matchTeam2, "matchTeamName")

			if matchTeam1Name == nil {
				matchTeam1Name = getFirstChildByClass(matchTeam1, "team")
			}

			if matchTeam2Name == nil {
				matchTeam2Name = getFirstChildByClass(matchTeam2, "team")
			}

			if matchTeam1Name == nil || matchTeam2Name == nil {
				continue
			}

			matchData := MatchData{
				HltvURL: hltvUrl,
				Format:  getElementText(ctx, matchMeta),
				Time:    getElementText(ctx, matchTime),
				Teams: [2]TeamData{
					{
						Name: getElementText(ctx, matchTeam1Name),
					},
					{
						Name: getElementText(ctx, matchTeam2Name),
					},
				},
				Event: EventData{
					Name: getElementText(ctx, matchEventName),
				},
			}

			matchesData[matchDayHeadlineText] = append(matchesData[matchDayHeadlineText], matchData)
		}
	}

	matchesJson, err := json.Marshal(matchesData)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(matchesJson))
}

func hasClass(node *cdp.Node, class string) bool {
	parts := strings.Split(node.AttributeValue("class"), " ")
	for _, part := range parts {
		if part == class {
			return true
		}
	}
	return false
}

func getFirstChildByClass(node *cdp.Node, class string) *cdp.Node {
	for _, child := range node.Children {
		if hasClass(child, class) {
			return child
		}
	}
	return nil
}

func getChildrenByClass(node *cdp.Node, class string) []*cdp.Node {
	children := []*cdp.Node{}
	for _, child := range node.Children {
		if hasClass(child, class) {
			children = append(children, child)
		}
	}
	return children
}

// recursively search for the first child node with the given class ><
func deepGetFirstChildByClass(node *cdp.Node, class string) *cdp.Node {
	if node == nil {
		return nil
	}

	if hasClass(node, class) {
		return node
	}

	for _, child := range node.Children {
		if result := deepGetFirstChildByClass(child, class); result != nil {
			return result
		}
	}

	return nil
}

func getElementText(ctx context.Context, node *cdp.Node) string {
	var text string
	err := c.Run(ctx,
		c.Text(node.FullXPath(), &text, c.NodeVisible),
	)

	if err != nil {
		log.Fatal(err)
	}

	return text
}
