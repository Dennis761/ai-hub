import React, { useState } from 'react';

const ApiKeyItem = ({ keyData, onChange }) => {
  const [editMode, setEditMode] = useState(false);
  const [loading, setLoading] = useState(false);
  const token = localStorage.getItem('token');

  const [formData, setFormData] = useState({
    keyName: keyData.keyName,
    keyValue: keyData.keyValue,
    modelName: keyData.modelName,
    provider: keyData.provider,
    usageEnv: keyData.usageEnv,
    balance: keyData.balance ?? '',
  });

  const handleChange = (e) =>
    setFormData((prev) => ({ ...prev, [e.target.name]: e.target.value }));

  const handleSave = async () => {
    setLoading(true);
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/api-keys/${keyData._id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          ...formData,
          balance: formData.balance === '' ? null : parseFloat(formData.balance),
        }),
      });
  
      const result = await res.json();
  
      if (!res.ok) {
        const errorMessage = result.details || result.error || 'Не вдалося оновити ключ';
        throw new Error(errorMessage);
      }
  
      setEditMode(false);
      onChange();
    } catch (err) {
      alert(err.message);
    } finally {
      setLoading(false);
    }
  };  

  const handleDelete = async () => {
    if (!window.confirm('Видалити ключ?')) return;
    setLoading(true);
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/api-keys/${keyData._id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) throw new Error('Не вдалося видалити ключ');
      onChange();
    } catch (err) {
      alert(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleToggleStatus = async () => {
    setLoading(true);
    try {
      const endpoint = `${import.meta.env.VITE_API_URL}/api/api-keys/${keyData._id}/${keyData.status === 'active' ? 'deactivate' : 'activate'}`;
      const res = await fetch(endpoint, {
        method: 'PATCH',
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) throw new Error('Не вдалося змінити статус');
      onChange();
    } catch (err) {
      alert(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ border: '1px solid #ccc', padding: '1rem', marginBottom: '1rem' }}>
      {!editMode ? (
        <>
          <p><strong>Назва:</strong> {keyData.keyName}</p>
          <p><strong>Тип:</strong> {keyData.usageEnv}</p>
          <p><strong>Статус:</strong> {keyData.status}</p>
          <p><strong>Баланс:</strong> {keyData.balance}</p>
          <button onClick={() => setEditMode(true)}>✏️ Редагувати</button>
          <button onClick={handleToggleStatus} style={{ marginLeft: '0.5rem' }}>
            {keyData.status === 'active' ? '🔴 Деактивувати' : '🟢 Активувати'}
          </button>
          <button onClick={handleDelete} style={{ marginLeft: '0.5rem' }}>
            ❌ Видалити
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
            placeholder="Model Name"
            value={formData.modelName}
            onChange={handleChange}
          />
          <input
            name="keyName"
            placeholder="Key Name"
            value={formData.keyName}
            onChange={handleChange}
          />
          <input
            name="keyValue"
            placeholder="Key Value"
            value={formData.keyValue}
            onChange={handleChange}
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
          />
          <div style={{ marginTop: '0.5rem' }}>
            <button onClick={handleSave} disabled={loading}>💾 Зберегти</button>
            <button onClick={() => setEditMode(false)} style={{ marginLeft: '0.5rem' }}>
              ❌ Скасувати
            </button>
          </div>
        </>

      )}
    </div>
  );
};

export default ApiKeyItem;
