import React, { useEffect, useState } from 'react';
import ApiKeyList from './ApiKeyList.jsx';
import ApiKeyForm from './ApiKeyForm.jsx';

const ApiKeys = () => {
  const [keys, setKeys] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const fetchKeys = async () => {
    setLoading(true);
    setError('');

    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/api-keys/my-keys`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
      });

      let data = null;
      try {
        data = await res.json();
      } catch (_) {
        data = null;
      }

      if (!res.ok) {
        const backendError = data && typeof data.error === 'string' ? data.error : '';
        setError(backendError);
        setKeys([]);
        return;
      }

      const list = Array.isArray(data) ? data : data?.value ?? [];
      setKeys(list);
      setError('');
    } catch (err) {
      setError(err?.message || '');
      setKeys([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchKeys();
  }, []);

  return (
    <div>
      <h2>API Key Management</h2>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      <ApiKeyForm onSuccess={fetchKeys} />
      <ApiKeyList keys={keys} loading={loading} onChange={fetchKeys} />
    </div>
  );
};

export default ApiKeys;
