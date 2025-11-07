import { useState, useEffect } from 'react'
import axios from 'axios'
import './App.css'

function App() {
  const [backendStatus, setBackendStatus] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // Check backend health
    axios.get('/api/health')
      .then(response => {
        setBackendStatus(response.data)
        setLoading(false)
      })
      .catch(error => {
        console.error('Error connecting to backend:', error)
        setBackendStatus({ status: 'DOWN', message: 'Cannot connect to backend' })
        setLoading(false)
      })
  }, [])

  return (
    <div className="App">
      <header className="App-header">
        <h1>Cowatching</h1>
        <p>Full-stack application with Java Spring Boot and React</p>

        <div className="status-card">
          <h2>Backend Status</h2>
          {loading ? (
            <p>Checking backend connection...</p>
          ) : (
            <>
              <p className={`status ${backendStatus?.status === 'UP' ? 'up' : 'down'}`}>
                Status: {backendStatus?.status || 'UNKNOWN'}
              </p>
              <p>{backendStatus?.message}</p>
            </>
          )}
        </div>
      </header>
    </div>
  )
}

export default App
