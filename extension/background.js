// Background service worker for Synapse extension

// Create context menu on installation
chrome.runtime.onInstalled.addListener(() => {
  chrome.contextMenus.create({
    id: 'saveSelectedText',
    title: 'Save selected text to Synapse',
    contexts: ['selection']
  });

  chrome.contextMenus.create({
    id: 'savePage',
    title: 'Save this page to Synapse',
    contexts: ['page']
  });
});

// Handle context menu clicks
chrome.contextMenus.onClicked.addListener(async (info, tab) => {
  try {
    let data = {
      url: tab.url,
      title: tab.title
    };

    if (info.menuItemId === 'saveSelectedText' && info.selectionText) {
      data.selectedText = info.selectionText.trim();
      data.content = info.selectionText.trim();
    }

    // Store data for popup to pick up
    await chrome.storage.local.set({
      quickSave: {
        title: data.title,
        content: data.selectedText || '',
        url: data.url,
        timestamp: Date.now()
      }
    });

    // Show badge to indicate data is ready
    chrome.action.setBadgeText({ text: '1', tabId: tab.id });
    chrome.action.setBadgeBackgroundColor({ color: '#4f46e5' });
  } catch (error) {
    console.error('Error in context menu:', error);
  }
});
