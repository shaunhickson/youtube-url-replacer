console.log("Background script running");

// Placeholder for background logic (e.g., API calls, caching)
chrome.runtime.onMessage.addListener((request, _sender, sendResponse) => {
  if (request.type === 'GET_TITLE') {
    // TODO: Fetch title from backend
    sendResponse({ title: "Placeholder Title" });
  }
});
