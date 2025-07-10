import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

export default function AdminLogin() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [step, setStep] = useState('login');
  const [code, setCode] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const navigate = useNavigate();

  const handleLogin = async () => {

    const res = await fetch(`${import.meta.env.VITE_API_URL}/api/admin/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    });
    const data = await res.json();
    if (res.ok) {
      localStorage.setItem('token', data.token);
      navigate('/ai-hub');
      alert('Вхід виконано успішно');
    } else {
      alert(data.message || 'Помилка входу');
    }
  };

  const handleForgotClick = async () => {
    if (!email) {
      alert('Будь ласка, введіть email перед відновленням паролю');
      return;
    }

    const res = await fetch(`${import.meta.env.VITE_API_URL}/api/admin/forgot-password`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email }),
    });
    const data = await res.json();

    if (res.ok) {
      setStep('verify');
      alert('Код надіслано на email');
    } else {
      alert(data.message || 'Помилка надсилання коду');
    }
  };

  const handleVerifyCode = async () => {
    const res = await fetch(`${import.meta.env.VITE_API_URL}/api/admin/verify-reset-code`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, code }),
    });
    const data = await res.json();
    if (res.ok) {
      setStep('newPassword');
      alert('Код підтверджено. Створіть новий пароль');
    } else {
      alert(data.message || 'Невірний або прострочений код');
    }
  };

  const handleSetNewPassword = async () => {
    const res = await fetch(`${import.meta.env.VITE_API_URL}/api/admin/set-new-password`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password: newPassword }),
    });
    const data = await res.json();
    if (res.ok) {
      alert('Пароль змінено. Тепер увійдіть');
      setStep('login');
    } else {
      alert(data.message || 'Помилка зміни паролю');
    }
  };

  return (
    <div style={{ maxWidth: '400px', margin: '0 auto', padding: '1rem' }}>
      {step === 'login' && (
        <>
          <h2>Вхід адміністратора</h2>
          <input
            placeholder='Email'
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className='form-control mb-2'
          />
          <input
            placeholder='Пароль'
            type='password'
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className='form-control mb-2'
          />
          <button className='btn btn-primary w-100 mb-2' onClick={handleLogin}>
            Увійти
          </button>
          <button className='btn btn-link w-100' onClick={handleForgotClick}>
            Забули пароль?
          </button>
        </>
      )}

      {step === 'verify' && (
        <>
          <h2>Підтвердження коду</h2>
          <input
            placeholder='Код з email'
            value={code}
            onChange={(e) => setCode(e.target.value)}
            className='form-control mb-2'
          />
          <button className='btn btn-success w-100' onClick={handleVerifyCode}>
            Підтвердити код
          </button>
        </>
      )}

      {step === 'newPassword' && (
        <>
          <h2>Новий пароль</h2>
          <input
            placeholder='Новий пароль'
            type='password'
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            className='form-control mb-2'
          />
          <button className='btn btn-success w-100' onClick={handleSetNewPassword}>
            Змінити пароль
          </button>
        </>
      )}
    </div>
  );
}
