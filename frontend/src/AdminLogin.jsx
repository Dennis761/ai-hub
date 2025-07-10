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

      if (res.ok) {
        localStorage.setItem('token', data.token);
        navigate('/ai-hub');
        alert('Вхід виконано успішно');
      } else {
        alert(data.message || 'Помилка входу');
      }
    } catch (err) {
      alert('Помилка зʼєднання з сервером');
    } finally {
      setLoading(false);
    }
  };

  const handleForgotClick = async () => {
    if (!email) {
      alert('Будь ласка, введіть email перед відновленням паролю');
      return;
    }

    setLoading(true);
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/admin/forgot-password`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email }),
      });
      const data = await safeJson(res);

      if (res.ok) {
        setStep('verify');
        alert('Код надіслано на email');
      } else {
        alert(data.message || 'Помилка надсилання коду');
      }
    } catch (err) {
      alert('Помилка зʼєднання з сервером');
    } finally {
      setLoading(false);
    }
  };

  const handleVerifyCode = async () => {
    setLoading(true);
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/admin/verify-reset-code`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, code }),
      });
      const data = await safeJson(res);

      if (res.ok) {
        setStep('newPassword');
        alert('Код підтверджено. Створіть новий пароль');
      } else {
        alert(data.message || 'Невірний або прострочений код');
      }
    } catch (err) {
      alert('Помилка зʼєднання з сервером');
    } finally {
      setLoading(false);
    }
  };

  const handleSetNewPassword = async () => {
    setLoading(true);
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/admin/set-new-password`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password: newPassword }),
      });
      const data = await safeJson(res);

      if (res.ok) {
        alert('Пароль змінено. Тепер увійдіть');
        setStep('login');
      } else {
        alert(data.message || 'Помилка зміни паролю');
      }
    } catch (err) {
      alert('Помилка зʼєднання з сервером');
    } finally {
      setLoading(false);
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
          <button
            className='btn btn-primary w-100 mb-2'
            onClick={handleLogin}
            disabled={loading}
          >
            Увійти
          </button>
          <button
            className='btn btn-link w-100'
            onClick={handleForgotClick}
            disabled={loading}
          >
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
          <button
            className='btn btn-success w-100'
            onClick={handleVerifyCode}
            disabled={loading}
          >
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
          <button
            className='btn btn-success w-100'
            onClick={handleSetNewPassword}
            disabled={loading}
          >
            Змінити пароль
          </button>
        </>
      )}
    </div>
  );
}
