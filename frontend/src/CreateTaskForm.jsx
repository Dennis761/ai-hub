import { useState } from 'react';

function CreateTaskForm({ projectId, onTaskCreated }) {
  const [form, setForm] = useState({
    name: '',
    description: '',
    apiMethod: '',
    status: 'active',
  });
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [loading, setLoading] = useState(false);

  const handleChange = (e) => {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (loading) return;

    setError('');
    setSuccess('');
    setLoading(true);

    try {
      const payload = {
        name: form.name.trim(),
        description: form.description.trim() || undefined,
        apiMethod: form.apiMethod.trim().toLowerCase(),
        status: form.status.trim().toLowerCase(),
        projectId,
      };

      const response = await fetch(`${import.meta.env.VITE_API_URL}/api/tasks`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
        body: JSON.stringify(payload),
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
        setError(backendError);
        return;
      }

      setSuccess(`✅ Task "${data.name || payload.name}" has been created`);
      setForm({
        name: '',
        description: '',
        apiMethod: '',
        status: 'active',
      });

      const createdId = data._id || data.id;
      onTaskCreated?.(projectId, createdId);
    } catch (err) {
      setError(err?.message || '');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ marginTop: '10px' }}>
      <form onSubmit={handleSubmit}>
        <input
          name="name"
          placeholder="Task name"
          value={form.name}
          onChange={handleChange}
          required
          disabled={loading}
        />

        <input
          name="apiMethod"
          placeholder="API method (e.g. chat, embeddings)"
          value={form.apiMethod}
          onChange={handleChange}
          required
          disabled={loading}
        />

        <input
          name="description"
          placeholder="Description"
          value={form.description}
          onChange={handleChange}
          disabled={loading}
        />

        <select
          name="status"
          value={form.status}
          onChange={handleChange}
          disabled={loading}
        >
          <option value="active">active</option>
          <option value="inactive">inactive</option>
          <option value="archived">archived</option>
        </select>

        <button type="submit" disabled={loading}>
          {loading ? 'Creating…' : 'Add task'}
        </button>
      </form>

      {error && <p style={{ color: 'red' }}>{error}</p>}
      {success && <p>{success}</p>}
    </div>
  );
}

export default CreateTaskForm;
