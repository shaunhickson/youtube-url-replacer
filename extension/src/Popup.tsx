import { useState, useEffect } from 'react';
import './Popup.css';

function Popup() {
  const [enabled, setEnabled] = useState(true);

  // Load state from storage on mount
  useEffect(() => {
    chrome.storage.local.get(['enabled'], (result) => {
      if (result.enabled !== undefined) {
        setEnabled(result.enabled);
      }
    });
  }, []);

  const toggleSwitch = (checked: boolean) => {
    setEnabled(checked);
    chrome.storage.local.set({ enabled: checked });
    // Reload current tab to apply changes immediately? 
    // Or send message to content script. For now, simple storage set is enough.
  };

  const handleReload = () => {
    chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
      if (tabs[0].id) {
        chrome.tabs.reload(tabs[0].id);
      }
    });
  };

  return (
    <div className="container">
      <div className="header">
        <img src="/icons/icon.svg" alt="Logo" style={{ width: '24px', height: '24px', marginRight: '8px' }} />
        <h1>YouTube Replacer</h1>
      </div>
      
      <div className="content">
        <div className="card">
          <div>
            <div className="label-text">Extension Status</div>
            <div className="status">{enabled ? 'Active' : 'Disabled'}</div>
          </div>
          <label className="switch">
            <input 
              type="checkbox" 
              checked={enabled} 
              onChange={(e) => toggleSwitch(e.target.checked)} 
            />
            <span className="slider"></span>
          </label>
        </div>

        <div style={{ textAlign: 'center' }}>
          <button className="btn" onClick={handleReload}>
            Reload Page
          </button>
        </div>
      </div>

      <div className="footer">
        Powered by Google Cloud Run & Firestore
      </div>
    </div>
  );
}

export default Popup;