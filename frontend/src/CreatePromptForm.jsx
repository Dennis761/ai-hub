import { useState, useEffect } from 'react';

function CreatePromptForm({ taskId, onPromptCreated }) {
  const [form, setForm] = useState({
    taskId: '',
    name: '',
    modelId: '',
    promptText: ''
  });
  const [models, setModels] = useState([]);
  const [message, setMessage] = useState('');

  useEffect(() => {
    setForm(prev => ({ ...prev, taskId: taskId || '' }));
  }, [taskId]);

  useEffect(() => {
    const fetchModels = async () => {
      try {
        const res = await fetch(`${import.meta.env.VITE_API_URL}/api/api-keys`, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`
          }
        });
        const data = await res.json();
        if (Array.isArray(data)) {
          setModels(data.filter(m => m.status === 'active'));
        }
      } catch (err) {
        console.error('Failed to fetch models:', err);
      }
    };

    fetchModels();
  }, []);

  const handleChange = (e) => {
    setForm(prev => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setMessage('');

    try {
      const response = await fetch(`${import.meta.env.VITE_API_URL}/api/prompts`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify(form)
      });

      const resultPrompt = await response.json();

      if (!response.ok) throw new Error(resultPrompt.details || resultPrompt.error);

      setMessage(`✅ Промпт "${resultPrompt.name}" створено!`);
      setForm({
        taskId: taskId || '',
        name: '',
        modelId: '',
        promptText: ''
      });

      await runPrompt(resultPrompt._id);

      if (onPromptCreated) onPromptCreated();

    } catch (err) {
      setMessage(`❌ Помилка: ${err.message}`);
    }
  };

  const runPrompt = async (promptId) => {
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/prompts/${promptId}/run`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
      });

      const data = await res.json();
      if (!res.ok) throw new Error(data.details || data.error);

    } catch (err) {
      console.error('❌ Помилка виконання:', err.message);
    }
  };

  return (
    <div style={{ marginTop: '20px' }}>
      <h3>Створити промпт</h3>
      <form onSubmit={handleSubmit}>
        <input
          name="name"
          placeholder="Назва промпта"
          value={form.name}
          onChange={handleChange}
          required
        />
        <textarea
          name="promptText"
          placeholder="Текст промпта з {{placeholders}}"
          value={form.promptText}
          onChange={handleChange}
          required
        />
        <select name="modelId" value={form.modelId} onChange={handleChange} required>
          <option value="">🔽 Оберіть модель</option>
          {models.map((model) => (
            <option key={model._id} value={model._id}>
              {model.modelName} ({model.provider})
            </option>
          ))}
        </select>
        <button type="submit">Створити</button>
      </form>

      {message && <p>{message}</p>}
    </div>
  );
}

export default CreatePromptForm;
