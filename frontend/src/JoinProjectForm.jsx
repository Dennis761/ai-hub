import { useState } from 'react';

function JoinProjectForm({ onJoined }) {
  const [form, setForm] = useState({ name: '', apiKey: '' });
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false);

  const handleChange = (e) => {
    setForm(prev => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setMessage('');
    if (loading) return;
    setLoading(true);
  
    try {
      const response = await fetch(`${import.meta.env.VITE_API_URL}/api/projects/join`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
        body: JSON.stringify(form),
      });
  
      const result = await response.json();
  
      if (!response.ok) throw new Error(result.details || result.error || 'Не вдалося приєднатись до проєкту');
  
      setMessage(result.message);
      setForm({ name: '', apiKey: '' });
  
    } catch (err) {
      setMessage(`❌ Помилка: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };  

  return (
    <div style={{ marginTop: '20px' }}>
      <h3>Приєднатися до проєкту</h3>
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
          required
        />
        <button type="submit" disabled={loading}>
          {loading ? 'Приєднання...' : 'Приєднатися'}
        </button>
      </form>
      {message && <p>{message}</p>}
    </div>
  );
}

export default JoinProjectForm;
