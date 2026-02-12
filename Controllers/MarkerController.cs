using Microsoft.AspNetCore.Mvc;
using MarkerServiceStandalone.Models;
using MarkerServiceStandalone.Services;
using System.Collections.Generic;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using System;

namespace MarkerServiceStandalone.Controllers;

[ApiController]
[Route("api/v1/markers")]
public class MarkerController : ControllerBase
{
    private readonly IMarkerService _markerService;
    private readonly ILogger<MarkerController> _logger;

    public MarkerController(IMarkerService markerService, ILogger<MarkerController> logger)
    {
        _markerService = markerService;
        _logger = logger;
    }

    [HttpPost]
    public async Task<ActionResult<MarkerResponse>> CreateMarker(
        [FromHeader(Name = "X-User-Id")] int userId,
        [FromBody] CreateMarkerRequest request)
    {
        try
        {
            var result = await _markerService.CreateMarkerAsync(userId, request);
            return CreatedAtAction(nameof(GetMarker), new { userId, markerId = result.Id }, result);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error creating marker for user {UserId}", userId);
            return BadRequest(new { error = ex.Message });
        }
    }

    [HttpPost("search")]
    public async Task<ActionResult<MarkerResponse>> SearchMarkers(
        [FromHeader(Name = "X-User-Id")] int userId,
        [FromBody] SearchMarkerRequest request)
    {
        try
        {
            if (string.IsNullOrWhiteSpace(request.Name))
                return BadRequest(new { error = "Name field should contain substring to search" });

            var result = await _markerService.SearchMarkersAsync(userId, request);
            return Ok(result);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error creating marker for user {UserId}", userId);
            return BadRequest(new { error = ex.Message });
        }
    }


    [HttpGet("{markerId}")]
    public async Task<ActionResult<MarkerResponse>> GetMarker(
        [FromHeader(Name = "X-User-Id")] int userId,
        [FromRoute] string markerId)
    {
        var marker = await _markerService.GetMarkerAsync(userId, markerId);
        if (marker == null)
            return NotFound(new { error = "Marker not found" });

        return Ok(marker);
    }

    [HttpGet]
    public async Task<ActionResult<ICollection<MarkerResponse>>> GetAllMarkers(
        [FromHeader(Name = "X-User-Id")] int userId,
        [FromQuery] MarkerCategory? category = null)
    {
        var markers = await _markerService.GetAllMarkersAsync(userId, category);
        return Ok(markers);
    }

    [HttpPatch("{markerId}")]
    public async Task<ActionResult<MarkerResponse>> UpdateMarker(
        [FromHeader(Name = "X-User-Id")] int userId,
        [FromRoute] string markerId,
        [FromBody] UpdateMarkerRequest request)
    {
        var result = await _markerService.UpdateMarkerAsync(userId, markerId, request);
        if (result == null)
            return NotFound(new { error = "Marker not found" });

        return Ok(result);
    }

    [HttpDelete("{markerId}")]
    public async Task<IActionResult> DeleteMarker(
        [FromHeader(Name = "X-User-Id")] int userId,
        [FromRoute] string markerId)
    {
        var result = await _markerService.DeleteMarkerAsync(userId, markerId);
        if (!result)
            return NotFound(new { error = "Marker not found" });

        return NoContent();
    }
}
