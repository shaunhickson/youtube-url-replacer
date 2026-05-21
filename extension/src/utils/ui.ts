/**
 * UI Manager for LinkLens
 * Handles Shadow DOM isolation and rich tooltip rendering
 */

const TOOLTIP_ID = 'link-lens-tooltip-root';

export interface TooltipData {
    title: string;
    description?: string;
    domain: string;
    platform: string;
}

class UIManager {
    private host: HTMLDivElement | null = null;
    private shadow: ShadowRoot | null = null;
    private tooltip: HTMLDivElement | null = null;

    constructor() {
        if (typeof document === 'undefined') return;
        this.createShadowRoot();
    }

    private createShadowRoot() {
        if (document.getElementById(TOOLTIP_ID)) return;

        this.host = document.createElement('div');
        this.host.id = TOOLTIP_ID;
        document.body.appendChild(this.host);

        this.shadow = this.host.attachShadow({ mode: 'closed' });
        
        // Inject Styles with CSS variables for theming
        const style = document.createElement('style');
        style.textContent = `
            :host {
                --ll-bg: #ffffff;
                --ll-text: #333333;
                --ll-header: #888888;
                --ll-title: #000000;
                --ll-border: #eeeeee;
                --ll-shadow: rgba(0,0,0,0.15);
            }

            .tooltip.dark {
                --ll-bg: #1e1e1e;
                --ll-text: #cccccc;
                --ll-header: #aaaaaa;
                --ll-title: #ffffff;
                --ll-border: #333333;
                --ll-shadow: rgba(0,0,0,0.5);
            }

            @media (prefers-color-scheme: dark) {
                .tooltip.system {
                    --ll-bg: #1e1e1e;
                    --ll-text: #cccccc;
                    --ll-header: #aaaaaa;
                    --ll-title: #ffffff;
                    --ll-border: #333333;
                    --ll-shadow: rgba(0,0,0,0.5);
                }
            }

            .tooltip {
                position: absolute;
                z-index: 1000000;
                background: var(--ll-bg);
                color: var(--ll-text);
                padding: 12px;
                border-radius: 8px;
                box-shadow: 0 4px 12px var(--ll-shadow);
                font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
                font-size: 14px;
                line-height: 1.4;
                width: 280px;
                pointer-events: none;
                opacity: 0;
                transition: opacity 0.2s ease-in-out;
                border: 1px solid var(--ll-border);
                visibility: hidden;
            }
            .tooltip.visible {
                opacity: 1;
                visibility: visible;
            }
            .header {
                font-size: 11px;
                text-transform: uppercase;
                letter-spacing: 0.5px;
                color: var(--ll-header);
                margin-bottom: 4px;
                display: flex;
                align-items: center;
            }
            .title {
                font-weight: 600;
                color: var(--ll-title);
                margin-bottom: 6px;
                display: -webkit-box;
                -webkit-line-clamp: 2;
                -webkit-box-orient: vertical;
                overflow: hidden;
            }
            .description {
                font-size: 13px;
                color: var(--ll-text);
                display: -webkit-box;
                -webkit-line-clamp: 3;
                -webkit-box-orient: vertical;
                overflow: hidden;
            }
            .platform-icon {
                width: 12px;
                height: 12px;
                margin-right: 4px;
            }
        `;
        this.shadow.appendChild(style);

        this.tooltip = document.createElement('div');
        this.tooltip.className = 'tooltip';
        this.shadow.appendChild(this.tooltip);
    }

    public show(target: HTMLElement, data: TooltipData, theme: string = 'system') {
        if (!this.tooltip) return;

        const rect = target.getBoundingClientRect();
        const scrollX = window.scrollX;
        const scrollY = window.scrollY;

        this.tooltip.innerHTML = `
            <div class="header">
                <span>${data.domain}</span>
            </div>
            <div class="title">${this.escape(data.title)}</div>
            ${data.description ? `<div class="description">${this.escape(data.description)}</div>` : ''}
        `;

        // Apply theme classes
        this.tooltip.classList.remove('light', 'dark', 'system');
        this.tooltip.classList.add(theme);

        // Position
        const top = rect.bottom + scrollY + 8;
        const left = rect.left + scrollX;

        this.tooltip.style.top = `${top}px`;
        this.tooltip.style.left = `${left}px`;
        this.tooltip.classList.add('visible');
    }

    public hide() {
        if (this.tooltip) {
            this.tooltip.classList.remove('visible');
        }
    }

    private escape(str: string): string {
        const div = document.createElement('div');
        div.textContent = str;
        return div.innerHTML;
    }
}

export const uiManager = new UIManager();

/**
 * Platform Icons as SVG Strings
 */
export const ICONS = {
    youtube: `<svg viewBox="0 0 24 24" fill="currentColor" width="1em" height="1em" style="vertical-align: middle; margin-right: 4px;"><path d="M23.498 6.186a3.016 3.016 0 0 0-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 0 0 .502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 0 0 2.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 0 0 2.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z"/></svg>`,
    twitter: `<svg viewBox="0 0 24 24" fill="currentColor" width="1em" height="1em" style="vertical-align: middle; margin-right: 4px;"><path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/></svg>`,
    linkedin: `<svg viewBox="0 0 24 24" fill="currentColor" width="1em" height="1em" style="vertical-align: middle; margin-right: 4px;"><path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/></svg>`,
    reddit: `<svg viewBox="0 0 24 24" fill="currentColor" width="1em" height="1em" style="vertical-align: middle; margin-right: 4px;"><path d="M12 22C6.477 22 2 17.523 2 12S6.477 2 12 2s10 4.477 10 10-4.477 10-10 10zm-3.23-5.26c-1.464 0-2.65-1.186-2.65-2.65 0-.58.188-1.116.51-1.554-1.353-.872-2.225-2.27-2.225-3.864 0-2.583 2.276-4.675 5.085-4.675 1.574 0 2.983.655 3.924 1.688l2.67-1.895.894 1.988-2.695 1.913c.094.343.144.707.144 1.082 0 2.582-2.276 4.674-5.086 4.674h-.57c0 1.463-1.187 2.65-2.652 2.65zm0-3.61c.53 0 .96.43.96.96s-.43.96-.96.96-.96-.43-.96-.96.43-.96.96-.96z"/></svg>`,
    wikipedia: `<svg viewBox="0 0 24 24" fill="currentColor" width="1em" height="1em" style="vertical-align: middle; margin-right: 4px;"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 14.5h-2v-2h2v2zm0-4h-2V7h2v5.5z"/></svg>`,
    spotify: `<svg viewBox="0 0 24 24" fill="currentColor" width="1em" height="1em" style="vertical-align: middle; margin-right: 4px;"><path d="M12 0C5.373 0 0 5.373 0 12s5.373 12 12 12 12-5.373 12-12S18.627 0 12 0zm5.496 17.316c-.23.364-.707.474-1.071.246-2.936-1.79-6.626-2.193-10.978-1.201-.412.094-.82-.165-.913-.575-.094-.412.165-.82.576-.913 4.757-1.08 8.825-.63 12.14 1.371.364.22.474.7.246 1.072zm1.488-3.29c-.292.476-.914.629-1.39.336-3.37-2.068-8.528-2.673-12.217-1.464-.542.176-1.119-.115-1.294-.658-.176-.543.115-1.119.658-1.294 4.248-1.391 10.026-.714 13.906 1.67.477.291.63.913.337 1.39zm.135-3.468c-4.045-2.4-10.741-2.62-14.622-1.448-.654.198-1.344-.173-1.542-.828-.198-.654.173-1.344.828-1.542 4.475-1.35 11.879-1.087 16.541 1.68.588.349.782 1.108.433 1.696-.349.588-1.108.782-1.638.442z"/></svg>`,
    generic: `<svg viewBox="0 0 24 24" fill="currentColor" width="1em" height="1em" style="vertical-align: middle; margin-right: 4px;"><path d="M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 1.99 2H18c1.1 0 2-.9 2-2V8l-6-6zm2 16H8v-2h8v2zm0-4H8v-2h8v2zm-3-5V3.5L18.5 9H13z"/></svg>`,
    link: `<svg viewBox="0 0 24 24" fill="currentColor" width="1em" height="1em" style="vertical-align: middle; margin-right: 4px;"><path d="M3.9 12c0-1.71 1.39-3.1 3.1-3.1h4V7H7c-2.76 0-5 2.24-5 5s2.24 5 5 5h4v-1.9H7c-1.71 0-3.1-1.39-3.1-3.1zM8 13h8v-2H8v2zm9-6h-4v1.9h4c1.71 0 3.1 1.39 3.1 3.1s-1.39 3.1-3.1 3.1h-4V17h4c2.76 0 5-2.24 5-5s-2.24-5-5-5z"/></svg>`
};
