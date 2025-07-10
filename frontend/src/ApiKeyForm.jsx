import React, { useState } from 'react';

const ApiKeyForm = ({ onSuccess }) => {
  const [formData, setFormData] = useState({
    provider: '',
    modelName: '',
    keyName: '',
    keyValue: '',
    status: 'inactive',
    balance: '',
    usageEnv: 'prod',
  });

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleChange = e => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async e => {
    e.preventDefault();
    setLoading(true);
    setError('');

    const body = {
      ...formData,
      balance: formData.balance === '' ? null : parseFloat(formData.balance),
    };

    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/api-keys`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
        body: JSON.stringify(body),
      });

      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.details || 'Failed to create key');
      }

      setFormData({
        provider: '',
        modelName: '',
        keyName: '',
        keyValue: '',
        status: 'inactive',
        balance: '',
        usageEnv: 'prod',
      });

      if (onSuccess) onSuccess();
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} style={{ marginBottom: '2rem' }}>
      <h3>Create New API Key</h3>
      {error && <p style={{ color: 'red' }}>{error}</p>}

      <input
        name="provider"
        placeholder="Provider"
        value={formData.provider}
        onChange={handleChange}
        required
      />
      <input
        name="modelName"
        placeholder="Model Name"
        value={formData.modelName}
        onChange={handleChange}
        required
      />
      <input
        name="keyName"
        placeholder="Key Name"
        value={formData.keyName}
        onChange={handleChange}
        required
      />
      <input
        name="keyValue"
        placeholder="Key Value"
        value={formData.keyValue}
        onChange={handleChange}
        required
      />

      <select name="usageEnv" value={formData.usageEnv} onChange={handleChange}>
        <option value="prod">Production</option>
        <option value="dev">Development</option>
        <option value="test">Testing</option>
      </select>

      <select name="status" value={formData.status} onChange={handleChange}>
        <option value="active">Active</option>
        <option value="inactive">Inactive</option>
      </select>

      <input
        name="balance"
        placeholder="Balance (optional)"
        type="number"
        value={formData.balance}
        onChange={handleChange}
        min="0"
      />

      <button type="submit" disabled={loading}>
        {loading ? 'Creating...' : 'Create Key'}
      </button>
    </form>
  );
};

export default ApiKeyForm;
