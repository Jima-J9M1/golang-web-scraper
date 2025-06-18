# Web Scraper Project

A concurrent web scraper written in Go that efficiently scrapes web pages and stores the extracted links in a SQLite database.

## Features

- Concurrent web scraping using worker goroutines
- Automatic worker scaling based on CPU cores
- SQLite database integration for persistent storage
- Robust error handling for network, parsing, and database operations
- Configurable timeout and context management
- Efficient resource utilization

## Project Structure

```
.
├── main.go           # Main application entry point
├── go.mod           # Go module definition
├── go.sum           # Go module checksums
├── scraper.db       # SQLite database file
├── configs/         # Configuration files
└── internal/        # Internal packages
    ├── scraper/     # Web scraping logic
    └── storage/     # Database operations
```

## Prerequisites

- Go 1.24.3 or later
- SQLite3

## Dependencies

- golang.org/x/net v0.41.0
- github.com/mattn/go-sqlite3 v1.14.28

## Usage

1. Clone the repository
2. Run the application:
   ```bash
   go run main.go
   ```

The scraper will process a predefined list of URLs and store the extracted links in the SQLite database.

## How It Works

1. The application initializes a SQLite database for storing scraped links
2. Creates a pool of worker goroutines based on available CPU cores
3. Each worker:
   - Fetches web pages concurrently
   - Parses links from the HTML content
   - Stores the links in the database
4. Results are processed and errors are handled appropriately

## Error Handling

The application handles various types of errors:
- Network fetch errors
- HTML parsing errors
- Database operation errors
- Timeout errors

## Performance

- Uses concurrent processing with worker goroutines
- Automatically scales based on available CPU cores
- Implements efficient resource management
- Uses context for timeout control

## License

This project is open source and available under the MIT License. 