import React, { useState } from 'react';
import axios from 'axios';

const AuthForm: React.FC<{ setUser: (user: any) => void }> = ({ setUser }) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const response = await axios.post(`${process.env.REACT_APP_HOST_URL}:${process.env.REACT_APP_SERVER_PORT}/login`, { email, password });
      setUser(response.data.Email);
    } 
    catch (error) {
      console.error('ERROR: could not log in:', error);
      alert('ERROR: Could not log in. Please try again or check the console logs.');
    }
  };

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const response = await axios.post(`${process.env.REACT_APP_HOST_URL}:${process.env.REACT_APP_SERVER_PORT}/register`, { email, password });
      setUser(response.data.Email);
    } 
    catch (error) {
      console.error('ERROR: could not register:', error);
      alert('There was an error registering. Please try again or check the console logs.');
    }
  };

  return (
    <form>
      <div>
        <label>Email:</label>
        <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} required />
      </div>
      <div>
        <label>Password:</label>
        <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} required />
      </div>
      <button type="submit" onClick={handleLogin}>Login</button>
      <button type="submit" onClick={handleRegister}>Register</button>
    </form>
  );
};

export default AuthForm;
