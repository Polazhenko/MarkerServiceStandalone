using MarkerServiceStandalone.Models;
using MarkerServiceStandalone.Repositories;
using MarkerServiceStandalone.Services;
using Moq;
using AutoMapper;
using Microsoft.Extensions.Logging;

namespace MarkerServiceStandalone.Tests;

public class MarkerServiceTests
{
    private readonly Mock<IMarkerRepository> _mockRepository;
    private readonly Mock<IMapper> _mockMapper;
    private readonly Mock<ILogger<MarkerService>> _mockLogger;
    private readonly MarkerService _service;

    public MarkerServiceTests()
    {
        _mockRepository = new Mock<IMarkerRepository>();
        _mockMapper = new Mock<IMapper>();
        _mockLogger = new Mock<ILogger<MarkerService>>();
        _service = new MarkerService(_mockRepository.Object, _mockMapper.Object, _mockLogger.Object);
    }

    [Fact]
    public async Task CreateMarkerAsync_ShouldReturnMarkerResponse_WhenValidRequest()
    {
        // Arrange
        var request = new CreateMarkerRequest
        {
            Name = "Test Marker",
            Category = MarkerCategory.General,
            Latitude = 45.5f,
            Longitude = -93.2f,
            Description = "Test description"
        };
        var userId = 1;
        var expectedResponse = new MarkerResponse
        {
            Id = "General_v1.0_123",
            Name = request.Name,
            Category = request.Category,
            Latitude = request.Latitude,
            Longitude = request.Longitude,
            Description = request.Description
        };

        _mockRepository.Setup(r => r.CreateAsync(It.IsAny<Marker>()))
            .ReturnsAsync((Marker m) => m);
        _mockMapper.Setup(m => m.Map<MarkerResponse>(It.IsAny<Marker>()))
            .Returns(expectedResponse);

        // Act
        var result = await _service.CreateMarkerAsync(userId, request);

        // Assert
        Assert.NotNull(result);
        Assert.Equal(request.Name, result.Name);
        _mockRepository.Verify(r => r.CreateAsync(It.IsAny<Marker>()), Times.Once);
    }

    [Fact]
    public async Task GetMarkerAsync_ShouldReturnMarkerResponse_WhenExists()
    {
        // Arrange
        var markerId = "test-id";
        var userId = 1;
        var marker = new Marker { Id = markerId, UserId = userId, Name = "Test" };
        var expectedResponse = new MarkerResponse { Id = markerId, Name = "Test" };

        _mockRepository.Setup(r => r.GetByIdAsync(markerId, userId))
            .ReturnsAsync(marker);
        _mockMapper.Setup(m => m.Map<MarkerResponse>(marker))
            .Returns(expectedResponse);

        // Act
        var result = await _service.GetMarkerAsync(userId, markerId);

        // Assert
        Assert.NotNull(result);
        Assert.Equal(markerId, result.Id);
    }

    [Fact]
    public async Task GetMarkerAsync_ShouldReturnNull_WhenNotExists()
    {
        // Arrange
        _mockRepository.Setup(r => r.GetByIdAsync(It.IsAny<string>(), It.IsAny<int>()))
            .ReturnsAsync((Marker?)null);

        // Act
        var result = await _service.GetMarkerAsync(1, "non-existent");

        // Assert
        Assert.Null(result);
    }

    [Fact]
    public async Task UpdateMarkerAsync_ShouldUpdateFields_WhenMarkerExists()
    {
        // Arrange
        var markerId = "test-id";
        var userId = 1;
        var existingMarker = new Marker 
        { 
            Id = markerId, 
            UserId = userId, 
            Name = "Old Name",
            Description = "Old Description"
        };
        var updateRequest = new UpdateMarkerRequest
        {
            Name = "New Name",
            Description = "New Description"
        };
        var expectedResponse = new MarkerResponse
        {
            Id = markerId,
            Name = "New Name",
            Description = "New Description"
        };

        _mockRepository.Setup(r => r.GetByIdAsync(markerId, userId))
            .ReturnsAsync(existingMarker);
        _mockRepository.Setup(r => r.UpdateAsync(It.IsAny<Marker>()))
            .ReturnsAsync((Marker m) => m);
        _mockMapper.Setup(m => m.Map<MarkerResponse>(It.IsAny<Marker>()))
            .Returns(expectedResponse);

        // Act
        var result = await _service.UpdateMarkerAsync(userId, markerId, updateRequest);

        // Assert
        Assert.NotNull(result);
        Assert.Equal("New Name", result.Name);
        Assert.Equal("New Description", result.Description);
    }

    [Fact]
    public async Task DeleteMarkerAsync_ShouldReturnTrue_WhenMarkerExists()
    {
        // Arrange
        var markerId = "test-id";
        var userId = 1;

        _mockRepository.Setup(r => r.DeleteAsync(markerId, userId))
            .ReturnsAsync(true);

        // Act
        var result = await _service.DeleteMarkerAsync(userId, markerId);

        // Assert
        Assert.True(result);
    }
}