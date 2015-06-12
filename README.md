### Romdex (rather Nesdex for now)

Collection of tools to indentify and scrape rom information, only NES
for now.

All code is work in progress.

* Goodimport; Scrape goodnes renamed files and extract the data to
  import it into a mongodb database
* Goodscrape: Scrape imported results with gamesdb data so we can lookup
  extra details via their API.
* GoodApi: Serve the imported/scrapes result over an RESTful API
* GoodFrontend: Combine all sources to build an emulator front-end launches fogleman/nes's emulator on click

#### Screenshot

![Screenshot of Alpha version](http://i.imgur.com/T0oOJyF.png)


#### Goals of the project

1. Index all known (NES) roms and have all the information about them in a public accessible API
2. Collect other meta-data (screenshots / cover art) from various sources and add this data to the API
3. Have one API to rule them all so all frontends have access to this data
4. Preserve history

#### Todo
1. Add more misc information about dumps
2. Find a way to group hacks of games with different names under their
   parent rom
3. Finish code
