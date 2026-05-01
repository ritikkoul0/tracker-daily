# Daily Tracker Backend

A Golang backend API for tracking daily activities including learning, sleep, and office hours.

## Setup

1. Install dependencies:
```bash
cd backend
go mod download
```

2. Run the server:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Activities
- `POST /api/activities` - Create a new activity
- `GET /api/activities` - Get all activities
- `GET /api/activities/:date` - Get activity by date (YYYY-MM-DD)
- `PUT /api/activities/:id` - Update an activity
- `DELETE /api/activities/:id` - Delete an activity

### Analysis
- `GET /api/analysis/weekly` - Get weekly analysis
- `GET /api/analysis/monthly` - Get monthly analysis
- `GET /api/analysis/yearly` - Get yearly analysis

## Example Request

```json
POST /api/activities
{
  "date": "2026-04-30",
  "learning_hours": 3.5,
  "sleep_hours": 7.5,
  "office_hours": 8.0,
  "notes": "Productive day"
}