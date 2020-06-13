package external

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"github.com/shuymn/nijisanji-db-collector/src/domain/entity"
	"github.com/shuymn/nijisanji-db-collector/src/domain/repository"
)

type ichikaraRepository struct {
	URL *url.URL
}

const URL = "https://nijisanji.ichikara.co.jp/member/"

func NewIchikaraRepository() (repository.IchikaraRepository, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	r := &ichikaraRepository{
		URL: u,
	}
	return r, nil
}

func (r *ichikaraRepository) FetchLivers() ([]*entity.Liver, error) {
	urls, err := r.fetchLiverProfilePageURLs()
	if err != nil {
		return nil, err
	}

	ls := make([]*entity.Liver, 0, len(urls))

	var wg sync.WaitGroup
	for _, u := range urls {
		u := u
		wg.Add(1)
		go func() {
			defer wg.Done()

			id, err := r.getIDByURL(u)
			if err != nil {
				log.Error().Err(err).Send()
				return
			}
			l := &entity.Liver{
				ID: strings.ToLower(id),
			}

			res, err := http.Get(u.String())
			if err != nil {
				log.Error().Err(err).Send()
				return
			}
			defer res.Body.Close()

			if res.StatusCode != 200 {
				// TODO: impl error
				return
			}

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				// TODO: impl error
				return
			}

			doc.Find("a.elementor-social-icon-youtube").Each(func(_ int, s *goquery.Selection) {
				id, err := r.getIDBySelection(s)
				if err != nil {
					log.Error().Err(err).Send()
					return
				}
				l.YouTube = append(l.YouTube, &entity.YouTube{ID: id})
			})

			doc.Find("a.elementor-social-icon-twitter").Each(func(_ int, s *goquery.Selection) {
				id, err := r.getIDBySelection(s)
				if err != nil {
					log.Error().Err(err).Send()
					return
				}
				l.Twitter = append(l.Twitter, &entity.Twitter{ID: id})
			})

			doc.Find("a.elementor-social-icon-").Each(func(_ int, s *goquery.Selection) {
				id, err := r.getIDBySelection(s)
				if err != nil {
					log.Error().Err(err).Send()
					return
				}
				l.Bilibili = append(l.Bilibili, &entity.Bilibili{ID: id})
			})

			doc.Find("a.elementor-social-icon-twitch").Each(func(_ int, s *goquery.Selection) {
				id, err := r.getIDBySelection(s)
				if err != nil {
					log.Error().Err(err).Send()
					return
				}
				l.Twitch = append(l.Twitch, &entity.Twitch{ID: id})
			})

			doc.Find("a.elementor-social-icon-facebook").Each(func(_ int, s *goquery.Selection) {
				u, err := r.getURLBySelection(s)
				if err != nil {
					log.Error().Err(err).Send()
					return
				}
				u.RawQuery = ""

				id, err := r.getIDByURL(u)
				if err != nil {
					log.Error().Err(err).Send()
					return
				}

				l.Facebook = append(l.Facebook, &entity.Facebook{ID: id, URL: u.String()})
			})

			doc.Find("a.elementor-social-icon-instagram").Each(func(_ int, s *goquery.Selection) {
				u, err := r.getURLBySelection(s)
				if err != nil {
					log.Error().Err(err).Send()
					return
				}
				u.RawQuery = ""

				id, err := r.getIDByURL(u)
				if err != nil {
					log.Error().Err(err).Send()
					return
				}

				l.Instagram = append(l.Instagram, &entity.Instagram{ID: id, URL: u.String()})
			})

			if err = l.Validate(); err != nil {
				log.Error().Err(err).Send()
				return
			}

			ls = append(ls, l)
		}()
	}

	wg.Wait()

	return ls, nil
}

func (r *ichikaraRepository) fetchLiverProfilePageURLs() ([]*url.URL, error) {
	res, err := http.Get(r.URL.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, nil // TODO: impl error
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	sel := doc.Find(".elementor-tab-content a")
	urls := make([]*url.URL, 0, sel.Length())

	var wg sync.WaitGroup
	for i := range sel.Nodes {
		s := sel.Eq(i)
		wg.Add(1)
		go func() {
			defer wg.Done()

			u, err := r.getURLBySelection(s)
			if err != nil {
				log.Error().Err(err).Send()
				return
			}

			if u.Hostname() == "" {
				u, err = url.Parse(r.URL.String() + u.String())
				if err != nil {
					log.Error().Err(err).Send()
					return
				}
			}

			urls = append(urls, u)
		}()
	}

	wg.Wait()

	return urls, nil
}

func (r *ichikaraRepository) getIDBySelection(s *goquery.Selection) (string, error) {
	u, err := r.getURLBySelection(s)
	if err != nil {
		return "", err
	}
	return r.getIDByURL(u)
}

func (r *ichikaraRepository) getIDByURL(u *url.URL) (string, error) {
	paths := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(paths) == 0 {
		return "", errors.New("paths is empty")
	}
	return paths[len(paths)-1], nil
}

func (r *ichikaraRepository) getURLBySelection(s *goquery.Selection) (*url.URL, error) {
	href, ok := s.Attr("href")
	if !ok {
		// TODO: impl error
		return nil, errors.New("")
	}
	return url.Parse(strings.TrimSpace(href))
}
