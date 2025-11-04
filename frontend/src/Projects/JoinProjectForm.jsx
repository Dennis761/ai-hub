import { useState } from 'react';

function JoinProjectForm({ onJoined }) {
  const [form, setForm] = useState({ name: '', apiKey: '' });
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false);

  const handleChange = (e) => {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (loading) return;

    setMessage('');
    setLoading(true);

    try {
      const body = {
        name: form.name.trim(),
        apiKey: form.apiKey,
      };

      const response = await fetch(`${import.meta.env.VITE_API_URL}/api/projects/join`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
        body: JSON.stringify(body),
      });

      let data = null;
      try {
        data = await response.json();
      } catch {
        data = null;
      }

      if (!response.ok) {
        const backendError =
          data && typeof data.error === 'string' ? data.error : '';
        setMessage(backendError);
        return;
      }

      setMessage(`✅ You have joined the project "${data.name}".`);
      setForm({ name: '', apiKey: '' });
      onJoined?.();
    } catch (err) {
      setMessage(err?.message || '');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ marginTop: '20px' }}>
      <h3>Join a Project</h3>
      <form onSubmit={handleSubmit}>
        <input
          name="name"
          value={form.name}
          onChange={handleChange}
          placeholder="Project name"
          required
          disabled={loading}
        />
        <input
          name="apiKey"
          value={form.apiKey}
          onChange={handleChange}
          placeholder="API key"
          type="password"
          autoComplete="new-password"
          required
          disabled={loading}
        />
        <button type="submit" disabled={loading}>
          {loading ? 'Joining...' : 'Join'}
        </button>
      </form>
      {message && <p>{message}</p>}
    </div>
  );
}

export default JoinProjectForm;
