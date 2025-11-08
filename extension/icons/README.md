# Extension Icons

Icons are optional for the extension to work. The extension will use Chrome's default icon if these files are not present.

## To Add Icons (Optional)

If you want custom icons, create these files:

- `icon16.png` (16x16 pixels)
- `icon48.png` (48x48 pixels)  
- `icon128.png` (128x128 pixels)

### Quick Ways to Create Icons:

1. **Online Icon Generator**: 
   - Visit https://www.favicon-generator.org/
   - Upload an image or create one
   - Download the generated icons

2. **Using the HTML Generator**:
   - Open `download_icons.html` in your browser
   - It will automatically download the icons

3. **Manual Creation**:
   - Use any image editor (Figma, Canva, Photoshop, etc.)
   - Create square images with the sizes above
   - Save as PNG format

4. **Using Python PIL** (if installed):
   ```bash
   pip install Pillow
   python3 -c "from PIL import Image; [Image.new('RGB', (s, s), '#4f46e5').save(f'icon{s}.png') for s in [16, 48, 128]]"
   ```

Once you have the icons, update `manifest.json` to include:
```json
"action": {
  "default_popup": "popup.html",
  "default_icon": {
    "16": "icons/icon16.png",
    "48": "icons/icon48.png",
    "128": "icons/icon128.png"
  }
},
"icons": {
  "16": "icons/icon16.png",
  "48": "icons/icon48.png",
  "128": "icons/icon128.png"
}
```
