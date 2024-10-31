## About The Project

This personal project is a backend RESTful api using CRUD operations to store data in a PostgreSQL database. The idea was a simple model for nursing students to store notes for studying.

The purpose of this project was for me to develop a full-scale project focusing on architectural implementations, using Chi router for endpoint paths, and GORM for database transactions.

## Reflection
  - What was the context of this project?
      * It was an experimental project working with Docker, PostgreSQL in Go for the first time. Its intended for me to continuously improve upon over time. 

### Built With

* [![Go][Go]][Go-url]
* [![Docker][Docker]][Docker-url]
* [![PostgreSQL][PostgreSQL]][PostgreSQL-url]

[Go]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[Go-url]: https://go.dev/
[Docker]: https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white
[Docker-url]: https://www.docker.com/
[PostgreSQL]: https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white
[PostgreSQL-url]: https://www.postgresql.org/

### Getting Started

Project is still under development.

#### Prerequisites
* Go 1.23 or higher
* Docker and Docker Compose

### Installation

#### Running Docker Environment Setup

**Make sure Docker is running**
1. Clone the repository:
```sh
git clone https://github.com/MichaelNguyenTT/nurse-api.git
cd <project location>
```
2. Run and start the container:
```sh
docker-compose up --build
```
3. Open another terminal and post to one of the endpoints:
```sh
curl -X POST http://localhost:8080/v1/service \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Dummy Name",
    "category": "Test Category",
    "notes": "Test Notes"
  }'
```
4. Get the response back:
```sh
curl -i http://localhost:8080/v1/service
```

## Roadmap
Area to focus on continuous improvement throughout the project.

Will add more ideas when I come across more
- [ ] Tiger-style assertions [Documentations](https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/TIGER_STYLE.md)
- [ ] Validation, Logging, Error middleware handlers
- [ ] API Documentation with Swagger
- [ ] Better structuring to centralize my errors
