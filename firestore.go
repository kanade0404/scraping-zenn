package scraping_zenn

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/firestore"

	"github.com/kanade0404/scraping-zenn/logger"
)

type FireStoreClient struct {
	ctx     context.Context
	client  *firestore.Client
	current time.Time
}

func NewClient(ctx context.Context, projectID string, current time.Time) (*FireStoreClient, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &FireStoreClient{
		ctx, client, current,
	}, err
}

func (s *FireStoreClient) Close() {
	err := s.client.Close()
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
func (s *FireStoreClient) UpdateArticles(collectionID string, articles []*ArticleContent) error {
	logger.Info("start UpdateArticles")
	defer logger.Info("end UpdateArticles")

	logger.Infof("start creating document. total: %d", len(articles))
	zennDocRef := s.client.Collection(collectionID).Doc("zenn")
	ch := make(chan error, len(articles))
	var wg sync.WaitGroup
	wg.Add(len(articles))
	for i := range articles {
		go func(v *ArticleContent, idx int) {
			defer wg.Done()
			logger.Infof("start creating document. url: %s %d/%d", v.URL, idx, len(articles))
			u, err := url.Parse(v.URL)
			if err != nil {
				ch <- err
				return
			}
			paths := strings.Split(u.Path, "/")
			if len(paths) < 4 {
				ch <- errors.New("invalid url")
				return
			}
			logger.Infof("paths: %+v %d/%d", paths, idx, len(articles))
			articleID := paths[3]
			logger.Infof("articleID: %s %d/%d", articleID, idx, len(articles))
			docRef := zennDocRef.Collection(articleID).Doc(strconv.FormatInt(s.current.UnixNano(), 10))
			jsonStr, err := json.Marshal(v)
			if err != nil {
				ch <- err
				return
			}
			var articleMap map[string]interface{}
			if err := json.Unmarshal(jsonStr, &articleMap); err != nil {
				ch <- err
				return
			}
			if _, err := docRef.Set(s.ctx, articleMap); err != nil {
				ch <- err
				return
			}
			logger.Infof("end creating document. url: %s %d/%d", v.URL, idx, len(articles))
		}(articles[i], i)
	}
	wg.Wait()
	close(ch)
	var errs error
	for err := range ch {
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}
	if errs != nil {
		return errs
	}
	return nil
}
