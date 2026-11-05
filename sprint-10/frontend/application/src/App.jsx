
import './App.css'

function App() {
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
            <span className="status">Running locally</span>
          </div>

          <div className="service-card">
            <div className="service-icon">⚙️</div>
            <h3>Backend API</h3>
            <p>Go REST API</p>
            <span className="status pending">Kubernetes deployment pending</span>
          </div>

          <div className="service-card">
            <div className="service-icon">🗄️</div>
            <h3>Database</h3>
            <p>PostgreSQL StatefulSet</p>
            <span className="status pending">Kubernetes deployment pending</span>
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
