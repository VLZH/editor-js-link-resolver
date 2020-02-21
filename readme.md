# Link resolver for @editorjs/link 🐳
Resolver parses opengraph markup on a site and returns information for [Link Tool](https://github.com/editor-js/link) in [Editor.js](https://editorjs.io)


**Response example:**
```json
{
    "meta": {
        "description": "#93: Версус: Node.js или PHP в 2019",
        "image": {
            "url": "https://miro.medium.com/max/1200/1*LKldKAfENCqlBbLOjfte3A.jpeg"
        },
        "title": "Девшахта-подкаст"
    },
    "success": 1
}
```
Request example:

```bash
# http is httpie
http http://localhost:9000/fetchUrl\?url\=https://medium.com/devschacht/devschacht-93-ac5e4b21e696
```
## Run in docker
```bash
docker pull vlzhvlzh/editor-js-link-resolver
docker run -p 9000:9000 -e PORT=9000 -e HOST=0.0.0.0 -e ALLOW_ORIGIN='*' vlzhvlzh/editor-js-link-resolver
```
### Env variables
- PORT
- HOST
- ALLOW_ORIGIN - `Access-Control-Allow-Origin` header value

## TODO
- Fix error handling
- Add support HTML meta tags
