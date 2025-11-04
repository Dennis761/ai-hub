import { useState, useEffect } from 'react';

function CreatePromptForm({ taskId, onPromptCreated }) {
  const [form, setForm] = useState({
    taskId: '',
    name: '',
    modelId: '',
    promptText: '',
  });
  const [models, setModels] = useState([]);
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setForm((prev) => ({ ...prev, taskId: taskId || '' }));
  }, [taskId]);

  useEffect(() => {
    const fetchModels = async () => {
      try {
        const res = await fetch(
          `${import.meta.env.VITE_API_URL}/api/api-keys/my-keys`,
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('token')}`,
            },
          }
        );

        let data = null;
        try {
          data = await res.json();
        } catch {
          data = null;
        }

        if (!res.ok) {
          const backendError =
            data && typeof data.error === 'string' ? data.error : '';
          console.error('Failed to fetch models:', backendError);
          return;
        }

        const list = Array.isArray(data) ? data : data?.value ?? [];
        setModels(list.filter((m) => m.status === 'active'));
      } catch (err) {
        console.error('Failed to fetch models:', err);
      }
    };

    fetchModels();
  }, []);

  const handleChange = (e) => {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (loading) return;

    setMessage('');
    setLoading(true);

    try {
      const payload = {
        taskId: taskId || form.taskId,
        name: form.name.trim(),
        modelId: form.modelId,
        promptText: form.promptText,
      };

      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/prompts`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${localStorage.getItem('token')}`,
        },
        body: JSON.stringify(payload),
      });

      let data = null;
      try {
        data = await res.json();
      } catch {
        data = null;
      }

      if (!res.ok) {
        const backendError =
          data && typeof data.error === 'string' ? data.error : '';
        setMessage(backendError);
        return;
      }

      const created = Array.isArray(data?.value)
        ? data.value[0]
        : data?.value ?? data;
      const promptId = created?._id || created?.id;
      const promptName = created?.name || payload.name;

      setMessage(`✅ Prompt "${promptName}" has been created!`);
      setForm({ taskId: taskId || '', name: '', modelId: '', promptText: '' });

      if (!promptId) {
        console.warn('No promptId in create response:', data);
      } else {
        await runPrompt(promptId);
      }

      onPromptCreated?.();
    } catch (err) {
      setMessage(err?.message || '');
    } finally {
      setLoading(false);
    }
  };

  const runPrompt = async (promptId) => {
    if (!promptId) {
      console.warn('No promptId for execution');
      return;
    }

    try {
      const res = await fetch(
        `${import.meta.env.VITE_API_URL}/api/prompts/${encodeURIComponent(
          promptId
        )}/run`,
        {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`,
            Accept: 'application/json',
          },
        }
      );

      let data = null;
      try {
        data = await res.json();
      } catch {
        data = null;
      }

      if (!res.ok || (data && data.isSuccess === false)) {
        const backendError =
          data && typeof data.error === 'string' ? data.error : '';
        setMessage(backendError);
        return;
      }
    } catch (err) {
      setMessage(err?.message || '');
    }
  };

  return (
    <div style={{ marginTop: '20px' }}>
      <h3>Create prompt</h3>
      <form onSubmit={handleSubmit}>
        <input
          name="name"
          placeholder="Prompt name"
          value={form.name}
          onChange={handleChange}
          required
          disabled={loading}
        />
        <textarea
          name="promptText"
          placeholder="Prompt text with {{placeholders}}"
          value={form.promptText}
          onChange={handleChange}
          required
          disabled={loading}
        />
        <select
          name="modelId"
          value={form.modelId}
          onChange={handleChange}
          required
          disabled={loading}
        >
          <option value="">🔽 Select model</option>
          {models.map((m) => (
            <option key={m._id || m.id} value={m._id || m.id}>
              {m.keyName}:{m.modelName} ({m.provider})
            </option>
          ))}
        </select>
        <button type="submit" disabled={loading}>
          {loading ? 'Creating…' : 'Create'}
        </button>
      </form>
      {message && <p>{message}</p>}
    </div>
  );
}

export default CreatePromptForm;
