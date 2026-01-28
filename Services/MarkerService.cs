using AutoMapper;
using MarkerServiceStandalone.Models;
using MarkerServiceStandalone.Repositories;
using Microsoft.Extensions.Logging;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace MarkerServiceStandalone.Services;

public interface IMarkerService
{
    Task<MarkerResponse> CreateMarkerAsync(int userId, CreateMarkerRequest request);
    Task<MarkerResponse?> GetMarkerAsync(int userId, string markerId);
    Task<IEnumerable<MarkerResponse>> GetAllMarkersAsync(int userId, MarkerCategory? category = null);
    Task<MarkerResponse?> UpdateMarkerAsync(int userId, string markerId, UpdateMarkerRequest request);
    Task<bool> DeleteMarkerAsync(int userId, string markerId);
}

public class MarkerService : IMarkerService
{
    private readonly IMarkerRepository _repository;
    private readonly IMapper _mapper;
    private readonly ILogger<MarkerService> _logger;

    public MarkerService(IMarkerRepository repository, IMapper mapper, ILogger<MarkerService> logger)
    {
        _repository = repository;
        _mapper = mapper;
        _logger = logger;
    }

    public async Task<MarkerResponse> CreateMarkerAsync(int userId, CreateMarkerRequest request)
    {
        var marker = new Marker
        {
            Id = GenerateMarkerId(request.Category),
            UserId = userId,
            Name = request.Name,
            Category = request.Category,
            Latitude = request.Latitude,
            Longitude = request.Longitude,
            Description = request.Description,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        await _repository.CreateAsync(marker);
        _logger.LogInformation("Created marker {MarkerId} for user {UserId}", marker.Id, userId);

        return _mapper.Map<MarkerResponse>(marker);
    }

    public async Task<MarkerResponse?> GetMarkerAsync(int userId, string markerId)
    {
        var marker = await _repository.GetByIdAsync(markerId, userId);
        return marker != null ? _mapper.Map<MarkerResponse>(marker) : null;
    }

    public async Task<IEnumerable<MarkerResponse>> GetAllMarkersAsync(int userId, MarkerCategory? category = null)
    {
        var markers = await _repository.GetByCategoryAsync(userId, category);
        return markers.Select(m => _mapper.Map<MarkerResponse>(m));
    }

    public async Task<MarkerResponse?> UpdateMarkerAsync(int userId, string markerId, UpdateMarkerRequest request)
    {
        var existing = await _repository.GetByIdAsync(markerId, userId);
        if (existing == null) return null;

        if (request.Name != null) existing.Name = request.Name;
        if (request.Category.HasValue) existing.Category = request.Category.Value;
        if (request.Latitude.HasValue) existing.Latitude = request.Latitude.Value;
        if (request.Longitude.HasValue) existing.Longitude = request.Longitude.Value;
        if (request.Description != null) existing.Description = request.Description;

        var updated = await _repository.UpdateAsync(existing);
        _logger.LogInformation("Updated marker {MarkerId} for user {UserId}", markerId, userId);

        return updated != null ? _mapper.Map<MarkerResponse>(updated) : null;
    }

    public async Task<bool> DeleteMarkerAsync(int userId, string markerId)
    {
        var result = await _repository.DeleteAsync(markerId, userId);
        if (result)
            _logger.LogInformation("Deleted marker {MarkerId} for user {UserId}", markerId, userId);
        return result;
    }

    private static string GenerateMarkerId(MarkerCategory category)
    {
        return $"{category}_v1.0_{DateTimeOffset.UtcNow.ToUnixTimeSeconds()}";
    }
}
