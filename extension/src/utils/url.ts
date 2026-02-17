// Regex to identify if text looks like a raw URL
export const RAW_URL_REGEX = /^(https?:\/\/)?([\w.-]+)\.([a-z]{2,})(:\d+)?(\/\S*)?$/i;

// Regex to extract video ID from various YouTube URL formats
// eslint-disable-next-line no-useless-escape
export const YOUTUBE_REGEX = /(?:youtube\.com\/(?:[^\/]+\/.+\/|(?:v|e(?:mbed)?|live|shorts)\/|.*[?&]v=)|youtu\.be\/)([^"&?\/\s]{11})/i;

/**
 * Checks if a string looks like a raw URL that should be resolved.
 */
export function isRawUrl(text: string): boolean {
    return RAW_URL_REGEX.test(text.trim());
}

/**
 * Checks if a URL is a YouTube URL.
 */
export function isYouTube(url: string): boolean {
    return YOUTUBE_REGEX.test(url);
}

/**
 * Normalizes a URL for comparison or display (optional helper)
 */
export function normalizeUrl(url: string): string {
    return url.replace(/^(https?:\/\/)?(www\.)?/, '').replace(/\/$/, '');
}
