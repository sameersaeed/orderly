import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import Home from './Home';
import Login from './Login';
import Navbar from './Navbar';
import Orders from './Orders';
import Register from './Register';
import Shop from './Shop';
import './../styles/App.css';

const App: React.FC = () => {
  const isAuthenticated = localStorage.getItem('auth') === 'true';
  const [auth, setAuth] = useState<boolean>(!!isAuthenticated);
  const user = JSON.parse(localStorage.getItem('userData') || '{}');

  const handleLogin = (authData: any) => {
    setAuth(true);
    localStorage.setItem('auth', 'true');
    localStorage.setItem('tokenExpiry', authData.expiresAt);
  };

  const handleLogout = () => {
    setAuth(false);
    localStorage.removeItem('auth');
    localStorage.removeItem('tokenExpiry');
  };

  const checkTokenExpiry = () => {
    const expiry = localStorage.getItem('tokenExpiry');
    if (expiry) {
      const tokenExpiry = new Date().getTime() > new Date(expiry).getTime();
      if (tokenExpiry) {
        handleLogout();
        toastr.warning('Your session has expired. Please log in again.');
        return false;
      }
      return true;
    }
    return false;
  };

  useEffect(() => {
    setAuth(checkTokenExpiry());
  }, []);

  return (
    <Router>
      <div className="App text-light bg-dark">
        <Navbar auth={auth} onLogout={handleLogout} />
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/login" element={<Login setAuth={handleLogin} />} />
          <Route path="/register" element={<Register />} />
          <Route
            path="/shop"
            element={auth ? <Shop user={user} /> : <Navigate to="/login" />}
          />
          <Route
            path="/orders"
            element={auth ? <Orders user={user} /> : <Navigate to="/login" />}
          />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
