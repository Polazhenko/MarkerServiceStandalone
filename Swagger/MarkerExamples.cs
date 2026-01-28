using MarkerServiceStandalone.Models;
using Swashbuckle.AspNetCore.Filters;
using System;
using System.Collections.Generic;

namespace MarkerServiceStandalone.Swagger;

public class CreateMarkerRequestExample : IExamplesProvider<CreateMarkerRequest>
{
    public CreateMarkerRequest GetExamples()
    {
        return new ()
        {
            Name = "Deer Sighting",
            Category = MarkerCategory.Category2,
            Latitude = 45.5231f,
            Longitude = -93.2467f,
            Description = "Large buck spotted near the oak tree"
        };
    }
}

public class UpdateMarkerRequestExample : IExamplesProvider<UpdateMarkerRequest>
{
    public UpdateMarkerRequest GetExamples()
    {
        return new ()
        {
            Name = "Updated Deer Sighting",
            Description = "Large 8-point buck spotted near the oak tree at dawn"
        };
    }
}

public class MarkerResponseExample : IExamplesProvider<MarkerResponse>
{
    public MarkerResponse GetExamples()
    {
        return new ()
        {
            Id = "marker_123456789",
            Name = "Deer Sighting",
            Category = MarkerCategory.Category2,
            Latitude = 45.5231f,
            Longitude = -93.2467f,
            Description = "Large buck spotted near the oak tree",
            CreatedAt = DateTimeOffset.UtcNow.AddDays(-1),
            UpdatedAt = DateTimeOffset.UtcNow
        };
    }
}

public class MarkerListResponseExample : IExamplesProvider<List<MarkerResponse>>
{
    public List<MarkerResponse> GetExamples()
    {
        return
        [
            new MarkerResponse
            {
                Id = "marker_123456789",
                Name = "Deer Sighting",
                Category = MarkerCategory.Category2,
                Latitude = 45.5231f,
                Longitude = -93.2467f,
                Description = "Large buck spotted near the oak tree",
                CreatedAt = DateTime.UtcNow.AddDays(-2),
                UpdatedAt = DateTime.UtcNow.AddDays(-1)
            },
            new MarkerResponse
            {
                Id = "marker_987654321",
                Name = "Weather Station",
                Category = MarkerCategory.Category4,
                Latitude = 45.5189f,
                Longitude = -93.2501f,
                Description = "Temperature and humidity monitoring station",
                CreatedAt = DateTime.UtcNow.AddDays(-5),
                UpdatedAt = DateTime.UtcNow.AddDays(-3)
            }
        ];
    }
}