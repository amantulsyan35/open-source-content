# Open Source Content API

For the past three years, I've been tracking the content I consume. It began as a [simple behavioral experiment](https://x.com/at_aman35/status/1489668741335896064) aimed at predicting how my consumption shapes my thinking and problem-solving approaches.

Over time, it evolved into a curious pursuit and a core thesis on how I operate,

[Personal Databases](https://www.amantulsyan.com/personal-databases)

This open API consolidates all the content I consume, making it embeddable and paving the way for innovative applications powered by this dataset—especially in an era dominated by LLMs.

## Features

- **Notion Integration**: Fetch data from multiple Notion databases
- **Pagination**: Cursor-based pagination for efficient data retrieval
- **Rate Limiting**: Protection against excessive API usage


## Prerequisites

- Go 1.19 or higher
- A Notion API key with access to your desired databases

## Installation

1. Clone the repository:
```bash
git clone https://github.com/amantulsyan35/open-source-content.git
cd open-source-content
```

2. Install dependencies:
```bash
go mod download
```

3. Create your environment configuration:
```bash
cp .env.example .env
```

4. Edit the `.env` file with your Notion API credentials:
```
NOTION_API_KEY=your_notion_api_key
NOTION_ROOT_PAGE_ID=your_root_page_id
SERVER_PORT=1323
```

## Usage

### Running the server

Start the server:
```bash
go run main.go
```

The API will be available at `http://localhost:1323` (or the port you specified in your configuration).

### API Endpoints

#### GET /v1

Retrieves content entries from Notion databases with pagination support.

**Query Parameters**:
- `pageSize` (optional): Number of entries to return per page (default: 20, max: 100)
- `cursor` (optional): Pagination cursor for fetching the next page of results

**Example Response**:
```json
{
  "entries": [
    {
      "title": "Example Entry",
      "url": "https://example.com",
      "createdTime": "2023-01-01T12:00:00Z"
    },
    ...
  ],
  "nextCursor": "20",
  "hasMore": true
}
```

**Pagination**:
To fetch the next page of results, include the `nextCursor` value from the previous response:
```
GET /v1?cursor=20
```

## Configuration

The following environment variables can be configured:

| Variable | Description | Default |
|----------|-------------|---------|
| NOTION_API_KEY | Your Notion integration token | None (Required) |
| NOTION_ROOT_PAGE_ID | ID of the root Notion page containing your databases | None (Required) |
| SERVER_PORT | Port for the API server | 1323 |

## Rate Limiting

The API implements rate limiting to prevent abuse:
- 3 requests per second per IP address
- Burst capacity of up to 9 requests (3 seconds worth)
- Rate limiter state resets after 1 minute of inactivity

When rate limits are exceeded, the API returns a 429 status code.

