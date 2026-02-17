import { describe, it, expect } from 'vitest';
import { isRawUrl, isYouTube, normalizeUrl } from './url';

describe('URL Utilities', () => {
    describe('isRawUrl', () => {
        it('identifies full https URLs', () => {
            expect(isRawUrl('https://example.com')).toBe(true);
            expect(isRawUrl('http://test.org/path?query=1')).toBe(true);
        });

        it('identifies URLs without protocol', () => {
            expect(isRawUrl('example.com')).toBe(true);
            expect(isRawUrl('sub.domain.co.uk/page')).toBe(true);
        });

        it('rejects non-URL text', () => {
            expect(isRawUrl('Click here')).toBe(false);
            expect(isRawUrl('My Website')).toBe(false);
            expect(isRawUrl('Visit google.com for more')).toBe(false);
        });

        it('handles whitespace', () => {
            expect(isRawUrl('  https://example.com  ')).toBe(true);
        });
    });

    describe('isYouTube', () => {
        it('identifies various YouTube formats', () => {
            expect(isYouTube('https://www.youtube.com/watch?v=dQw4w9WgXcQ')).toBe(true);
            expect(isYouTube('https://youtu.be/dQw4w9WgXcQ')).toBe(true);
            expect(isYouTube('https://youtube.com/shorts/12345678901')).toBe(true);
            expect(isYouTube('https://youtube.com/live/12345678901')).toBe(true);
        });

        it('rejects non-YouTube URLs', () => {
            expect(isYouTube('https://google.com')).toBe(false);
            expect(isYouTube('https://vimeo.com/123')).toBe(false);
        });
    });

    describe('normalizeUrl', () => {
        it('removes protocol and www', () => {
            expect(normalizeUrl('https://www.example.com/')).toBe('example.com');
            expect(normalizeUrl('http://example.com')).toBe('example.com');
        });
    });
});
