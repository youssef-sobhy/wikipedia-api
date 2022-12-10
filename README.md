# Wikipedia API
This is an API to get short descriptions of Wikipedia articles.

Available at https://wikipedia.youssefsobhy.com/api/v1

Documentation available at https://wikipedia.youssefsobhy.com/api/v1/docs/index.html

![Go Version](https://img.shields.io/github/go-mod/go-version/youssef1337/wikipedia-api)
![Uptime](https://img.shields.io/uptimerobot/ratio/m793223758-f73506a770999c5e13ade54f)
![Status](https://img.shields.io/uptimerobot/status/m793223758-f73506a770999c5e13ade54f)
![B&T](https://github.com/youssef1337/wikipedia-api/actions/workflows/build-and-run-tests.yaml/badge.svg)
![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)

## Table of contents
- [Wikipedia API](#wikipedia-api)
  - [Table of contents](#table-of-contents)
  - [Installation](#installation)
  - [Usage](#usage)
  - [API Reference and Documentation](#api-reference-and-documentation)
  - [Built With](#built-with)
  - [Deployment](#deployment)
  - [Monitoring](#monitoring)
  - [License](#license)
  - [Authors and acknowledgment](#authors-and-acknowledgment)

## Installation

- Make sure you have Golang installed on your machine (https://go.dev/doc/install)
- Clone the repository and change directory to the project
  ```bash
  git clone https://github.com/youssef1337/wikipedia-api.git
  cd wikipedia-api
  ```
- Install the dependencies
  ```bash
  go get -t ./...
  ```
- Run the server
  ```bash
  go run cmd/main.go
  ```
- Open your browser and go to http://localhost:3000 to see the API in action

## Usage

- To get a short description of an article, send a GET request to http://localhost:3000/api/v1/search with the article name as a query parameter
  ```bash
  curl http://localhost:3000/api/v1/search?query=Yoshua_Bengio
  ```
- To check if the API is running, send a GET request to http://localhost:3000/api/v1
  ```bash
  curl http://localhost:3000/api/v1
  ```

## API Reference and Documentation
- [Wikipedia API](https://en.wikipedia.org/w/api.php) - The Wikipedia API I used to get the short descriptions
- [API Documentation](https://wikipedia.youssefsobhy.com/api/v1/docs/index.html) - The API documentation of this project

## Built With
- [Golang](https://golang.org/) - The programming language used
- [Gin](https://github.com/gin-gonic/gin) - The web framework used
- [Swag](https://github.com/swaggo/swag) - The API documentation generator used

## Deployment
- [Render](https://render.com/) - The cloud platform used to deploy the API
- [Cloudflare](https://www.cloudflare.com/) - The service used to manage the DNS records

## Monitoring
- [UptimeRobot](https://uptimerobot.com/) - The service used to monitor the API uptime and status

## License
[Apache License 2.0](https://choosealicense.com/licenses/apache-2.0/)

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## Authors and acknowledgment
- Youssef Sobhy - [youssef1337](https://github.com/youssef1337)
