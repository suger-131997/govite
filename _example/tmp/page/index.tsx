
import { StrictMode } from "react";
import { createRoot } from 'react-dom/client'
import App from '~/page/index.tsx'
import '~/global.css'


createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App {...(window.APP_PROPS || {})}/>
  </StrictMode>
)
