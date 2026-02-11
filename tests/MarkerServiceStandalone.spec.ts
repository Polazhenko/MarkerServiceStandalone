import { test, expect } from '@playwright/test';

const baseURL = process.env.BASE_URL ?? 'http://localhost:59817';

test.describe('MarkerServiceStandalone API (end-to-end)', () => {
  test('create ? get ? delete ? not found', async ({ request }) => {
    const headers = { 'X-User-Id': '42' };

    // Create marker
    const createPayload = {
      name: 'Playwright E2E Marker',
      category: 1, // MarkerCategory.General
      latitude: 45.0,
      longitude: -93.0,
      description: 'Created by Playwright test'
    };

    const post = await request.post(`${baseURL}/api/v1/markers`, {
      data: createPayload,
      headers
    });

    expect(post.status()).toBe(201);
    const created = await post.json();
    expect(created).toHaveProperty('id');
    expect(created.name).toBe(createPayload.name);

    const markerId = created.id as string;

    // Get marker
    const get = await request.get(`${baseURL}/api/v1/markers/${encodeURIComponent(markerId)}`, { headers });
    expect(get.status()).toBe(200);
    const fetched = await get.json();
    expect(fetched.id).toBe(markerId);

    // Delete marker
    const del = await request.delete(`${baseURL}/api/v1/markers/${encodeURIComponent(markerId)}`, { headers });
    expect(del.status()).toBe(204);

    // Verify deletion - should return 404
    const getAfterDelete = await request.get(`${baseURL}/api/v1/markers/${encodeURIComponent(markerId)}`, { headers });
    expect(getAfterDelete.status()).toBe(404);
  });
});
