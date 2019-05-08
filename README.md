# **The little detective**
## function:
1. The application can crawl data form IMDB ranking
2. Use elastic search to indexing and searching

## How to run:
1. Run all container to run application
    `make compose`
2. Crawl data `localhost:8080/make`
3. Sync data to Elastic search `localhost:8080/sync`
4. Search a file `localhost:8080/`