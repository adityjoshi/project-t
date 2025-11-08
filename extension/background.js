// Background service worker for the extension

chrome.runtime.onInstalled.addListener(() => {
  console.log('Synapse extension installed');
});

