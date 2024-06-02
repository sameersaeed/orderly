import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import axios from 'axios';
import * as toastr from 'toastr';

const Login: React.FC<{ setAuth: (user: any) => void }> = ({ setAuth }) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const navigate = useNavigate();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const login = await axios.post(`${process.env.REACT_APP_HOST_URL}:${process.env.REACT_APP_SERVER_PORT}/login`, { email, password });
      const jwt = await axios.get(`${process.env.REACT_APP_HOST_URL}:${process.env.REACT_APP_SERVER_PORT}/jwt`, {
        headers: { Authorization: `Bearer ${login.data.Token}` }
      });

      // getting jwt token information to keep user logged in while session is valid
      const { token, expiresAt } = jwt.data;
      if(token) {
        localStorage.setItem('tokenExpiry', expiresAt.toString());
        toastr.success('Login successful! Welcome, ' + login.data.Email + '!');
        setAuth({ expiresAt });

        // saving user information for validation once they're logged in
        localStorage.setItem('userData', JSON.stringify(login.data));

        navigate('/shop');
      }
    } 
    catch (error) {
      console.error('ERROR: could not log in:', error);
      toastr.error('ERROR: Could not log in. Please ensure that your credentials are correct.');
    }
  };

  return (
    <div className="container mt-5 text-light">
      <form onSubmit={handleLogin}>
        <div className="mb-3 row">
          <label className="col-sm-2 col-form-label">Email:</label>
          <div className="col-sm-10">
            <input type="email" className="form-control" value={email} onChange={(e) => setEmail(e.target.value)} required />
          </div>
        </div>
        <div className="mb-3 row">
          <label className="col-sm-2 col-form-label">Password:</label>
          <div className="col-sm-10">
            <input type="password" className="form-control" value={password} onChange={(e) => setPassword(e.target.value)} required />
          </div>
        </div>
        <button type="submit" className="btn btn-primary">Login</button>
      </form>
      <p className="mt-3">
        Don't have an account? <Link to="/register">Register here</Link>
      </p>
    </div>
  );
};

export default Login;
