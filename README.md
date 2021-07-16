# fork URL shortener
test it here -> https://fork.pw
### Project structure
The project is split into two parts, the backend which is written in Go, and the frontend in Vuejs 3. The frontend part is pretty simple since it's only used for ease of testing. For the backend, it uses firestore as a DB and global counter management.

### Infrastructure
The project is host in two separated parts, the frontend is on google cloud storage (which doesn't provide TLS) with Cloudflare in front to upgrade the connection to HTTPS. For the backend, it's on google cloud run with CD enable when pushing on the main branch of the git project.

### Architecture
For the short link id, I've used base62 encoding since there's less chance for collision than md5 and also added an atomic global counter which is in firestore, so that multiple instances of the service can run simultaneously without generating the same short id. I've also use cloud run so that I don't have to scale manually if the service was getting a lot of traffic, but one downside is that after some time the first request can be a little bit longer to respond since it needs to startup the container first. But if it was a real service with constant traffic it probably won't be a problem.

## API
The API consists of 4 primary endpoints.

- `/gen-link` is use to generate new link, it needs a json payload with the url you want to shortened.
It accepts `POST, OPTIONS` methods.
```json
{ "url": "https://www.google.com" }
```
- `/update-link` is use for updating an existing link with a new url, it also needs a json payload.
It accepts `POST, OPTIONS` methods.
``` json
{"shortUrl": "https://fork.pw/aBre3", "newUrl": "https://github.com"}
```
- `/info-link/{id}` returns the info about an existing link, `{id}` needs to be replace with a short url id.
- It accepts `GET` method.
``` bash
	curl -XGET https://fork.pw/info-link/yBbPUt
```
- `/{id}` redirects the short link to the long url. `{id}` needs to be replace with a short url id.
In case that the upstream site is down it will fallback to a WaybackMachine snapshot.
It accept `GET` method.
``` bash
	curl -XGET https://fork.pw/yBbPUt
```
