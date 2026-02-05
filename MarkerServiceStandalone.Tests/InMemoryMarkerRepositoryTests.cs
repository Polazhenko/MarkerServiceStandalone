using MarkerServiceStandalone.Models;
using MarkerServiceStandalone.Repositories;

namespace MarkerServiceStandalone.Tests;

public class InMemoryMarkerRepositoryTests
{
    private readonly InMemoryMarkerRepository _repository;

    public InMemoryMarkerRepositoryTests()
    {
        _repository = new InMemoryMarkerRepository();
    }

    [Fact]
    public async Task CreateAsync_ShouldStoreMarker()
    {
        // Arrange
        var marker = new Marker
        {
            Id = "test-id",
            UserId = 1,
            Name = "Test Marker",
            Category = MarkerCategory.General,
            Latitude = 45.5f,
            Longitude = -93.2f,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        // Act
        var result = await _repository.CreateAsync(marker);

        // Assert
        Assert.NotNull(result);
        Assert.Equal(marker.Id, result.Id);
        Assert.Equal(marker.Name, result.Name);
    }

    [Fact]
    public async Task GetByIdAsync_ShouldReturnMarker_WhenExists()
    {
        // Arrange
        var marker = new Marker { UserId = 1, Name = "Test" };
        var created = await _repository.CreateAsync(marker);

        // Act
        var result = await _repository.GetByIdAsync(created.Id, 1);

        // Assert
        Assert.NotNull(result);
        Assert.Equal(created.Id, result.Id);
    }

    [Fact]
    public async Task GetByIdAsync_ShouldReturnNull_WhenWrongUser()
    {
        // Arrange
        var marker = new Marker { UserId = 1, Name = "Test" };
        var created = await _repository.CreateAsync(marker);

        // Act
        var result = await _repository.GetByIdAsync(created.Id, 2);

        // Assert
        Assert.Null(result);
    }

    [Fact]
    public async Task GetAllAsync_ShouldReturnUserMarkers_FilteredByCategory()
    {
        // Arrange
        await _repository.CreateAsync(new Marker { UserId = 1, Id = "General_v1.0_123", Category = MarkerCategory.General });
        await _repository.CreateAsync(new Marker { UserId = 1, Id = "GroupCategory1_Marker2_v1.0_1343", Category = MarkerCategory.GroupCategory1_Marker2 });
        await _repository.CreateAsync(new Marker { UserId = 2, Id = "General_v1.0_12323", Category = MarkerCategory.General });

        // Act
        var result = await _repository.GetByCategoryAsync(1, MarkerCategory.General);

        // Assert
        Assert.Single(result);
        Assert.All(result, m => Assert.Equal(MarkerCategory.General, m.Category));
    }

    [Fact]
    public async Task UpdateAsync_ShouldModifyMarker_AndUpdateTimestamp()
    {
        // Arrange
        var marker = new Marker { UserId = 1, Name = "Original" };
        var created = await _repository.CreateAsync(marker);
        var originalUpdatedAt = created.UpdatedAt;
        
        await Task.Delay(1); // Ensure timestamp difference
        created.Name = "Updated";

        // Act
        var result = await _repository.UpdateAsync(created);

        // Assert
        Assert.NotNull(result);
        Assert.Equal("Updated", result.Name);
        Assert.True(result.UpdatedAt > originalUpdatedAt);
    }

    [Fact]
    public async Task DeleteAsync_ShouldReturnTrue_WhenMarkerExists()
    {
        // Arrange
        var marker = new Marker { UserId = 1, Name = "Test" };
        var created = await _repository.CreateAsync(marker);

        // Act
        var result = await _repository.DeleteAsync(created.Id, 1);

        // Assert
        Assert.True(result);
        
        // Verify deletion
        var deleted = await _repository.GetByIdAsync(created.Id, 1);
        Assert.Null(deleted);
    }

    [Fact]
    public async Task DeleteAsync_ShouldReturnFalse_WhenMarkerNotExists()
    {
        // Act
        var result = await _repository.DeleteAsync("non-existent", 1);

        // Assert
        Assert.False(result);
    }
}