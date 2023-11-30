# Github Events

## Overview

This is a monorepo of a duel microservices system.
The system purpose is to allow clients to get the following functionality:
1. List all events.
2. Count the amount of events.
3. List the 20 most recent actors that were involved in the
   events.
4. List the 20 most recent repositories that were involved in
   the events, including the amount of
   stars that each one of them has.

## Architecture

The system is divided to two microservices: Collector and Api.
The collector collects the events from the github api (https://docs.github.com/en/rest/activity/events).
The api provides the mentioned above functionality over HTTP protocol.
Both of the services uses the same database, which is MongoDB.

This architecture solves the following problems:
1. Data will be consistently provided by the api service even if github api is unreachable.
2. Data will be delivered without latency caused by github api.
3. Computational time will be decreased due to the use of query over the MongoDB server.

<b><u>Note</u></b>:
This architecture was designed for a home challenge exercise.
It doesn't take in consideration any overload problems, and tries to provide the necessary functionality without making the design too complicated.

## How To Run The Application Locally

This version doesn't contain a docker-compose support.
To run this application, follow the next steps:
1. Install GO.
2. Download this project.
3. Install all the requirements mentioned in `go.mod`.
4. Create an instance of MongoDB (using Docker should be the easiest way).
5. Put the relevant values as environment variables. The keys are available at `example.env` file.
   1. To run the api server, put the value "API" in the `Role` variable.
   Otherwise, the collector will run.
   2. This code supports the using of `.env` file.
6. Run the code with two instances - one for each service (to run the collector periodically see the mention in the collector section).

## Dive Into The Services

### Database

The chosen database is MongoDB.
The reason for that is the structure of the data received by the github API.
There are no complicated connections between the received entities (events, actors, repos).
This way it's easier to save the data, without any processing.

The DB contains one model: Events.
This model contains the exact same fields as the model from github api, with the addition of the field `repo.stars` which is collected separately.

This code uses the `mongo-driver` package for communicating with MongoDB.
This package was chosen despite the fact that this package is not an ORM package, due to the following reasons: 
1. Its simplicity.
2. Its aggregation functionality.
3. Its support as an official package of MongoDB company.
4. The size of the project (only one model is stored in the DB).

### Collector

The collector is responsible for collecting the events.
Its flow contains three main stages, where each one runs as `goroutine` and they uses `channels` to transfer data between themselves.
1. <u>Collect events stage</u> - Collects the events from github api by using another three `goroutines`, one for each page.
After that it stores each event in a `channel`.
2. <u>Add repo stars stage</u> - Receives each event and uses workers as `goroutines` to get the stars of the repo mentioned in the events.
After that it stores each event (after the stars repo data was added) to another `channel`.
3. <u>Store in DB</u> - Receives each event and uses the `ReplaceOne` MongoDB functionality with `Upsert` option to add all new events and update the exists ones.

This service does not contains a scheduler.
To run this service periodically you will need to add an external scheduler.
This way the service is not responsible for scheduling and remains with one responsibility (as it should due to Single Responsibility Principle).
To make this code run periodically, use the package `github.com/robfig/cron` by adding the next code block:

```go
import (
    "github.com/robfig/cron"
)

func main() {
    c := cron.New()
    c.AddFunc("0 * * * * *", collector.Main)
    c.Start()

    select {} // Block forever
}
```

### Api

The api is responsible for providing the functionality to the client.
Client can use an HTTP request to this server to get the data he needs.
This web servers also provides a `swagger` that is reachable by `/swagger` uri.
This api contains four relevant routes:
1. <u>GET /api/v1/events</u> - Lists all collected events.
This route uses pagination to return a page of events with the wanted size.
2. <u>GET /api/v1/events/count</u> - Counts all collected events.
3. <u>GET /api/v1/events/actors</u> - Gets the actors of recent events.
This route gets the amount of wanted actors and uses aggregation to provides the data as an actor scheme.  
4. <u>GET /api/v1/events/repos</u> - Gets the repos of recent events.
This route gets the amount of wanted repos and uses aggregation to provides the data as a repo scheme.

## Extra Possible Features

1. <u>Make the repo stars update work consistently</u> -
Github api has a rate limit of 60 requests per hour for unauthenticated requests.
This amount of requests limit doesn't meet with the amount of api requests that is needed in order to update all the repos stars.
To make this work properly, we should add an option to authenticate to github api, so the rate limit will be 5,000 requests per hour.
2. <u>Add external scheduler</u> -
There are a lot of ways to schedule a service, and it's better to choose according to the service deployment.  
For example, service that is about to be deployed on Openshift environment, can be scheduled by a Kubernetes cron provided by Openshift.
3. <u>Add sorting options</u> -
Add another functionality to /api/v1/events route that sorts the data (using a DB query) and then paginates it.
4. <u>Add easy local deployment solution</u> -
Docker-compose can provide an easy deployment for the system.
After writing two Dockerfiles (one for each service), add a docker-compose file that contains the deployment config for those two services in addition to the DB configuration.
5. <u>Add web UI</u> -
This would provide a better user experience.
