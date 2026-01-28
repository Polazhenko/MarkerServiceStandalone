using MarkerServiceStandalone.Models;
using System;
using System.Collections.Generic;
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
    private readonly List<Marker> _markers = new();
    private readonly object _lock = new();

    public Task<Marker?> GetByIdAsync(string id, int userId)
    {
        lock (_lock)
        {
            return Task.FromResult(_markers.FirstOrDefault(m => m.Id == id && m.UserId == userId));
        }
    }

    public Task<IEnumerable<Marker>> GetByUserIdAsync(int userId)
    {
        lock (_lock)
        {
            return Task.FromResult(_markers.Where(m => m.UserId == userId).AsEnumerable());
        }
    }

    public Task<IEnumerable<Marker>> GetByCategoryAsync(int userId, MarkerCategory? category)
    {
        lock (_lock)
        {
            var query = _markers.Where(m => m.UserId == userId);
            if (category.HasValue)
                query = query.Where(m => m.Category == category.Value);
            return Task.FromResult(query.AsEnumerable());
        }
    }

    public Task<Marker> CreateAsync(Marker marker)
    {
        lock (_lock)
        {
            _markers.Add(marker);
            return Task.FromResult(marker);
        }
    }

    public Task<Marker?> UpdateAsync(Marker marker)
    {
        lock (_lock)
        {
            var existing = _markers.FirstOrDefault(m => m.Id == marker.Id && m.UserId == marker.UserId);
            if (existing == null) return Task.FromResult<Marker?>(null);

            existing.Name = marker.Name;
            existing.Category = marker.Category;
            existing.Latitude = marker.Latitude;
            existing.Longitude = marker.Longitude;
            existing.Description = marker.Description;
            existing.UpdatedAt = DateTime.UtcNow;

            return Task.FromResult<Marker?>(existing);
        }
    }

    public Task<bool> DeleteAsync(string id, int userId)
    {
        lock (_lock)
        {
            var marker = _markers.FirstOrDefault(m => m.Id == id && m.UserId == userId);
            if (marker == null) return Task.FromResult(false);
            _markers.Remove(marker);
            return Task.FromResult(true);
        }
    }
}
