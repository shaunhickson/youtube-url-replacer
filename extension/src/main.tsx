import React from 'react'
import ReactDOM from 'react-dom/client'
import Popup from './Popup.tsx'
import Options from './Options.tsx'

const rootElement = document.getElementById('root');
if (rootElement) {
  ReactDOM.createRoot(rootElement as HTMLElement).render(
    <React.StrictMode>
      <Popup />
    </React.StrictMode>,
  )
}

const optionsRoot = document.getElementById('options-root');
if (optionsRoot) {
  ReactDOM.createRoot(optionsRoot as HTMLElement).render(
    <React.StrictMode>
      <Options />
    </React.StrictMode>,
  )
}
