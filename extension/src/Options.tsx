import React, { useEffect, useState } from 'react';
import './Popup.css'; // Reuse basic styles
import { getSettings, saveSettings, Settings, DEFAULT_SETTINGS, FilterMode } from './utils/settings';

const Options: React.FC = () => {
    const [settings, setSettings] = useState<Settings>(DEFAULT_SETTINGS);
    const [newDomain, setNewDomain] = useState('');

    useEffect(() => {
        getSettings().then(setSettings);
    }, []);

    const handleSave = async (updated: Partial<Settings>) => {
        const newSettings = { ...settings, ...updated };
        await saveSettings(newSettings);
        setSettings(newSettings);
    };

    const addDomain = () => {
        const domain = newDomain.trim().toLowerCase();
        if (domain && !settings.domainList.includes(domain)) {
            handleSave({ domainList: [...settings.domainList, domain] });
            setNewDomain('');
        }
    };

    const removeDomain = (domain: string) => {
        handleSave({ domainList: settings.domainList.filter(d => d !== domain) });
    };

    return (
        <div className="container" style={{ width: '400px', margin: '20px auto' }}>
            <div className="header">
                <h1>LinkLens Settings</h1>
            </div>

            <div className="content">
                <div className="card">
                    <div>
                        <div className="label-text">Filtering Mode</div>
                        <div className="status" style={{ fontSize: '12px' }}>
                            {settings.filterMode === 'blocklist' 
                                ? 'Resolving all EXCEPT these domains' 
                                : 'ONLY resolving these domains'}
                        </div>
                    </div>
                    <select 
                        value={settings.filterMode} 
                        onChange={(e) => handleSave({ filterMode: e.target.value as FilterMode })}
                        style={{ padding: '4px' }}
                    >
                        <option value="blocklist">Blocklist</option>
                        <option value="allowlist">Allowlist</option>
                    </select>
                </div>

                <div className="card" style={{ flexDirection: 'column', alignItems: 'flex-start' }}>
                    <div className="label-text">Domain List</div>
                    <div style={{ display: 'flex', width: '100%', marginTop: '8px' }}>
                        <input 
                            type="text" 
                            placeholder="example.com" 
                            value={newDomain}
                            onChange={(e) => setNewDomain(e.target.value)}
                            onKeyDown={(e) => e.key === 'Enter' && addDomain()}
                            style={{ flex: 1, padding: '8px', border: '1px solid #ccc', borderRadius: '4px 0 0 4px' }}
                        />
                        <button 
                            onClick={addDomain}
                            className="btn"
                            style={{ margin: 0, borderRadius: '0 4px 4px 0', padding: '8px 16px' }}
                        >
                            Add
                        </button>
                    </div>
                    <ul style={{ width: '100%', padding: 0, marginTop: '12px', listStyle: 'none' }}>
                        {settings.domainList.map(domain => (
                            <li key={domain} style={{ display: 'flex', justifyContent: 'space-between', padding: '8px 0', borderBottom: '1px solid #eee' }}>
                                <span>{domain}</span>
                                <button 
                                    onClick={() => removeDomain(domain)}
                                    style={{ border: 'none', background: 'none', color: '#ff4444', cursor: 'pointer' }}
                                >
                                    Remove
                                </button>
                            </li>
                        ))}
                    </ul>
                </div>

                <div className="card">
                    <div>
                        <div className="label-text">Match Subdomains</div>
                        <div className="status" style={{ fontSize: '12px' }}>e.g. corp.com also blocks dev.corp.com</div>
                    </div>
                    <label className="switch">
                        <input 
                            type="checkbox" 
                            checked={settings.matchSubdomains} 
                            onChange={(e) => handleSave({ matchSubdomains: e.target.checked })} 
                        />
                        <span className="slider"></span>
                    </label>
                </div>
            </div>

            <div className="footer">
                Privacy settings are stored locally in your browser.
            </div>
        </div>
    );
};

export default Options;
