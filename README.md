# Orion

**Orion** is a file storage and delivery service designed for the TeaCoder platform, using S3 and PostgreSQL for cloud storage and data management.

## Tech Stack

- **Programming Language**: Go (Golang)
- **Web Framework**: [Gin](https://gin-gonic.com/)
- **Cloud Storage**: S3
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **ORM**: [GORM](https://gorm.io/)

## API Documentation

| Method     | URL               | Description            |
|------------|-------------------|------------------------|
| `GET`      | `/`               | API health check       |
| `POST`     | `/upload`         | Upload a file          |
| `GET`      | `/:tag/:id`       | Retrieve a file        |

## License

This project is licensed under the [GNU Affero General Public License v3.0](https://github.com/teacoder-team/storage/blob/master/LICENSE).
