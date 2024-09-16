package scraper

import (
	"context"
	"log"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/joaopugsley/hltv-scraper/internal/models"
	"github.com/joaopugsley/hltv-scraper/internal/utils"
)

func ScrapeMatches(ctx context.Context) (map[string][]models.MatchData, error) {
	var upcomingMatchesSection []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.hltv.org/matches"),
		chromedp.WaitReady(`div.upcomingMatchesAll`),
		chromedp.WaitReady(`div.matchInfo`),
		chromedp.WaitReady(`div.matchTime`),
		chromedp.WaitReady(`div.matchTeamName`),
		chromedp.WaitReady(`div.matchEventName`),
		chromedp.Nodes(`div.upcomingMatchesAll div.upcomingMatchesSection`, &upcomingMatchesSection),
	)

	if err != nil {
		log.Fatal(err)
	}

	matchesData := make(map[string][]models.MatchData)

	for _, section := range upcomingMatchesSection {
		matchDayHeadline := utils.GetFirstChildByClass(section, "matchDayHeadline")
		if matchDayHeadline == nil {
			continue
		}

		var matchDayHeadlineText string
		err := chromedp.Run(ctx,
			chromedp.Text(matchDayHeadline.FullXPath(), &matchDayHeadlineText, chromedp.NodeVisible),
		)

		if err != nil {
			log.Fatal(err)
		}

		matches := utils.GetChildrenByClass(section, "upcomingMatch")
		if len(matches) == 0 {
			continue
		}

		for _, matchElement := range matches {
			match := matchElement.Children[0]
			if match == nil {
				continue
			}

			matchData := extractMatchData(ctx, match)

			if matchData != nil {
				matchesData[matchDayHeadlineText] = append(matchesData[matchDayHeadlineText], *matchData)
			}
		}
	}

	return matchesData, nil
}

func extractMatchData(ctx context.Context, match *cdp.Node) *models.MatchData {
	hltvUrl := match.AttributeValue("href")
	matchTime := utils.DeepGetFirstChildByClass(match, "matchTime")
	matchMeta := utils.DeepGetFirstChildByClass(match, "matchMeta")

	matchEventName := utils.DeepGetFirstChildByClass(match, "matchEventName")
	if matchEventName == nil {
		return nil
	}

	matchTeams := utils.DeepGetFirstChildByClass(match, "matchTeams")
	if matchTeams == nil {
		return nil
	}

	matchTeam1 := utils.GetFirstChildByClass(matchTeams, "team1")
	matchTeam2 := utils.GetFirstChildByClass(matchTeams, "team2")
	if matchTeam1 == nil || matchTeam2 == nil {
		return nil
	}

	matchTeam1Name := utils.GetFirstChildByClass(matchTeam1, "matchTeamName")
	matchTeam2Name := utils.GetFirstChildByClass(matchTeam2, "matchTeamName")

	if matchTeam1Name == nil {
		matchTeam1Name = utils.GetFirstChildByClass(matchTeam1, "team")
	}

	if matchTeam2Name == nil {
		matchTeam2Name = utils.GetFirstChildByClass(matchTeam2, "team")
	}

	if matchTeam1Name == nil || matchTeam2Name == nil {
		return nil
	}

	return &models.MatchData{
		HltvURL: hltvUrl,
		Format:  utils.GetElementText(ctx, matchMeta),
		Time:    utils.GetElementText(ctx, matchTime),
		Teams: [2]models.TeamData{
			{
				Name: utils.GetElementText(ctx, matchTeam1Name),
			},
			{
				Name: utils.GetElementText(ctx, matchTeam2Name),
			},
		},
		Event: models.EventData{
			Name: utils.GetElementText(ctx, matchEventName),
		},
	}
}
