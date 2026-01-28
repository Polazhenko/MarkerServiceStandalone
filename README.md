# Marker Service

A simplified, standalone version of the MarkerService microservice for portfolio demonstration.

## Features

- **CRUD Operations**: Create, Read, Update, Delete markers
- **Category Management**: Organize markers by categories (General, Scouting, Harvests, Weather)
- **Geolocation**: Store latitude/longitude coordinates for each marker
- **In-Memory Storage**: Simple repository pattern with in-memory data storage
- **RESTful API**: Clean API design with proper HTTP methods
- **Swagger Documentation**: Interactive API documentation

## Technology Stack

- **.NET 8.0**: Modern C# web framework
- **ASP.NET Core Web API**: RESTful API framework
- **AutoMapper**: Object-to-object mapping
- **Swagger/OpenAPI**: API documentation

## Architecture

The project follows clean architecture principles:

```
├── Controllers/        # API endpoints
├── Services/          # Business logic layer
├── Repositories/      # Data access layer
├── Models/            # Domain models and DTOs
└── Program.cs         # Application entry point
```

## API Endpoints

### Create Marker
```http
POST /api/v1/markers
Headers: X-User-Id: 1
Body: {
  "name": "Deer Sighting",
  "category": 2,
  "latitude": 45.5,
  "longitude": -93.2,
  "description": "Large buck spotted"
}
```

### Get Marker
```http
GET /api/v1/markers/{markerId}
Headers: X-User-Id: 1
```

### Get All Markers
```http
GET /api/v1/markers?category=2
Headers: X-User-Id: 1
```

### Update Marker
```http
PATCH /api/v1/markers/{markerId}
Headers: X-User-Id: 1
Body: {
  "name": "Updated Name",
  "description": "Updated description"
}
```

### Delete Marker
```http
DELETE /api/v1/markers/{markerId}
Headers: X-User-Id: 1
```

## Running the Application

```bash
cd MarkerServiceStandalone
dotnet restore
dotnet run
```

Navigate to `https://localhost:5001/swagger` to view the API documentation.

## Key Differences from Original

This standalone version simplifies the original MarkerService by:

- Removing AWS dependencies (S3, DynamoDB, SQS)
- Removing Redis caching
- Removing external service integrations
- Using in-memory storage instead of DynamoDB
- Simplified authentication (header-based user ID)
- Removed image upload functionality
- Removed map tools and advanced features
- Minimal logging configuration

## Original Service Features

The original MarkerService includes:
- DynamoDB for persistent storage
- AWS S3 for image storage
- Redis for caching
- Integration with Device Manager and Subscription services
- Image upload and thumbnail generation
- Map tools (areas, lines, shapes)
- Soil data integration
- Weather station data
- Rain station information
- Sharing functionality
- OpenTelemetry metrics

## Use Case

This simplified version demonstrates:
- Clean architecture and separation of concerns
- Repository pattern implementation
- Service layer with business logic
- RESTful API design
- Dependency injection
- Async/await patterns
- AutoMapper usage
- Swagger documentation
