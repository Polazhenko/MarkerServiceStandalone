# MarkerServiceStandalone — Agents Notes

This document documents the `MarkerServiceStandalone` project for automation agents and contributors. It summarizes architecture, runtime and test commands, important implementation details, and troubleshooting notes.

## Project overview

- Purpose: simple marker management API (create, read, update, delete) used for demos and end-to-end tests.
- Tech stack: .NET 8 (C# 12), ASP.NET Core Web API, AutoMapper, in-memory repository, xUnit tests.

## Structure (key files)

- `Program.cs` — ASP.NET Core app startup and DI registration.
- `Controllers/MarkerController.cs` — HTTP endpoints under `api/v1/markers`.
- `Services/MarkerService.cs` — business logic; maps DTOs to domain model.
- `Repositories/MarkerRepository.cs` — `InMemoryMarkerRepository` implementation used for tests and local runs.
- `Models/Marker.cs` — domain model and DTOs (`CreateMarkerRequest`, `MarkerResponse`, `MarkerCategory` enum).
- `tests/MarkerServiceStandalone.spec.ts` — Playwright API tests.
- `MarkerServiceStandalone.Tests` — xUnit integration/unit tests.

## Running locally

1. Build:
   - `dotnet build`
2. Run the API (default listens on `http://localhost:59817` in test configs):
   - `dotnet run --project ./MarkerServiceStandalone` (or run from solution root)
3. Run tests:
   - Unit/integration: `dotnet test`
   - Playwright (API tests): ensure API is running or add Playwright `webServer` to `playwright.config.ts` so the runner starts it.

## API surface (routes)

- `POST /api/v1/markers` — create marker. Requires header `X-User-Id: <userId>`.
- `GET /api/v1/markers/{markerId}` — get marker by id.
- `GET /api/v1/markers?category={n}` — list markers by optional category.
- `PATCH /api/v1/markers/{markerId}` — update marker.
- `DELETE /api/v1/markers/{markerId}` — delete marker.

Request/response DTOs are in `Models/Marker.cs`.

## Important implementation notes

- `MarkerCategory` is an enum. The project accepts string names in incoming JSON for creation by decorating `CreateMarkerRequest.Category` with `[JsonConverter(typeof(JsonStringEnumConverter))]`. Responses may be numeric or string depending on JSON options — see `Program.cs`.

- `InMemoryMarkerRepository` stores markers keyed by `userId` and `marker.Id`. Make sure `MarkerService.CreateMarkerAsync` sets `marker.UserId = userId` when creating markers so subsequent GET/DELETE calls with the same `X-User-Id` can find them.

- `Program` contains top-level statements. A `public partial class Program { }` is present to allow test hosts to locate the `Program` type for integration tests.

## Debugging tips

- To debug tests that use `WebApplicationFactory<Program>`, ensure `Program` is discoverable (the `partial` class stub exists) and the test host can construct the app.

## Known fixes included in the repo

- Added per-property `JsonStringEnumConverter` for `CreateMarkerRequest.Category` to accept string names in incoming JSON.

## Contributing / extending

- To persist data beyond process lifetime, replace `InMemoryMarkerRepository` with a DB-backed implementation and register it in `Program.cs`.
- Consider adding a `Health` endpoint (e.g. `/health`) to improve readiness checks used by CI.
