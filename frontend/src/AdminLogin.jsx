import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

export default function AdminLogin() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [step, setStep] = useState('login');
  const [code, setCode] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const safeJson = async (res) => {
    try {
      return await res.json();
    } catch {
      return {};
    }
  };

  const handleLogin = async () => {
    setLoading(true);
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/admin/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });
      const data = await safeJson(res);

      if (!res.ok) {
        const msg = data?.error;
        throw new Error(msg);
      }

      if (res.ok) {
        localStorage.setItem('token', data.token);
        navigate('/ai-hub');
        alert('Login successful');
      }
    } catch (err) {
      alert(err);
    } finally {
      setLoading(false);
    }
  };

  const handleForgotClick = async () => {
    if (!email) {
      alert('Please enter your email before resetting the password');
      return;
    }

    setLoading(true);
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/admin/reset/start`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email }),
      });
      const data = await safeJson(res);
      
      if (!res.ok) {
        const msg = data?.error;
        throw new Error(msg);
      } else {
        setStep('verify');
        alert('A code has been sent to your email');
      }
      
    } catch (err) {
      alert(err);
    } finally {
      setLoading(false);
    }
  };

  const handleVerifyCode = async () => {
    setLoading(true);
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/admin/reset/confirm`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, code }),
      });
      const data = await safeJson(res);

      if (res.ok) {
        setStep('newPassword');
        alert('Code confirmed. Please create a new password');
      } else {
        alert(data.error || 'Invalid or expired code');
      }
    } catch (err) {
      alert(err);
    } finally {
      setLoading(false);
    }
  };

  const handleSetNewPassword = async () => {
    setLoading(true);
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/admin/reset/change`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, newPassword }),
      });
      const data = await safeJson(res);

      if (res.ok) {
        alert('Password has been changed. Please log in');
        setStep('login');
      } else {
        alert(data.error);
      }
    } catch (err) {
      alert('Error connecting to the server');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: '400px', margin: '0 auto', padding: '1rem' }}>
      {step === 'login' && (
        <>
          <h2>Admin login</h2>
          <input
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="form-control mb-2"
          />
          <input
            placeholder="Password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="form-control mb-2"
          />
          <button
            className="btn btn-primary w-100 mb-2"
            onClick={handleLogin}
            disabled={loading}
          >
            Log in
          </button>
          <button
            className="btn btn-link w-100"
            onClick={handleForgotClick}
            disabled={loading}
          >
            Forgot password?
          </button>
        </>
      )}

      {step === 'verify' && (
        <>
          <h2>Code verification</h2>
          <input
            placeholder="Code from email"
            value={code}
            onChange={(e) => setCode(e.target.value)}
            className="form-control mb-2"
          />
          <button
            className="btn btn-success w-100"
            onClick={handleVerifyCode}
            disabled={loading}
          >
            Confirm code
          </button>
        </>
      )}

      {step === 'newPassword' && (
        <>
          <h2>New password</h2>
          <input
            placeholder="New password"
            type="password"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            className="form-control mb-2"
          />
          <button
            className="btn btn-success w-100"
            onClick={handleSetNewPassword}
            disabled={loading}
          >
            Change password
          </button>
        </>
      )}
    </div>
  );
}
