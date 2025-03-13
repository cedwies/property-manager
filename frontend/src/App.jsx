import React, { useState, useEffect } from 'react';
import './styles/App.css';
import Houses from './pages/Houses';
import Apartments from './pages/Apartments';
import Tenants from './pages/Tenants';
import Payments from './pages/Payments';
import { GetAppInfo } from '../wailsjs/go/main/App';

function App() {
  const [activePage, setActivePage] = useState('welcome');
  const [appInfo, setAppInfo] = useState({});

  useEffect(() => {
    // Fetch application info from the backend
    const fetchAppInfo = async () => {
      try {
        const info = await GetAppInfo();
        setAppInfo(info);
      } catch (error) {
        console.error('Error fetching app info:', error);
      }
    };

    fetchAppInfo();
  }, []);

  const renderContent = () => {
    switch (activePage) {
      case 'houses':
        return <Houses />;
      case 'apartments':
        return <Apartments />;
      case 'tenants':
        return <Tenants />;
      case 'payments':
        return <Payments />;
      case 'welcome':
      default:
        return (
          <div className="welcome-screen">
            <h2>Welcome to Your Property Management Dashboard</h2>
            <p>Use the sidebar to navigate between features.</p>
            <div className="app-info">
              <p><strong>App Name:</strong> {appInfo.name}</p>
              <p><strong>Version:</strong> {appInfo.version}</p>
              <p><strong>Status:</strong> {appInfo.status}</p>
            </div>
          </div>
        );
    }
  };

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
              <li 
                className={`nav-item ${activePage === 'welcome' ? 'active' : ''}`}
                onClick={() => setActivePage('welcome')}
              >
                Dashboard
              </li>
              <li 
                className={`nav-item ${activePage === 'houses' ? 'active' : ''}`}
                onClick={() => setActivePage('houses')}
              >
                Houses
              </li>
              <li 
                className={`nav-item ${activePage === 'apartments' ? 'active' : ''}`}
                onClick={() => setActivePage('apartments')}
              >
                Apartments
              </li>
              <li 
                className={`nav-item ${activePage === 'tenants' ? 'active' : ''}`}
                onClick={() => setActivePage('tenants')}
              >
                Tenants
              </li>
              <li 
                className={`nav-item ${activePage === 'payments' ? 'active' : ''}`}
                onClick={() => setActivePage('payments')}
              >
                Payments
              </li>
            </ul>
          </nav>
        </div>
        
        <div className="content">
          {renderContent()}
        </div>
      </main>
      
      <footer className="app-footer">
        <p>Property Management System &copy; {new Date().getFullYear()}</p>
      </footer>
    </div>
  );
}

export default App;