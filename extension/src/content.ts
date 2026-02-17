import { isRawUrl, isYouTube } from './utils/url';

console.log("LinkLens Content Script Loaded");

// Cache processed links to avoid re-fetching
const processedLinks = new Set<string>();

async function scanAndReplace() {
    const links = Array.from(document.querySelectorAll('a'));
    const pendingLinks: { element: HTMLAnchorElement; url: string; isYT: boolean }[] = [];
    const urlsToFetch: string[] = [];

    // 1. Identify valid, unprocessed "raw" links
    links.forEach(link => {
        // Skip if already processed
        if (processedLinks.has(link.href)) return;

        const href = link.href;
        const text = link.innerText.trim();

        // Heuristic: If text looks like a URL, it's a candidate
        if (isRawUrl(text)) {
            const isYT = isYouTube(href);
            pendingLinks.push({ element: link, url: href, isYT });
            urlsToFetch.push(href);
            processedLinks.add(href); // Mark as seen
        }
    });

    if (urlsToFetch.length === 0) return;

    // 2. Fetch titles from backend
    try {
        const uniqueUrls = [...new Set(urlsToFetch)];
        const response = await fetch('https://youtube-replacer-backend-542312799814.us-east1.run.app/resolve', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ urls: uniqueUrls })
        });

        if (!response.ok) {
            console.error('Failed to fetch titles:', response.statusText);
            return;
        }

        const data = await response.json();
        const titles = data.titles || {};

        // 3. Update DOM
        pendingLinks.forEach(({ element, url, isYT }) => {
            if (titles[url]) {
                const title = titles[url];
                const prefix = isYT ? "[YT] " : "";
                element.innerText = `${prefix}${title}`;
                element.title = title; // Also set tooltip
            }
        });

    } catch (err) {
        console.error('Error resolving titles:', err);
    }
}

// Initial scan
scanAndReplace();

// Observe DOM for infinite scrolling / dynamic content
let timeout: ReturnType<typeof setTimeout> | undefined;
const observer = new MutationObserver(() => {
    // Debounce to avoid spamming on every tiny DOM change
    clearTimeout(timeout);
    timeout = setTimeout(scanAndReplace, 1000);
});

observer.observe(document.body, { childList: true, subtree: true });