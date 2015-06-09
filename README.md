##### Romdex (rather Nesdex for now)

Collection of tools to indentify and scrape rom information, only NES
for now.

All code is work in progress.

* Goodimport; Scrape goodnes renamed files and extract the data to
  import it into a mongodb database
* Goodscrape: Scrape imported results with gamesdb data so we can lookup
  extra details via their API.
* GoodApi: Serve the imported/scrapes result over an RESTful API
* GoodFrontend: Combine all sources to build an emulator front-end
