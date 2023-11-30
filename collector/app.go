package collector

import (
	"githubEvents/collector/github"
	"githubEvents/shared"
	"githubEvents/shared/dal"
	"githubEvents/shared/models"
	"log"
	"sync"
)

func Main() {
	client, err := dal.ConnectToMongo(shared.MongoUri)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := dal.DisconnectFromMongo(client)
		if err != nil {
			log.Fatal(err)
		}
	}()

	mongoEventsClient := dal.NewMongoClient[models.Event](*client, shared.MongoDatabaseName, shared.MongoEventsCollectionName)
	githubClient := github.NewGithubClient()

	collectEventsChannel := make(chan models.Event)
	addRepoEventsChannel := make(chan models.Event)

	var wg sync.WaitGroup
	var wgCollect sync.WaitGroup
	var wgRepoAdd sync.WaitGroup

	wg.Add(1)
	go func() {
		collectEvents(githubClient, shared.NumberOfGithubEventsPages, shared.MaximumGithubEventsPageSize, collectEventsChannel, &wg, &wgCollect)
		wgCollect.Wait()
		close(collectEventsChannel)
	}()

	wg.Add(1)
	go func() {
		getEventsRepoStars(shared.GithubRepoWorkers, githubClient, collectEventsChannel, addRepoEventsChannel, &wg, &wgRepoAdd)
		wgRepoAdd.Wait()
		close(addRepoEventsChannel)
	}()

	for i := 1; i <= shared.WriteDBWorkers; i++ {
		wg.Add(1)
		go saveEvents(mongoEventsClient, addRepoEventsChannel, &wg)
	}

	wg.Wait()
}

func collectEvents(githubClient github.GithubClient, numOfPages, perPage int, ch chan<- models.Event, wg, wgCollect *sync.WaitGroup) {
	defer wg.Done()

	for page := 1; page <= numOfPages; page++ {
		wgCollect.Add(1)
		go func(githubClient github.GithubClient, page, perPage int) {
			events, err := githubClient.GetEvents(shared.GithubEventsUrl, page, perPage)
			if err != nil {
				log.Fatal(err)
				return
			}
			for _, event := range events {
				ch <- event
			}
			log.Printf("Page %v collected, got %v events.", page, len(events))
			wgCollect.Done()
		}(githubClient, page, perPage)
	}
}

func getEventsRepoStars(numOfWorkers int, githubClient github.GithubClient, collectedEvents <-chan models.Event, addedRepoEvents chan<- models.Event, wg, wgAddRepo *sync.WaitGroup) {
	defer wg.Done()

	for i := 1; i <= numOfWorkers; i++ {
		wgAddRepo.Add(1)
		go func(githubClient github.GithubClient) {
			for {
				select {
				case event, more := <-collectedEvents:
					if !more {
						wgAddRepo.Done()
						return
					}
					repo, err := githubClient.GetRepo(event.Repo.Url)
					if err != nil {
						log.Print(err)
						wgAddRepo.Done()
						return
					}
					event.Repo.Stars = repo.StargazersCount
					addedRepoEvents <- event
				}
			}
		}(githubClient)
	}

}

func saveEvents(mongoEventsClient dal.MongoClient[models.Event], ch <-chan models.Event, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case event, more := <-ch:
			if !more {
				return
			}
			err := mongoEventsClient.ReplaceOneObj(event, event.ID)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
