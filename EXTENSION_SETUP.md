# Browser Extension Setup Guide

## Installation

1. **Open Chrome/Edge Extensions Page**
   - Chrome: Go to `chrome://extensions/`
   - Edge: Go to `edge://extensions/`

2. **Enable Developer Mode**
   - Toggle the "Developer mode" switch in the top right corner

3. **Load the Extension**
   - Click "Load unpacked"
   - Navigate to the `extension` folder in this project
   - Select the folder and click "Select Folder"

4. **Verify Installation**
   - You should see "Synapse - Capture Your Thoughts" in your extensions list
   - The extension icon should appear in your browser toolbar

## Creating Icons (Optional)

The extension needs icons to display properly. You can:

1. **Quick Option**: Create simple colored squares (16x16, 48x48, 128x128 pixels)
2. **Online Tool**: Use https://www.favicon-generator.org/ or similar
3. **Design Tool**: Create icons in Figma, Canva, or any image editor

Place the icons in `extension/icons/`:
- `icon16.png` (16x16 pixels)
- `icon48.png` (48x48 pixels)
- `icon128.png` (128x128 pixels)

The extension will work without icons, but Chrome will show a default placeholder.

## Usage

### Capturing Amazon Products

1. Navigate to any Amazon product page
2. Click the Synapse extension icon
3. The extension will automatically detect it's an Amazon product
4. Click "Fill from Page" to extract:
   - Product title
   - Price
   - Rating
   - Product image
   - Description
   - ASIN
5. Click "Save" to add to your knowledge base

### Capturing Blog Posts

1. Navigate to any blog post
2. Click the extension icon
3. Click "Fill from Page" to extract:
   - Blog title
   - Author
   - Publication date
   - Full content
   - Featured image
4. Click "Save"

### Capturing Videos

1. Navigate to YouTube, Vimeo, or any video page
2. Click the extension icon
3. Click "Fill from Page" to extract:
   - Video title
   - Channel/creator
   - Thumbnail
4. Click "Save"

### Taking Screenshots

1. Navigate to any webpage
2. Click the extension icon
3. Click "📸 Screenshot" button
4. The screenshot will be automatically saved to your knowledge base

### Capturing Selected Text

1. Highlight any text on a webpage
2. Click the extension icon
3. The selected text will be pre-filled in the content field
4. Add a title and click "Save"

### General Web Page Capture

1. Navigate to any webpage
2. Click the extension icon
3. Click "Fill from Page" to extract page content
4. Review and edit the extracted content
5. Click "Save"

## Features

- ✅ **Auto-detection**: Automatically detects content type (Amazon, Blog, Video, etc.)
- ✅ **Smart extraction**: Extracts relevant metadata for each content type
- ✅ **Screenshot capture**: Full-page screenshots
- ✅ **Selected text**: Captures highlighted text
- ✅ **Image extraction**: Automatically extracts product images, blog featured images, video thumbnails
- ✅ **Metadata capture**: Price, ratings, authors, dates, etc.

## Troubleshooting

### Extension not loading
- Make sure Developer Mode is enabled
- Check that all files are in the `extension` folder
- Look for errors in the Extensions page

### "Failed to save" error
- Make sure the backend is running (`docker-compose up`)
- Check that the API is accessible at `http://localhost:8080`
- Open browser console (F12) to see detailed error messages

### Content not extracting properly
- Some websites may have different HTML structures
- Try manually editing the extracted content
- The extension works best with standard website layouts

### Screenshot not working
- Make sure you've granted the extension necessary permissions
- Try refreshing the page and trying again
- Check browser console for errors

## Permissions Explained

- **activeTab**: Access to current tab content for extraction
- **storage**: Store extension settings
- **tabs**: Get tab information (title, URL)
- **desktopCapture**: Capture screenshots

All permissions are only used when you actively use the extension.

