package scraping_zenn

import (
	"context"
	"github.com/kanade0404/scraping-zenn/logger"
)

func SearchArticleCard(ctx context.Context, url string) (*ArticleContent, error) {
	logger.Infof("start SearchArticleCard. url: %s", url)
	defer logger.Infof("end SearchArticleCard. url: %s", url)
	if article, err := SearchArticle(ctx, url); err != nil {
		return nil, err
	} else {
		return article, nil
	}
}
