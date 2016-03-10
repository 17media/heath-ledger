```
.__                   __  .__              .__             .___
|  |__   ____ _____ _/  |_|  |__           |  |   ____   __| _/ ____   ___________
|  |  \_/ __ \\__  \\   __\  |  \   ______ |  | _/ __ \ / __ | / ___\_/ __ \_  __ \
|   Y  \  ___/ / __ \|  | |   Y  \ /_____/ |  |_\  ___// /_/ |/ /_/  >  ___/|  | \/
|___|  /\___  >____  /__| |___|  /         |____/\___  >____ |\___  / \___  >__|
     \/     \/     \/          \/                    \/     \/_____/      \/
```
 
## Docker
**To build**
`docker build -t heath-ledger .`

**To run**
```
docker run -v `pwd`:/go/src/app -e "dev=1" -p 3000:3000 --name="heath-ledger" heath-ledger`
```
*Dev mode* `-e dev=1 to turn on dev mode`

## Note
- please put `_test.go` files adjacent to the file, e.g. `controllers/index.go` will have `controllers/index_test.go` for testing
