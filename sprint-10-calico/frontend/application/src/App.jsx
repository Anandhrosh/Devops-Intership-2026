import { useEffect, useState } from 'react'
import './App.css'
import './index.css'

function App() {
  const [images, setImages] = useState([])
  const [selectedFile, setSelectedFile] = useState(null)
  const [loading, setLoading] = useState(false)
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')
  const [backendStatus, setBackendStatus] = useState('Checking...')
  const [databaseStatus, setDatabaseStatus] = useState('Checking...')

  const loadImages = async () => {
    try {
      const response = await fetch('/api/images')
      if (!response.ok) throw new Error('Could not load images')
      const data = await response.json()
      setImages(Array.isArray(data) ? data : [])
    } catch (err) {
      setError(err.message)
    }
  }

  const checkStatus = async () => {
    try {
      const response = await fetch('/api/status')
      if (!response.ok) throw new Error('Backend unavailable')
      const data = await response.json()
      setBackendStatus(data.backend || 'Unknown')
      setDatabaseStatus(data.database || 'Unknown')
    } catch {
      setBackendStatus('Unavailable')
      setDatabaseStatus('Unavailable')
    }
  }

  useEffect(() => {
    checkStatus()
    loadImages()
  }, [])

  const handleUpload = async (event) => {
    event.preventDefault()

    if (!selectedFile) {
      setError('Please select an image first.')
      return
    }

    if (selectedFile.size > 5 * 1024 * 1024) {
      setError('Image must be 5 MB or smaller.')
      return
    }

    setLoading(true)
    setMessage('')
    setError('')

    try {
      const formData = new FormData()
      formData.append('image', selectedFile)

      const response = await fetch('/api/images', {
        method: 'POST',
        body: formData,
      })
      const result = await response.json()

      if (!response.ok) {
        throw new Error(result.error || 'Image upload failed')
      }

      setMessage('Image uploaded successfully!')
      setSelectedFile(null)
      event.target.reset()
      await loadImages()
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (id) => {
    if (!window.confirm('Are you sure you want to delete this image?')) return

    setError('')
    setMessage('')

    try {
      const response = await fetch(`/api/images/${id}`, {
        method: 'DELETE',
      })
      const result = await response.json()

      if (!response.ok) {
        throw new Error(result.error || 'Image deletion failed')
      }

      setMessage('Image deleted successfully!')
      await loadImages()
    } catch (err) {
      setError(err.message)
    }
  }

  return (
    <div className="app">
      <header className="header">
        <div>
          <h1>DevOps Capstone</h1>
          <p>Kubernetes 3-Tier Microservice Application</p>
        </div>
        <span className="environment">CAPSTONE ENVIRONMENT</span>
      </header>

      <main className="container">

        {/* SECTION 1: IMAGE UPLOAD */}
        <section className="image-section">
          <h2>📤 Image Upload</h2>
          <p>Choose an image and upload it to the PostgreSQL database.</p>

          <form className="upload-form" onSubmit={handleUpload}>
            <input
              type="file"
              accept="image/jpeg,image/png,image/gif,image/webp"
              onChange={(event) =>
                setSelectedFile(event.target.files[0] || null)
              }
            />
            <button type="submit" disabled={loading || !selectedFile}>
              {loading ? 'Uploading...' : 'Upload Image'}
            </button>
          </form>

          <p className="file-hint">
            Supported formats: JPEG, PNG, GIF, WebP. Maximum size: 5 MB.
          </p>

          {selectedFile && (
            <p className="selected-file">
              Selected: {selectedFile.name}
            </p>
          )}
        </section>

        {message && <p className="success-message">{message}</p>}
        {error && <p className="error-message">{error}</p>}

        {/* SECTION 2: IMAGE GALLERY VIEWER */}
        <section className="image-section">
          <div className="gallery-heading">
            <div>
              <h2>🖼️ Image Gallery Viewer</h2>
              <p>Browse all images saved in your database.</p>
            </div>
            <button
              type="button"
              className="refresh-button"
              onClick={loadImages}
            >
              Refresh Gallery
            </button>
          </div>

          {images.length === 0 ? (
            <p className="empty-gallery">
              No images available. Upload your first image above.
            </p>
          ) : (
            <div className="image-grid">
              {images.map((image) => (
                <div className="image-card" key={image.id}>
                  <img
                    src={`/api/images/${image.id}`}
                    alt={image.filename}
                  />
                  <div className="image-info">
                    <p title={image.filename}>{image.filename}</p>
                    <span>Image ID: {image.id}</span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </section>

        {/* SECTION 3: DELETE IMAGES */}
        <section className="image-section delete-section">
          <h2>🗑️ Delete Images</h2>
          <p>
            Select an image from this list when you want to remove it
            permanently from PostgreSQL.
          </p>

          {images.length === 0 ? (
            <p className="empty-gallery">There are no images to delete.</p>
          ) : (
            <div className="delete-list">
              {images.map((image) => (
                <div className="delete-item" key={image.id}>
                  <div>
                    <strong>{image.filename}</strong>
                    <span>Image ID: {image.id}</span>
                  </div>
                  <button
                    type="button"
                    className="delete-button"
                    onClick={() => handleDelete(image.id)}
                  >
                    Delete Image
                  </button>
                </div>
              ))}
            </div>
          )}
        </section>

        {/* APPLICATION STATUS */}
        <section className="welcome">
          <h2>Application Overview</h2>
          <p>
            Monitor the frontend, backend API, and PostgreSQL database
            from one dashboard.
          </p>
        </section>

        <section className="services">
          <div className="service-card">
            <div className="service-icon">🌐</div>
            <h3>Frontend</h3>
            <p>React + Vite</p>
            <span className="status">Running</span>
          </div>

          <div className="service-card">
            <div className="service-icon">⚙️</div>
            <h3>Backend API</h3>
            <p>Go REST API</p>
            <span className={backendStatus === 'Running' ? 'status' : 'status pending'}>
              {backendStatus}
            </span>
          </div>

          <div className="service-card">
            <div className="service-icon">🗄️</div>
            <h3>Database</h3>
            <p>PostgreSQL StatefulSet</p>
            <span className={databaseStatus === 'Connected' ? 'status' : 'status pending'}>
              {databaseStatus}
            </span>
          </div>
        </section>

        <section className="architecture">
          <h2>Application Architecture</h2>
          <div className="architecture-flow">
            <div className="architecture-box">🌐 Browser</div>
            <div className="arrow">↓</div>
            <div className="architecture-box">React Frontend</div>
            <div className="arrow">↓</div>
            <div className="architecture-box">Go Backend API</div>
            <div className="arrow">↓</div>
            <div className="architecture-box">PostgreSQL Database</div>
          </div>
        </section>

        <section className="details">
          <h2>Deployment Components</h2>
          <ul>
            <li>NGINX Ingress – External application access</li>
            <li>Kubernetes Deployment – Application workloads</li>
            <li>Horizontal Pod Autoscaler – Backend scaling</li>
            <li>Persistent Volume Claim – Database storage</li>
            <li>NetworkPolicy – Database access protection</li>
            <li>ArgoCD – GitOps deployment synchronization</li>
          </ul>
        </section>
      </main>

      <footer>
        <p>DevOps Internship 2026 | Sprint 10 Capstone</p>
      </footer>
    </div>
  )
}

export default App