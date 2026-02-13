console.log("YouTube URL Replacer Content Script Loaded");

// Regex to extract video ID from various YouTube URL formats
const YOUTUBE_REGEX = /(?:youtube\.com\/(?:[^\/]+\/.+\/|(?:v|e(?:mbed)?|live|shorts)\/|.*[?&]v=)|youtu\.be\/)([^"&?\/\s]{11})/i;

// Cache processed links to avoid re-fetching
const processedLinks = new Set<string>();

async function scanAndReplace() {
    const links = Array.from(document.querySelectorAll('a'));
    const pendingLinks: { element: HTMLAnchorElement; videoId: string }[] = [];
    const videoIdsToFetch: string[] = [];

    // 1. Identify valid, unprocessed YouTube links
    links.forEach(link => {
        // Skip if already processed
        if (processedLinks.has(link.href)) return;

        // Check if the HREF is a YouTube link
        const hrefMatch = link.href.match(YOUTUBE_REGEX);
        if (hrefMatch && hrefMatch[1]) {
            const videoId = hrefMatch[1];
            
            // NEW: Only proceed if the visible text also contains a YouTube URL
            if (YOUTUBE_REGEX.test(link.innerText)) {
                pendingLinks.push({ element: link, videoId });
                videoIdsToFetch.push(videoId);
                processedLinks.add(link.href); // Mark as seen
            }
        }
    });

    if (videoIdsToFetch.length === 0) return;

    // 2. Fetch titles from backend
    try {
        const uniqueIds = [...new Set(videoIdsToFetch)];
        const response = await fetch('https://youtube-replacer-backend-542312799814.us-east1.run.app/resolve', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ videoIds: uniqueIds })
        });

        if (!response.ok) {
            console.error('Failed to fetch titles:', response.statusText);
            return;
        }

        const data = await response.json();
        const titles = data.titles || {};

        // 3. Update DOM
        pendingLinks.forEach(({ element, videoId }) => {
            if (titles[videoId]) {
                // Determine if we should replace innerText or add a tooltip
                // For now, let's append the title or replace if it's just a raw URL
                
                // If text is the URL itself or "click here" etc, replace it. 
                // Otherwise maybe append? Let's just replace for the MVP as requested.
                // You can refine this logic later.
                element.innerText = `[YT] ${titles[videoId]}`;
                element.title = titles[videoId]; // Also set tooltip
            }
        });

    } catch (err) {
        console.error('Error resolving YouTube titles:', err);
    }
}

// Initial scan
scanAndReplace();

// Observe DOM for infinite scrolling / dynamic content
let timeout: any;
const observer = new MutationObserver(() => {
    // Debounce to avoid spamming on every tiny DOM change
    clearTimeout(timeout);
    timeout = setTimeout(scanAndReplace, 1000);
});

observer.observe(document.body, { childList: true, subtree: true });