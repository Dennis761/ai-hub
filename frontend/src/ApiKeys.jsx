import React, { useEffect, useState } from 'react';
import ApiKeyList from './ApiKeyList.jsx';
import ApiKeyForm from './ApiKeyForm.jsx';

const ApiKeys = () => {
  const [keys, setKeys] = useState([]);
  const [loading, setLoading] = useState(false);

  const fetchKeys = async () => {
    setLoading(true);
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/api-keys`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
      });

      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.details || 'Failed to fetch keys');
      }

      const data = await res.json();
      setKeys(data);
    } catch (err) {
      console.error('Fetch error:', err.message);
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
      <ApiKeyForm onSuccess={fetchKeys} />
      <ApiKeyList keys={keys} loading={loading} onChange={fetchKeys} />
    </div>
  );
};

export default ApiKeys;