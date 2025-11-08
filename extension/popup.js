const API_BASE_URL = 'http://localhost:8080/api';

document.addEventListener('DOMContentLoaded', async () => {
  const form = document.getElementById('captureForm');
  const titleInput = document.getElementById('title');
  const contentInput = document.getElementById('content');
  const fillPageBtn = document.getElementById('fillPage');
  const screenshotBtn = document.getElementById('screenshotBtn');
  const saveBtn = document.getElementById('saveBtn');
  const statusDiv = document.getElementById('status');
  const selectedTextDiv = document.getElementById('selectedText');
  const contentTypeDiv = document.getElementById('contentType');

  // Get current tab info
  const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
  
  // Check for quick save data from context menu
  const quickSave = await chrome.storage.local.get('quickSave');
  if (quickSave.quickSave) {
    titleInput.value = quickSave.quickSave.title || tab.title || '';
    contentInput.value = quickSave.quickSave.content || '';
    // Clear the quick save data and badge
    chrome.storage.local.remove('quickSave');
    chrome.action.setBadgeText({ text: '', tabId: tab.id });
    // Auto-fill if content exists
    if (quickSave.quickSave.content) {
      showStatus('Content ready to save!', 'success');
    }
  } else {
    // Set default title
    titleInput.value = tab.title || '';
  }

  // Detect and show content type
  chrome.tabs.sendMessage(tab.id, { action: 'extractContent' }, (response) => {
    if (response) {
      if (response.type) {
        contentTypeDiv.textContent = `Type: ${response.type.charAt(0).toUpperCase() + response.type.slice(1)}`;
        contentTypeDiv.style.display = 'block';
      }
    }
  });

  // Check for selected text
  chrome.tabs.sendMessage(tab.id, { action: 'getSelectedText' }, (response) => {
    if (response && response.selectedText) {
      selectedTextDiv.textContent = `Selected: "${response.selectedText.substring(0, 100)}..."`;
      selectedTextDiv.style.display = 'block';
      contentInput.value = response.selectedText;
    }
  });

  // Fill from page button
  fillPageBtn.addEventListener('click', async () => {
    try {
      fillPageBtn.disabled = true;
      fillPageBtn.textContent = 'Extracting...';
      showStatus('Extracting content...', 'info');
      
      const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
      
      // Use promise-based message sending
      const response = await new Promise((resolve, reject) => {
        chrome.tabs.sendMessage(tab.id, { action: 'extractContent' }, (response) => {
          if (chrome.runtime.lastError) {
            reject(new Error(chrome.runtime.lastError.message));
          } else {
            resolve(response);
          }
        });
      });
      
      if (response) {
        if (response.title) titleInput.value = response.title;
        if (response.content) {
          contentInput.value = response.content;
        }
        
        // Show content type
        if (response.type) {
          contentTypeDiv.textContent = `Type: ${response.type.charAt(0).toUpperCase() + response.type.slice(1)}`;
          contentTypeDiv.style.display = 'block';
        }
        
        showStatus('Page content extracted!', 'success');
      } else {
        showStatus('No content found on this page', 'error');
      }
    } catch (error) {
      console.error('Error extracting content:', error);
      showStatus(`Failed to extract content: ${error.message}`, 'error');
    } finally {
      fillPageBtn.disabled = false;
      fillPageBtn.textContent = '📄 Fill from Page';
    }
  });

  // Screenshot button
  screenshotBtn.addEventListener('click', async () => {
    try {
      const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
      
      screenshotBtn.disabled = true;
      screenshotBtn.textContent = 'Capturing...';
      showStatus('Capturing screenshot...', 'info');

      // Capture visible tab
      const dataUrl = await chrome.tabs.captureVisibleTab(null, {
        format: 'png',
        quality: 100
      });

      // Convert to blob and then to base64 for API
      const response = await fetch(dataUrl);
      const blob = await response.blob();
      const reader = new FileReader();
      
      reader.onloadend = async () => {
        const base64data = reader.result.split(',')[1];
        
        // Save screenshot
        const saveResponse = await fetch(`${API_BASE_URL}/items`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            title: tab.title || 'Screenshot',
            content: `Screenshot captured from: ${tab.url}`,
            source_url: tab.url,
            type: 'image',
            image_url: dataUrl, // Store as data URL
          }),
        });

        if (!saveResponse.ok) {
          throw new Error('Failed to save screenshot');
        }

        showStatus('Screenshot saved! ✓', 'success');
        setTimeout(() => {
          window.close();
        }, 1000);
      };
      
      reader.readAsDataURL(blob);
    } catch (error) {
      console.error('Error capturing screenshot:', error);
      showStatus('Failed to capture screenshot', 'error');
    } finally {
      screenshotBtn.disabled = false;
      screenshotBtn.textContent = '📸 Screenshot';
    }
  });

  // Form submission
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const title = titleInput.value.trim();
    const content = contentInput.value.trim();
    const url = tab.url;

    if (!title && !content) {
      showStatus('Please enter a title or content', 'error');
      return;
    }

    saveBtn.disabled = true;
    saveBtn.textContent = 'Saving...';
    showStatus('Saving...', 'info');

    try {
      // Get enhanced content with metadata
      const extractResponse = await new Promise((resolve) => {
        chrome.tabs.sendMessage(tab.id, { action: 'extractContent' }, resolve);
      });

      const payload = {
        title: title || extractResponse?.title || 'Untitled',
        content: content || extractResponse?.content || title,
        source_url: url,
        type: extractResponse?.type || 'url',
      };

      // Add metadata if available (important for video descriptions, product info, etc.)
      if (extractResponse?.metadata) {
        payload.metadata = extractResponse.metadata;
        // Include image from metadata if available
        if (extractResponse.metadata.image) {
          payload.image_url = extractResponse.metadata.image;
        }
        if (extractResponse.metadata.thumbnail) {
          payload.image_url = extractResponse.metadata.thumbnail;
        }
      }

      const response = await fetch(`${API_BASE_URL}/items`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Failed to save');
      }

      showStatus('Saved successfully! ✓', 'success');
      
      // Clear form after a delay
      setTimeout(() => {
        titleInput.value = '';
        contentInput.value = '';
        window.close();
      }, 1000);
    } catch (error) {
      console.error('Error saving:', error);
      showStatus(`Failed to save: ${error.message}`, 'error');
    } finally {
      saveBtn.disabled = false;
      saveBtn.textContent = 'Save';
    }
  });
});

function showStatus(message, type) {
  const statusDiv = document.getElementById('status');
  statusDiv.textContent = message;
  statusDiv.className = `status ${type}`;
  statusDiv.style.display = 'block';
  
  if (type === 'success') {
    setTimeout(() => {
      statusDiv.style.display = 'none';
    }, 3000);
  }
}
