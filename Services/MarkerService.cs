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
    Task<IEnumerable<MarkerResponse>> SearchMarkersAsync(int userId, SearchMarkerRequest request);
    Task<MarkerResponse?> GetMarkerAsync(int userId, string markerId);
    Task<IEnumerable<MarkerResponse>> GetAllMarkersAsync(int userId, MarkerCategory? category = null);
    Task<MarkerResponse?> UpdateMarkerAsync(int userId, string markerId, UpdateMarkerRequest request);
    Task<bool> DeleteMarkerAsync(int userId, string markerId);
}

public class MarkerService(IMarkerRepository repository, IMapper mapper, ILogger<MarkerService> logger) : IMarkerService
{
    private readonly IMarkerRepository _repository = repository;
    private readonly IMapper _mapper = mapper;
    private readonly ILogger<MarkerService> _logger = logger;

    public async Task<MarkerResponse> CreateMarkerAsync(int userId, CreateMarkerRequest request)
    {
        var marker = _mapper.Map<Marker>(request);

        marker.Id = GenerateMarkerId(request.Category);
        marker.CreatedAt = DateTime.UtcNow;
        // Ensure the marker is associated with the creating user
        marker.UserId = userId;

        await _repository.CreateAsync(marker);
        _logger.LogInformation("Created marker {MarkerId} for user {UserId}", marker.Id, userId);

        return _mapper.Map<MarkerResponse>(marker);
    }

    public async Task<IEnumerable<MarkerResponse>> SearchMarkersAsync(int userId, SearchMarkerRequest request)
    {
        return _mapper.Map<IEnumerable<MarkerResponse>>(await _repository.SearchAsync(userId, request));
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
        // category cannot be changed
        // if (request.Category.HasValue) existing.Category = request.Category.Value;
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
