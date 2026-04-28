import { useState } from 'react'
import goLogo from '~/assets/go-logo.svg'
import viteLogo from '~/assets/vite-logo.svg'
import './app.css'

function App() {
  const [count, setCount] = useState(0)

  return (
    <>
      <section id="center">
        <div className="hero">
          <img src={goLogo} className="base" width="120" height="120" alt="Go logo" />
          <img src={viteLogo} className="framework" width="120" height="120" alt="Vite logo" />
        </div>
        <div>
          <h1>Get started</h1>
          <p>
            Edit <code>page/app.tsx</code> and save to test <code>HMR</code>
          </p>
        </div>
        <button
          type="button"
          className="counter"
          onClick={() => setCount((count) => count + 1)}
        >
          Count is {count}
        </button>
      </section>

      <div className="ticks"></div>
      <section id="spacer"></section>
    </>
  )
}

export default App
