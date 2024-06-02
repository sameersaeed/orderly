import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import axios from 'axios';
import * as toastr from 'toastr';

const Register: React.FC = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const navigate = useNavigate();

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await axios.post(`${process.env.REACT_APP_HOST_URL}:${process.env.REACT_APP_SERVER_PORT}/register`, { email, password });
      toastr.success('Registration successful!');
      navigate('/login');
    } 
    catch (error) {
      console.error('ERROR: could not register:', error);
      toastr.error('There was an error registering. Please try again or check the console logs.');
    }
  };

  return (
    <div className="container mt-5 text-light">
      <form onSubmit={handleRegister}>
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
        <button type="submit" className="btn btn-secondary">Register</button>
      </form>
      <p className="mt-3">
        Already have an account? <Link to="/login">Login here</Link>
      </p>
    </div>
  );
};

export default Register;
