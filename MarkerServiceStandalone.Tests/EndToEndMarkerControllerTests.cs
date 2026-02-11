using System.Net;
using System.Net.Http.Json;
using Microsoft.AspNetCore.Mvc.Testing;
using MarkerServiceStandalone.Models;

namespace MarkerServiceStandalone.Tests;

public class EndToEndMarkerControllerTests : IAsyncLifetime
{
    private readonly WebApplicationFactory<Program> _factory;
    private HttpClient? _client;

    public EndToEndMarkerControllerTests()
    {
        _factory = new WebApplicationFactory<Program>();
    }

    public Task InitializeAsync()
    {
        _client = _factory.CreateClient();
        return Task.CompletedTask;
    }

    public Task DisposeAsync()
    {
        _client?.Dispose();
        _factory.Dispose();
        return Task.CompletedTask;
    }

    [Fact]
    public async Task Create_Get_Delete_Flow_Works()
    {
        // Arrange
        var createPayload = new CreateMarkerRequest
        {
            Name = "E2E Test",
            Category = MarkerCategory.General,
            Latitude = 45.0f,
            Longitude = -93.0f,
            Description = "End to end test marker"
        };

        _client!.DefaultRequestHeaders.Remove("X-User-Id");
        _client.DefaultRequestHeaders.Add("X-User-Id", "42");

        // Act - Create (use HttpClient.BaseAddress to build full request URI)
        var createUri = new Uri(_client!.BaseAddress!, "/api/v1/markers");
        var postResponse = await _client.PostAsJsonAsync(createUri, createPayload);
        Assert.Equal(HttpStatusCode.Created, postResponse.StatusCode);

        var created = await postResponse.Content.ReadFromJsonAsync<MarkerResponse>();
        Assert.NotNull(created);
        Assert.Equal("E2E Test", created!.Name);

        // Act - Get
        var getUri = new Uri(_client.BaseAddress!, $"/api/v1/markers/{Uri.EscapeDataString(created.Id)}");
        var getResponse = await _client.GetAsync(getUri);
        Assert.Equal(HttpStatusCode.OK, getResponse.StatusCode);
        var fetched = await getResponse.Content.ReadFromJsonAsync<MarkerResponse>();
        Assert.Equal(created.Id, fetched!.Id);

        // Act - Delete
        var deleteUri = new Uri(_client.BaseAddress!, $"/api/v1/markers/{Uri.EscapeDataString(created.Id)}");
        var deleteResponse = await _client.DeleteAsync(deleteUri);
        Assert.Equal(HttpStatusCode.NoContent, deleteResponse.StatusCode);

        // Verify deletion
        var getAfterDelete = await _client.GetAsync(getUri);
        Assert.Equal(HttpStatusCode.NotFound, getAfterDelete.StatusCode);
    }
}
