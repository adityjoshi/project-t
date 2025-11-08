#!/bin/bash

# Script to refresh images for all items with broken or missing images

echo "Refreshing images for items with broken or missing images..."

# Get all items
ITEMS=$(curl -s "http://localhost:8080/api/items")

# Extract item IDs that need image refresh
echo "$ITEMS" | python3 << 'PYTHON'
import sys
import json
import subprocess

items = json.load(sys.stdin)
fixed = 0
failed = 0

for item in items:
    image_url = item.get('image_url', '')
    item_id = item['id']
    
    # Check if image is broken (Unsplash Source) or missing
    needs_fix = False
    if not image_url:
        needs_fix = True
    elif 'source.unsplash.com' in image_url:
        needs_fix = True
    
    if needs_fix:
        print(f"Refreshing image for: {item['title'][:50]}...")
        result = subprocess.run(
            ['curl', '-s', '-X', 'POST', f'http://localhost:8080/api/items/{item_id}/refresh-image'],
            capture_output=True,
            text=True
        )
        if result.returncode == 0:
            print(f"  ✓ Fixed")
            fixed += 1
        else:
            print(f"  ✗ Failed: {result.stderr}")
            failed += 1

print(f"\nDone! Fixed: {fixed}, Failed: {failed}")
PYTHON

echo "Image refresh complete!"

