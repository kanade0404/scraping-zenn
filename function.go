package scraping_zenn

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/kanade0404/scraping-zenn/logger"
)

type request struct {
	Name string `json:"name"`
}

func Scraping(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed) // 405
		Response(w, http.StatusMethodNotAllowed, map[string]string{"error": "This method is not allowed"})
		return
	}
	logger.Info("start Scraping")
	defer logger.Info("end Scraping")
	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Init()
	kService := os.Getenv("K_SERVICE")
	var isLocal bool
	if kService == "" {
		isLocal = true
	}
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", isLocal),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
		chromedp.Flag("hide-scrollbars", false),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("mute-audio", false))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	alc, _ := chromedp.NewExecAllocator(ctx, opts...)
	ctx, cancel = chromedp.NewContext(alc, chromedp.WithLogf(log.Printf))
	defer cancel()
	client, err := NewClient(ctx, os.Getenv("PROJECT_ID"), currentTime)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer client.Close()
	var (
		parents          []*cdp.Node
		articleCardNodes []*cdp.Node
	)
	url := fmt.Sprintf("https://zenn.dev/%s", req.Name)
	logger.Infof("start chromedp. target: %v", url)
	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible("#__next > div.View_contents___goft > div > div.FadeInUp_fadeInUp__U9uWt > div", chromedp.ByQuery),
		chromedp.Nodes("#__next > div.View_contents___goft > div > div.FadeInUp_fadeInUp__U9uWt > div", &parents, chromedp.NodeVisible),
	}); err != nil {
		logger.Error(err.Error())
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(parents) < 1 {
		err := errors.New("cannot query parents elements")
		logger.Errorf(err.Error())
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	logger.Infof("parent articleCardNodes leb: %d\n", len(parents))
	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Nodes("#__next > div.View_contents___goft > div > div.FadeInUp_fadeInUp__U9uWt > div > article > a.ArticleCard_mainLink__X2TOE", &articleCardNodes, chromedp.ByQueryAll),
	}); err != nil {
		logger.Error(err.Error())
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	logger.Infof("articleCardNodes: %d", len(articleCardNodes))
	logger.Infof("end chromedp. target: %v\n", url)
	var articles []*ArticleContent
	for nodeIdx := range articleCardNodes {
		logger.Infof("exec SearchArticleCard: %d/%d", nodeIdx+1, len(articleCardNodes)+1)
		relLink, ok := articleCardNodes[nodeIdx].Attribute("href")
		if !ok {
			logger.Errorf("not found href. index: %d", nodeIdx)
		}
		if article, err := SearchArticleCard(ctx, fmt.Sprintf("https://zenn.dev%s", relLink)); err != nil {
			logger.Errorf("failure SearchArticleCard: %d/%d, error: %s", nodeIdx+1, len(articleCardNodes)+1, err.Error())
		} else {
			articles = append(articles, article)
			logger.Infof("success SearchArticleCard: %d/%d", nodeIdx+1, len(articleCardNodes)+1)
		}
	}
	if err := client.UpdateArticles(req.Name, articles); err != nil {
		logger.Error(err.Error())
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	type response struct {
		Articles []*ArticleContent `json:"articles"`
	}
	Response(w, http.StatusOK, response{
		Articles: articles,
	})
}
