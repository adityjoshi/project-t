#!/bin/bash
# Create simple placeholder icons using ImageMagick or sips (macOS)

for size in 16 48 128; do
  if command -v convert &> /dev/null; then
    convert -size ${size}x${size} xc:'#4f46e5' icon${size}.png
  elif command -v sips &> /dev/null; then
    # macOS sips method - create a simple colored image
    python3 << PYEOF
from PIL import Image
img = Image.new('RGB', ($size, $size), color='#4f46e5')
img.save('icon${size}.png')
PYEOF
  else
    echo "No image tool found. Please install ImageMagick or PIL"
  fi
done
