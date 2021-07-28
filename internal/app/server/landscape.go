package main

import (
	"github.com/kyohmizu/cicd-landscape/internal/pkg/common"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

const (
	LandscapeUrl = "https://raw.githubusercontent.com/cncf/landscape/master/landscape.yml"
)

var client http.Client

type LandScape struct {
	Landscape []struct {
		Category      string
		Name          string
		Subcategories []struct {
			Subcategory string
			Name        string
			Items       []struct {
				Extra struct {
					Accepted    string
					DevStatsUrl string `yaml:"dev_stats_url"`
					ArtworkUrl  string `yaml:"artwork_url"`
					StackOverflowUrl string `yaml:"stack_overflow_url"`
					BlogUrl string `yaml:"blog_url"`
					SlackUrl string `yaml:"slack_url"`
					YoutubeUrl string `yaml:"youtube_url"`
				}
				Name        string
				Description string
				HomepageUrl string `yaml:"homepage_url"`
				Project     string
				RepoUrl     string `yaml:"repo_url"`
				Logo        string
				Crunchbase  string
			}
		}
	}
}

type Project common.Project

func getCicdProjects() ([]Project, error) {
	resp, err := client.Get(LandscapeUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return findCicdProjects(bodyBytes)
	}
	return nil, nil
}

func findCicdProjects(data []byte) ([]Project, error) {
	l := LandScape{}

	var err = yaml.Unmarshal([]byte(data), &l)
	if err != nil {
		return nil, err
	}

	for _, category := range l.Landscape {
		if category.Name == "App Definition and Development" {
			for _, sub := range category.Subcategories {
				if sub.Name == "Continuous Integration & Delivery" {
					var list []Project
					for _, proj := range sub.Items {
						list = append(list, Project{
							Name:        proj.Name,
							Description: proj.Description,
							HomepageUrl: proj.HomepageUrl,
							Project:     proj.Project,
							RepoUrl:     proj.RepoUrl,
							Crunchbase:  proj.Crunchbase,
						})
					}
					return list, nil
				}
			}
		}
	}
	return nil, nil
}
