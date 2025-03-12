import React from 'react';
import './styles/App.css';

function App() {
  return (
    <div className="container">
      <header className="app-header">
        <h1>Property Management System</h1>
        <p>Manage your multi-tenant properties efficiently</p>
      </header>
      
      <main className="app-main">
        <div className="sidebar">
          <nav>
            <ul>
              <li className="nav-item">Houses</li>
            </ul>
          </nav>
        </div>
        
        <div className="content">
          <div className="welcome-screen">
            <h2>Welcome to Your Property Management Dashboard</h2>
            <p>This is the initial setup. Features will be implemented step by step.</p>
          </div>
        </div>
      </main>
      
      <footer className="app-footer">
        <p>Property Management System &copy; {new Date().getFullYear()}</p>
      </footer>
    </div>
  );
}

export default App;