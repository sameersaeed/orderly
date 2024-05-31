import React, { useState } from 'react';
import AuthForm from './AuthForm';
import Shop from './Shop';
import './App.css';

const App: React.FC = () => {
  const [user, setUser] = useState(null);

  // on login, identify user with their user data (i.e., email, name, etc.)
  const handleLogin = (userData: any) => {
    setUser(userData);
  };

  const handleLogout = () => {
    setUser(null);
  };

  return (
    <div className="App">
      {/* display shop to authorized users, otherwise show login / registration page */}
      {user ? <Shop user={user} onLogout={handleLogout} /> : <AuthForm setUser={handleLogin} />}
    </div>
  );
}

export default App;
