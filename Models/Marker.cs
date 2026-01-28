using System;

namespace MarkerServiceStandalone.Models;

public enum MarkerCategory
{
    General = 1,
    Category2 = 2,
    Category3 = 3,
    Category4 = 4
}

public class Marker
{
    public string Id { get; set; } = string.Empty;
    public int UserId { get; set; }
    public string Name { get; set; } = string.Empty;
    public MarkerCategory Category { get; set; }
    public float Latitude { get; set; }
    public float Longitude { get; set; }
    public string? Description { get; set; }
    public DateTimeOffset CreatedAt { get; set; }
    public DateTimeOffset UpdatedAt { get; set; }
}

public record CreateMarkerRequest
{
    public string Name { get; set; } = string.Empty;
    public MarkerCategory Category { get; set; }
    public float Latitude { get; set; }
    public float Longitude { get; set; }
    public string? Description { get; set; }
}

public record UpdateMarkerRequest
{
    public string? Name { get; set; }
    public MarkerCategory? Category { get; set; }
    public float? Latitude { get; set; }
    public float? Longitude { get; set; }
    public string? Description { get; set; }
}

public record MarkerResponse
{
    public string Id { get; set; } = string.Empty;
    public string Name { get; set; } = string.Empty;
    public MarkerCategory Category { get; set; }
    public float Latitude { get; set; }
    public float Longitude { get; set; }
    public string? Description { get; set; }
    public DateTimeOffset CreatedAt { get; set; }
    public DateTimeOffset UpdatedAt { get; set; }
}
