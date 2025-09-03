import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { ErrorProvider } from './components/UserInfo/Error/Error.jsx'
import { InformationProvider } from './components/UserInfo/information/InformationField.jsx'
import App from './App.jsx'

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <ErrorProvider>
      <InformationProvider>
        <App />
      </InformationProvider>
    </ErrorProvider>
  </StrictMode>,
)
