import requests
import json

BASE_URL = "https://localhost:5001/api/v1/markers"
USER_ID = "1"
HEADERS = {"X-User-Id": USER_ID, "Content-Type": "application/json"}

def create_marker():
    """Create a new marker"""
    data = {
        "name": "Deer Sighting",
        "category": 2,
        "latitude": 45.5231,
        "longitude": -93.2467,
        "description": "Large buck spotted near the oak tree"
    }
    response = requests.post(BASE_URL, json=data, headers=HEADERS, verify=False)
    print(f"Create marker: {response.status_code}")
    if response.status_code == 201:
        marker = response.json()
        print(f"Created marker ID: {marker['id']}")
        return marker['id']
    return None

def get_all_markers():
    """Get all markers"""
    response = requests.get(BASE_URL, headers=HEADERS, verify=False)
    print(f"Get all markers: {response.status_code}")
    if response.status_code == 200:
        markers = response.json()
        print(f"Found {len(markers)} markers")
        return markers
    return []

def get_marker_by_id(marker_id):
    """Get specific marker by ID"""
    response = requests.get(f"{BASE_URL}/{marker_id}", headers=HEADERS, verify=False)
    print(f"Get marker {marker_id}: {response.status_code}")
    return response.json() if response.status_code == 200 else None

def update_marker(marker_id):
    """Update marker"""
    data = {
        "name": "Updated Deer Sighting",
        "description": "Large 8-point buck spotted at dawn"
    }
    response = requests.patch(f"{BASE_URL}/{marker_id}", json=data, headers=HEADERS, verify=False)
    print(f"Update marker {marker_id}: {response.status_code}")
    return response.json() if response.status_code == 200 else None

def delete_marker(marker_id):
    """Delete marker"""
    response = requests.delete(f"{BASE_URL}/{marker_id}", headers=HEADERS, verify=False)
    print(f"Delete marker {marker_id}: {response.status_code}")

if __name__ == "__main__":
    # Create marker
    marker_id = create_marker()
    
    if marker_id:
        # Get all markers
        get_all_markers()
        
        # Get specific marker
        get_marker_by_id(marker_id)
        
        # Update marker
        update_marker(marker_id)
        
        # Delete marker
        delete_marker(marker_id)