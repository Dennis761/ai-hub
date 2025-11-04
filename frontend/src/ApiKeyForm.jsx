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

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    const parsedBalance =
      formData.balance === ''
        ? null
        : Number.isNaN(parseFloat(formData.balance))
        ? null
        : parseFloat(formData.balance);

    const body = {
      provider: formData.provider.trim(),
      modelName: formData.modelName.trim(),
      keyName: formData.keyName.trim(),
      keyValue: formData.keyValue,
      usageEnv: formData.usageEnv,
      status: formData.status,
      balance: parsedBalance,
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

      let data = null;
      try {
        data = await res.json();
      } catch (_) {
        data = null;
      }

      if (!res.ok) {
        const backendError = data && typeof data.error === 'string' ? data.error : '';
        setError(backendError);
        return;
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
      setError('');
      onSuccess?.();
    } catch (err) {
      setError(err?.message || '');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} style={{ marginBottom: '2rem' }} autoComplete="off">
      <h3>Create New API Key</h3>
      {error && <p style={{ color: 'red' }}>{error}</p>}

      <input
        name="provider"
        placeholder="Provider"
        value={formData.provider}
        onChange={handleChange}
        required
        autoComplete="off"
      />
      <input
        name="modelName"
        placeholder="Model name"
        value={formData.modelName}
        onChange={handleChange}
        required
        autoComplete="off"
      />
      <input
        name="keyName"
        placeholder="Key name (alias)"
        value={formData.keyName}
        onChange={handleChange}
        required
        autoComplete="off"
      />
      <input
        name="keyValue"
        placeholder="Key value"
        type="password"
        value={formData.keyValue}
        onChange={handleChange}
        required
        autoComplete="new-password"
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
        placeholder="Balance"
        type="number"
        value={formData.balance}
        onChange={handleChange}
        min="0"
        step="0.01"
        inputMode="decimal"
      />

      <button type="submit" disabled={loading}>
        {loading ? 'Creating...' : 'Create Key'}
      </button>
    </form>
  );
};

export default ApiKeyForm;
