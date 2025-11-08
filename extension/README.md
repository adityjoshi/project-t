# Synapse Browser Extension

A powerful browser extension to capture content from various websites to your Synapse knowledge base.

## Features

- **Amazon Product Capture**: Automatically extracts product title, price, rating, image, and description
- **Blog Post Capture**: Extracts blog title, author, date, content, and featured image
- **Video Capture**: Captures YouTube, Vimeo, and other video information with thumbnails
- **Screenshot Capture**: Takes full-page screenshots and saves them
- **Smart Content Detection**: Automatically detects content type and extracts relevant metadata
- **Selected Text**: Captures highlighted text from any page

## Installation

1. Open Chrome/Edge and go to `chrome://extensions/` (or `edge://extensions/`)
2. Enable "Developer mode" (toggle in top right)
3. Click "Load unpacked"
4. Select the `extension` folder from this project
5. The extension icon should appear in your toolbar

## Usage

1. Navigate to any webpage you want to capture
2. Click the Synapse extension icon
3. The popup will automatically detect the content type (Amazon, Blog, Video, etc.)
4. Click "Fill from Page" to auto-extract content
5. Or click "Screenshot" to capture the current page
6. Click "Save" to send to your Synapse knowledge base

## Supported Sites

- **Amazon**: Product pages (extracts price, rating, images, ASIN)
- **Blogs**: WordPress, Medium, and most blog platforms
- **YouTube/Vimeo**: Video information and thumbnails
- **Any Website**: General content extraction

## Permissions

- `activeTab`: To access current tab content
- `storage`: To store extension settings
- `tabs`: To get tab information
- `desktopCapture`: For screenshot functionality

## Development

The extension consists of:
- `manifest.json`: Extension configuration
- `popup.html/js`: Extension popup UI
- `content.js`: Content script that runs on web pages to extract data
- `background.js`: Background service worker

