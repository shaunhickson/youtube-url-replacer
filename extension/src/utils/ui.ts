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
        
        // Inject Styles
        const style = document.createElement('style');
        style.textContent = `
            .tooltip {
                position: absolute;
                z-index: 1000000;
                background: #ffffff;
                color: #333333;
                padding: 12px;
                border-radius: 8px;
                box-shadow: 0 4px 12px rgba(0,0,0,0.15);
                font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
                font-size: 14px;
                line-height: 1.4;
                width: 280px;
                pointer-events: none;
                opacity: 0;
                transition: opacity 0.2s ease-in-out;
                border: 1px solid #eeeeee;
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
                color: #888888;
                margin-bottom: 4px;
                display: flex;
                align-items: center;
            }
            .title {
                font-weight: 600;
                color: #000000;
                margin-bottom: 6px;
                display: -webkit-box;
                -webkit-line-clamp: 2;
                -webkit-box-orient: vertical;
                overflow: hidden;
            }
            .description {
                font-size: 13px;
                color: #666666;
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

    public show(target: HTMLElement, data: TooltipData) {
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
    generic: `<svg viewBox="0 0 24 24" fill="currentColor" width="1em" height="1em" style="vertical-align: middle; margin-right: 4px;"><path d="M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 1.99 2H18c1.1 0 2-.9 2-2V8l-6-6zm2 16H8v-2h8v2zm0-4H8v-2h8v2zm-3-5V3.5L18.5 9H13z"/></svg>`,
    link: `<svg viewBox="0 0 24 24" fill="currentColor" width="1em" height="1em" style="vertical-align: middle; margin-right: 4px;"><path d="M3.9 12c0-1.71 1.39-3.1 3.1-3.1h4V7H7c-2.76 0-5 2.24-5 5s2.24 5 5 5h4v-1.9H7c-1.71 0-3.1-1.39-3.1-3.1zM8 13h8v-2H8v2zm9-6h-4v1.9h4c1.71 0 3.1 1.39 3.1 3.1s-1.39 3.1-3.1 3.1h-4V17h4c2.76 0 5-2.24 5-5s-2.24-5-5-5z"/></svg>`
};
