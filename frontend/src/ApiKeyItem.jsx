import React, { useState } from 'react';

const ApiKeyItem = ({ keyData, onChange }) => {
  const [editMode, setEditMode] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const token = localStorage.getItem('token');

  const [formData, setFormData] = useState({
    keyName: keyData.keyName,
    keyValue: '',
    modelName: keyData.modelName,
    provider: keyData.provider,
    usageEnv: keyData.usageEnv,
    balance: keyData.balance ?? '',
    status: keyData.status,
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSave = async () => {
    setLoading(true);
    setError('');

    const payload = {
      keyName: formData.keyName.trim(),
      provider: formData.provider.trim(),
      modelName: formData.modelName.trim(),
      usageEnv: formData.usageEnv,
      status: formData.status,
      balance:
        formData.balance === ''
          ? null
          : Number.isNaN(parseFloat(formData.balance))
          ? null
          : parseFloat(formData.balance),
    };

    if (formData.keyValue && formData.keyValue.trim() !== '') {
      payload.keyValue = formData.keyValue;
    }

    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/api-keys/${keyData.id}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(payload),
      });

      let result = null;
      try {
        result = await res.json();
      } catch (_) {
        result = null;
      }

      if (!res.ok) {
        const backendError = result && typeof result.error === 'string' ? result.error : '';
        setError(backendError);
        return;
      }

      setError('');
      setEditMode(false);
      onChange && onChange();
    } catch (err) {
      setError(err?.message || '');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!window.confirm('Delete this key?')) return;
    setLoading(true);
    setError('');
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/api-keys/${keyData.id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` },
      });

      if (!res.ok) {
        let data = null;
        try {
          data = await res.json();
        } catch (_) {
          data = null;
        }
        const backendError = data && typeof data.error === 'string' ? data.error : '';
        setError(backendError);
        return;
      }

      setError('');
      onChange && onChange();
    } catch (err) {
      setError(err?.message || '');
    } finally {
      setLoading(false);
    }
  };

  const handleToggleStatus = async () => {
    setLoading(true);
    setError('');
    try {
      const newStatus = keyData.status === 'active' ? 'inactive' : 'active';
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/api-keys/${keyData.id}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ status: newStatus }),
      });

      let result = null;
      try {
        result = await res.json();
      } catch (_) {
        result = null;
      }

      if (!res.ok) {
        const backendError = result && typeof result.error === 'string' ? result.error : '';
        setError(backendError);
        return;
      }

      setError('');
      onChange && onChange();
    } catch (err) {
      setError(err?.message || '');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ border: '1px solid #ccc', padding: '1rem', marginBottom: '1rem' }}>
      {error && <p style={{ color: 'red', marginBottom: '0.5rem' }}>{error}</p>}

      {!editMode ? (
        <>
          <p>
            <strong>Name:</strong> {keyData.keyName}
          </p>
          <p>
            <strong>Environment:</strong> {keyData.usageEnv}
          </p>
          <p>
            <strong>Status:</strong> {keyData.status}
          </p>
          <p>
            <strong>Balance:</strong> {keyData.balance}
          </p>
          <button onClick={() => setEditMode(true)} disabled={loading}>
            ✏️ Edit
          </button>
          <button
            onClick={handleToggleStatus}
            style={{ marginLeft: '0.5rem' }}
            disabled={loading}
          >
            {keyData.status === 'active' ? '🔴 Deactivate' : '🟢 Activate'}
          </button>
          <button
            onClick={handleDelete}
            style={{ marginLeft: '0.5rem' }}
            disabled={loading}
          >
            ❌ Delete
          </button>
        </>
      ) : (
        <>
          <input
            name="provider"
            placeholder="Provider"
            value={formData.provider}
            onChange={handleChange}
          />
          <input
            name="modelName"
            placeholder="Model name"
            value={formData.modelName}
            onChange={handleChange}
          />
          <input
            name="keyName"
            placeholder="Key name"
            value={formData.keyName}
            onChange={handleChange}
          />
          <input
            name="keyValue"
            placeholder="Key value (enter to update)"
            type="password"
            value={formData.keyValue}
            onChange={handleChange}
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
          <div style={{ marginTop: '0.5rem' }}>
            <button onClick={handleSave} disabled={loading}>
              💾 Save
            </button>
            <button
              onClick={() => {
                setEditMode(false);
                setError('');
              }}
              style={{ marginLeft: '0.5rem' }}
              disabled={loading}
            >
              ❌ Cancel
            </button>
          </div>
        </>
      )}
    </div>
  );
};

export default ApiKeyItem;
