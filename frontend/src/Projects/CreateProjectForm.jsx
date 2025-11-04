import { useState } from 'react';

function CreateProjectForm({ onProjectCreated }) {
  const [form, setForm] = useState({ name: '', status: 'active', apiKey: '' });
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (loading) return;

    setMessage('');
    setLoading(true);

    try {
      const body = {
        name: form.name.trim(),
        status: form.status.trim().toLowerCase(),
      };

      if (form.apiKey && form.apiKey.trim() !== '') {
        body.apiKey = form.apiKey;
      }

      const response = await fetch(`${import.meta.env.VITE_API_URL}/api/projects`, {
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

      setMessage(`✅ Project "${data.name}" has been created with status "${data.status}".`);
      setForm({ name: '', status: 'active', apiKey: '' });
      onProjectCreated?.(data);
    } catch (err) {
      setMessage(err?.message || '');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ marginTop: '20px' }}>
      <h3>Create a new project</h3>
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
          disabled={loading}
        />
        <select name="status" value={form.status} onChange={handleChange} disabled={loading}>
          <option value="active">active</option>
          <option value="inactive">inactive</option>
          <option value="archived">archived</option>
        </select>
        <button type="submit" disabled={loading}>
          {loading ? 'Creating...' : 'Create'}
        </button>
      </form>
      {message && <p>{message}</p>}
    </div>
  );
}

export default CreateProjectForm;
