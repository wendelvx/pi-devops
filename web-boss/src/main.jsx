import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.jsx'

// Sem StrictMode, o Socket conecta apenas UMA vez de forma limpa.
createRoot(document.getElementById('root')).render(
    <App />
)