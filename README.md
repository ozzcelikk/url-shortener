# url-shortener

url-shortener is a golang app for generate random short url and return.

## Installation
Use the docker url-shortener-go:1.0.0 to install url-shortener. You can pull image from dockerhub

```bash
docker pull kozcelik/url-shortener-go
```

## How it works

```
# -GET "/service"
returns app is alive.

# -POST "/" body: { "url": "http://xyz.com" }
generate random short-url and match with given long url. Generated value is saved in store on the memory and return short-url.

# -GET "/"
returns exist long url value by given short-url and increase visit count.

# -GET "/list"
returns all saved data with visit counts.
```

## Contributing
Pull requests are welcome.

## License
License does not require you to take the license with you to your project.
