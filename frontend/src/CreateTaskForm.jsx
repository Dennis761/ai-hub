import { useState } from 'react';

function CreateTaskForm({ projectId, onTaskCreated }) {
  const [form, setForm] = useState({
    name: '',
    description: '',
    apiMethod: '',
    version: '1.0.0',
    status: 'active'
  });
  const [message, setMessage] = useState('');

  const handleChange = (e) => {
    setForm(prev => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setMessage('');

    try {
      const response = await fetch(`${import.meta.env.VITE_API_URL}/api/tasks`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({
          ...form,
          projectId
        })
      });

      const result = await response.json();

      if (!response.ok) throw new Error(result.details || result.error);

      setMessage(`✅ Задачу "${result.name}" створено`);
      setForm({
        name: '',
        description: '',
        apiMethod: '',
        status: 'active'
      });

      if (onTaskCreated) onTaskCreated(projectId, result._id);

    } catch (err) {
      setMessage(`❌ Помилка: ${err.message}`);
    }
  };

  return (
    <div style={{ marginTop: '10px' }}>
      <form onSubmit={handleSubmit}>
        <input
          name="name"
          placeholder="Task Name"
          value={form.name}
          onChange={handleChange}
          required
        />
        <input
          name="apiMethod"
          placeholder="API Method"
          value={form.apiMethod}
          onChange={handleChange}
          required
        />
        <input
          name="description"
          placeholder="Description"
          value={form.description}
          onChange={handleChange}
        />
        <button type="submit">Додати задачу</button>
      </form>
      {message && <p>{message}</p>}
    </div>
  );
}

export default CreateTaskForm;
