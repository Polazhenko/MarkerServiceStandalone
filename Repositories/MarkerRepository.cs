using MarkerServiceStandalone.Models;
using System;
using System.Collections.Generic;
using System.Collections.Concurrent;
using System.Linq;
using System.Threading.Tasks;

namespace MarkerServiceStandalone.Repositories;

public interface IMarkerRepository
{
    Task<Marker?> GetByIdAsync(string id, int userId);
    Task<IEnumerable<Marker>> GetByUserIdAsync(int userId);
    Task<IEnumerable<Marker>> GetByCategoryAsync(int userId, MarkerCategory? category);
    Task<Marker> CreateAsync(Marker marker);
    Task<Marker?> UpdateAsync(Marker marker);
    Task<bool> DeleteAsync(string id, int userId);
}

public class InMemoryMarkerRepository : IMarkerRepository
{
    private readonly ConcurrentDictionary<int, ConcurrentDictionary<string, Marker>> _markers = new();
    private static readonly IEnumerable<Marker> EmptyMarkerList = new List<Marker>().AsEnumerable();

    public Task<Marker?> GetByIdAsync(string id, int userId)
    {
        Marker? retVal = null;

        if (_markers.TryGetValue(userId, out var userMarkers) &&
            userMarkers.TryGetValue(id, out retVal))
        {
            return Task.FromResult(retVal);
        }

        return Task.FromResult(retVal);
    }

    public Task<IEnumerable<Marker>> GetByUserIdAsync(int userId)
    {
        return GetByCategoryAsync(userId, null);
    }

    public Task<IEnumerable<Marker>> GetByCategoryAsync(int userId, MarkerCategory? category)
    {
        if (!_markers.TryGetValue(userId, out var userMarkers))
            return Task.FromResult(EmptyMarkerList);

        if (category == null)
            return Task.FromResult(userMarkers.Values.AsEnumerable());

        return Task.FromResult(userMarkers.Values.Where(m => m.Category == category).ToArray().AsEnumerable());
    }

    public Task<Marker> CreateAsync(Marker marker)
    {
        if (!_markers.TryGetValue(marker.UserId, out var userMarkers))
        {
            userMarkers = new();
            _markers[marker.UserId] = userMarkers; 
        }

        userMarkers[marker.Id] = marker; 
        return Task.FromResult(marker);
    }

    public Task<Marker?> UpdateAsync(Marker marker)
    {
        var existing = GetByIdAsync(marker.Id, marker.UserId).Result;
        if (existing == null) return Task.FromResult<Marker?>(null);

        existing.Name = marker.Name;
        existing.Category = marker.Category;
        existing.Latitude = marker.Latitude;
        existing.Longitude = marker.Longitude;
        existing.Description = marker.Description;
        existing.UpdatedAt = DateTime.UtcNow;

        return Task.FromResult<Marker?>(existing);
    }

    public Task<bool> DeleteAsync(string id, int userId)
    {
        if (_markers.TryGetValue(userId, out var userMarkers))
        {
            // ConcurrentDictionary does not have Remove; use TryRemove
            return Task.FromResult(userMarkers.TryRemove(id, out var _));
        }

        return Task.FromResult(false);
    }
}
