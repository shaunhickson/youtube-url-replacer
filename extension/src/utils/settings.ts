export type FilterMode = 'blocklist' | 'allowlist';

export interface Settings {
    enabled: boolean;
    filterMode: FilterMode;
    domainList: string[];
    matchSubdomains: boolean;
}

export const DEFAULT_SETTINGS: Settings = {
    enabled: true,
    filterMode: 'blocklist',
    domainList: [],
    matchSubdomains: true,
};

export async function getSettings(): Promise<Settings> {
    return new Promise((resolve) => {
        chrome.storage.local.get(DEFAULT_SETTINGS, (items) => {
            resolve(items as Settings);
        });
    });
}

export async function saveSettings(settings: Partial<Settings>): Promise<void> {
    return new Promise((resolve) => {
        chrome.storage.local.set(settings, () => {
            resolve();
        });
    });
}

/**
 * Checks if a given domain should be processed based on the user's settings.
 */
export function isDomainAllowed(domain: string, settings: Settings): boolean {
    if (!settings.enabled) return false;

    const normalizedDomain = domain.toLowerCase();
    let isMatch = false;

    for (const item of settings.domainList) {
        const normalizedItem = item.toLowerCase();
        if (normalizedDomain === normalizedItem) {
            isMatch = true;
            break;
        }
        if (settings.matchSubdomains && normalizedDomain.endsWith('.' + normalizedItem)) {
            isMatch = true;
            break;
        }
    }

    if (settings.filterMode === 'blocklist') {
        return !isMatch;
    } else {
        return isMatch;
    }
}

/**
 * Extracts the domain from a URL string.
 */
export function getDomain(url: string): string {
    try {
        const u = new URL(url);
        return u.hostname;
    } catch (e) {
        return '';
    }
}
