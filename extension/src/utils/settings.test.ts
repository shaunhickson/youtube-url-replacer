import { describe, it, expect } from 'vitest';
import { isDomainAllowed, Settings } from './settings';

describe('Settings Utilities', () => {
    const baseSettings: Settings = {
        enabled: true,
        filterMode: 'blocklist',
        domainList: ['example.com'],
        matchSubdomains: true,
    };

    describe('isDomainAllowed', () => {
        it('blocks domains in blocklist mode', () => {
            expect(isDomainAllowed('example.com', baseSettings)).toBe(false);
            expect(isDomainAllowed('google.com', baseSettings)).toBe(true);
        });

        it('blocks subdomains if enabled', () => {
            expect(isDomainAllowed('sub.example.com', baseSettings)).toBe(false);
        });

        it('does not block subdomains if disabled', () => {
            const noSubSettings = { ...baseSettings, matchSubdomains: false };
            expect(isDomainAllowed('sub.example.com', noSubSettings)).toBe(true);
        });

        it('allows only listed domains in allowlist mode', () => {
            const allowSettings: Settings = {
                ...baseSettings,
                filterMode: 'allowlist',
                domainList: ['trusted.com'],
            };
            expect(isDomainAllowed('trusted.com', allowSettings)).toBe(true);
            expect(isDomainAllowed('sub.trusted.com', allowSettings)).toBe(true);
            expect(isDomainAllowed('example.com', allowSettings)).toBe(false);
        });

        it('returns false if extension is disabled', () => {
            const disabledSettings = { ...baseSettings, enabled: false };
            expect(isDomainAllowed('any.com', disabledSettings)).toBe(false);
        });
        
        it('is case-insensitive', () => {
            expect(isDomainAllowed('EXAMPLE.COM', baseSettings)).toBe(false);
        });
    });
});
