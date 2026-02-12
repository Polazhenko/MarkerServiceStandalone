using System;
using System.Text.Json.Serialization;

namespace MarkerServiceStandalone.Models;


public enum MarkerCategory
{
    All,
    General = 1,
    GroupCategory1_Marker1 = 2,
    GroupCategory1_Marker2 = 3,
    GroupCategory2_Marker1 = 4,
}

public class Marker
{
    public string Id { get; set; } = string.Empty;
    public int UserId { get; set; }
    public string Name { get; set; } = string.Empty;
    public MarkerCategory Category { get; set; } = MarkerCategory.General;
    public float Latitude { get; set; }
    public float Longitude { get; set; }
    public string? Description { get; set; }
    public DateTimeOffset CreatedAt { get; set; } = DateTimeOffset.Now;
    public DateTimeOffset UpdatedAt { get; set; }
}

public record CreateMarkerRequest
{
    public string Name { get; set; } = string.Empty;

    [JsonConverter(typeof(JsonStringEnumConverter))]
    public MarkerCategory Category { get; set; }

    public float Latitude { get; set; }
    public float Longitude { get; set; }
    public string? Description { get; set; }
}

public record UpdateMarkerRequest
{
    public string? Name { get; set; }

    [JsonConverter(typeof(JsonStringEnumConverter))]
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
