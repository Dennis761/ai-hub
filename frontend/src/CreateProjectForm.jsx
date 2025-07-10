import { useState } from 'react';

function CreateProjectForm({ onProjectCreated }) {
  const [form, setForm] = useState({ name: '', status: 'active', apiKey: '' });
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false);

  const handleChange = (e) => {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setMessage('');
    if (loading) return;

    setLoading(true);
    try {
      const response = await fetch(`${import.meta.env.VITE_API_URL}/api/projects`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify(form),
      });

      const result = await response.json();

      if (!response.ok) throw new Error(result.details || result.error || 'Помилка створення проєкту');

      setMessage(`✅ Проєкт "${result.name}" створено!`);
      setForm({ name: '', status: 'active', apiKey: '' });

      if (onProjectCreated) onProjectCreated(result);
    } catch (err) {
      setMessage(`❌ Помилка: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ marginTop: '20px' }}>
      <h3>Створити новий проєкт</h3>
      <form onSubmit={handleSubmit}>
        <input
          name="name"
          value={form.name}
          onChange={handleChange}
          placeholder="Назва проєкту"
          required
        />
        <input
          name="apiKey"
          value={form.apiKey}
          onChange={handleChange}
          placeholder="API ключ"
        />
        <select name="status" value={form.status} onChange={handleChange}>
          <option value="active">active</option>
          <option value="inactive">inactive</option>
          <option value="archived">archived</option>
        </select>
        <button type="submit" disabled={loading}>
          {loading ? 'Створення...' : 'Створити'}
        </button>
      </form>
      {message && <p>{message}</p>}
    </div>
  );
}

export default CreateProjectForm;
