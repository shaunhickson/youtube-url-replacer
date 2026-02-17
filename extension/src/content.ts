import { isRawUrl, isYouTube } from './utils/url';

console.log("LinkLens Content Script Loaded");

interface PendingLink {
    element: HTMLAnchorElement;
    url: string;
    isYT: boolean;
}

class LinkLensOptimizer {
    private processedLinks = new Set<string>();
    private pendingResolutions = new Set<string>();
    private scanScheduled = false;
    private batchTimeout: ReturnType<typeof setTimeout> | null = null;
    private linksToProcess: PendingLink[] = [];

    constructor() {
        this.init();
    }

    private init() {
        // Initial scan of the whole body
        this.enqueueLinks(Array.from(document.querySelectorAll('a')));

        // Setup MutationObserver for targeted scans
        const observer = new MutationObserver((mutations) => {
            const newLinks: HTMLAnchorElement[] = [];
            
            for (const mutation of mutations) {
                for (const node of Array.from(mutation.addedNodes)) {
                    if (node instanceof HTMLElement) {
                        // If it's a link, add it
                        if (node instanceof HTMLAnchorElement) {
                            newLinks.push(node);
                        }
                        // Also find links inside the added node
                        const nestedLinks = node.querySelectorAll('a');
                        nestedLinks.forEach(l => newLinks.push(l));
                    }
                }
            }

            if (newLinks.length > 0) {
                this.enqueueLinks(newLinks);
            }
        });

        observer.observe(document.body, {
            childList: true,
            subtree: true
        });
    }

    private enqueueLinks(links: HTMLAnchorElement[]) {
        let addedAny = false;
        
        links.forEach(link => {
            const href = link.href;
            if (this.processedLinks.has(href) || this.pendingResolutions.has(href)) {
                return;
            }

            const text = link.innerText.trim();
            if (isRawUrl(text)) {
                this.linksToProcess.push({
                    element: link,
                    url: href,
                    isYT: isYouTube(href)
                });
                this.pendingResolutions.add(href);
                addedAny = true;
            }
        });

        if (addedAny) {
            this.scheduleBatch();
        }
    }

    private scheduleBatch() {
        if (this.scanScheduled) return;
        this.scanScheduled = true;

        const runScan = () => {
            if (this.batchTimeout) clearTimeout(this.batchTimeout);
            
            this.batchTimeout = setTimeout(() => {
                this.processBatch();
                this.scanScheduled = false;
            }, 500); // Batch every 500ms
        };

        if ('requestIdleCallback' in window) {
            window.requestIdleCallback(() => runScan(), { timeout: 2000 });
        } else {
            setTimeout(runScan, 100);
        }
    }

    private async processBatch() {
        if (this.linksToProcess.length === 0) return;

        const currentBatch = [...this.linksToProcess];
        this.linksToProcess = [];

        const urlsToFetch = [...new Set(currentBatch.map(l => l.url))];

        try {
            const response = await fetch('https://youtube-replacer-backend-542312799814.us-east1.run.app/resolve', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ urls: urlsToFetch })
            });

            if (!response.ok) {
                throw new Error(`Backend returned ${response.status}`);
            }

            const data = await response.json();
            const titles = data.titles || {};

            currentBatch.forEach(({ element, url, isYT }) => {
                this.pendingResolutions.delete(url);
                this.processedLinks.add(url);

                if (titles[url]) {
                    const title = titles[url];
                    const prefix = isYT ? "[YT] " : "";
                    element.innerText = `${prefix}${title}`;
                    element.title = title;
                }
            });

        } catch (err) {
            console.error('LinkLens: Error resolving batch:', err);
            // Cleanup pending so they can be retried in next scan if mutation occurs
            currentBatch.forEach(({ url }) => this.pendingResolutions.delete(url));
        }
    }
}

// Start the optimizer
new LinkLensOptimizer();
