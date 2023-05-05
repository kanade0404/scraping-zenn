package scraping_zenn

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"

	"github.com/kanade0404/scraping-zenn/logger"
)

type ArticleContent struct {
	URL         string    `json:"url"`
	Emoji       string    `json:"emoji"`
	Title       string    `json:"title"`
	PublishedAt time.Time `json:"published_at"`
	Topics      []string  `json:"topics"`
	GoodCount   int       `json:"good_count"`
}

func SearchArticle(ctx context.Context, url string) (*ArticleContent, error) {
	logger.Infof("start SearchArticle. url: %s\n", url)
	defer logger.Infof("end SearchArticle. url: %s\n", url)
	const (
		emojiSelector       = "#__next > article > header > div > div > div.ArticleHeader_emoji__GEgLg > span.Emoji_nativeEmoji__JRjFi"
		titleSelector       = "#__next > article > header > div > div > h1 > span"
		publishedAtSelector = "#__next > article > div.Container_wide__i6D5f.Container_common__bSTKj > div > div > div.View_columnsContainer__mDFW9 > aside > div > div.ArticleSidebar_articleInfo__YqxUh.ArticleSidebar_sidebarCard__w397P > dl > div:nth-child(2) > dd > span > time"
		topicsSelector      = "#__next > article > div.Container_wide__i6D5f.Container_common__bSTKj > div > div > div.View_columnsContainer__mDFW9 > section > div.View_main__es7o7 > div > div.View_topics__OVMdM > a > div.View_topicName__rxKth"
		goodCountSelector   = "#share > div.LikeButton_container__lhHqp.style-large > span"
	)
	var (
		emoji            string
		title            string
		publishedAt      string
		existPublishedAt bool
		topicsNodes      []*cdp.Node
		goodCount        string
		topics           []string
	)
	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(url),
		// header
		chromedp.WaitVisible("#__next > article > header > div > div"),
		// article content
		chromedp.WaitVisible("#__next > article > div.Container_wide__i6D5f.Container_common__bSTKj > div > div"),
		chromedp.Text(emojiSelector, &emoji),
		chromedp.Text(titleSelector, &title),
		chromedp.AttributeValue(publishedAtSelector, "datetime", &publishedAt, &existPublishedAt, chromedp.ByQuery),
		chromedp.Nodes(topicsSelector, &topicsNodes, chromedp.ByQueryAll),
		chromedp.Text(goodCountSelector, &goodCount),
	}); err != nil {
		return nil, err
	}
	for i := range topicsNodes {
		logger.Infof("topic: %+v\n", topicsNodes[i])
		var text string
		if err := chromedp.Run(ctx, chromedp.Text(topicsNodes[i].FullXPath(), &text)); err != nil {
			return nil, err
		}
		topics = append(topics, text)
	}
	articleContent := &ArticleContent{
		URL:    url,
		Emoji:  emoji,
		Title:  title,
		Topics: topics,
	}
	if existPublishedAt {
		if t, err := time.Parse(time.RFC3339, publishedAt); err != nil {
			return nil, err
		} else {
			articleContent.PublishedAt = t
		}
	} else {
		return nil, fmt.Errorf("'publishedAt' is not found. selector: %s", publishedAtSelector)
	}
	if v, err := strconv.Atoi(goodCount); err != nil {
		return nil, err
	} else {
		articleContent.GoodCount = v
	}
	return articleContent, nil
}
