import { useEffect, useState } from 'react';
import CreateProjectForm from './CreateProjectForm.jsx';
import CreateTaskForm from '../CreateTaskForm.jsx';
import CreatePromptForm from '../CreatePromptForm.jsx';
import JoinProjectForm from './JoinProjectForm.jsx';

import { SortablePromptItem } from '../SortablePromptItem.jsx';
import PromptReorderList from '../PromptReorderList.jsx';

import ProjectRow from './ProjectRow.jsx';
import TaskBlock from './TaskBlock.jsx';

export default function Projects() {
  const [projects, setProjects] = useState([]);
  const [selectedProjectId, setSelectedProjectId] = useState(null);
  const [expandedTasks, setExpandedTasks] = useState({});
  const [tasks, setTasks] = useState({});
  const [prompts, setPrompts] = useState({});

  const [editedPromptTexts, setEditedPromptTexts] = useState({});
  const [editedPromptNames, setEditedPromptNames] = useState({});
  const [editingPromptIds, setEditingPromptIds] = useState({});

  const [editingProjectId, setEditingProjectId] = useState(null);
  const [editingProjectName, setEditingProjectName] = useState('');
  const [editingProjectApiKey, setEditingProjectApiKey] = useState('');

  const [editingTaskId, setEditingTaskId] = useState(null);
  const [editingTaskName, setEditingTaskName] = useState('');
  const [editingTaskApiMethod, setEditingTaskApiMethod] = useState('');
  const [editingTaskDescription, setEditingTaskDescription] = useState('');

  const [error, setError] = useState('');

  const token = localStorage.getItem('token');

  const arr = (val) => (Array.isArray(val) ? val : []);
  const toArray = (data, fallbackKey = 'value') =>
    Array.isArray(data) ? data : arr(data?.[fallbackKey]);

  const fetchProjects = async () => {
    setError('');
    try {
      const res = await fetch(
        `${import.meta.env.VITE_API_URL}/api/projects/my-projects`,
        {
          headers: { Authorization: `Bearer ${token}` },
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
        setError(backendError);
        setProjects([]);
        return;
      }

      const list = toArray(data?.projects) || toArray(data) || [];
      setProjects(list);
    } catch (err) {
      setError(err?.message || '');
      setProjects([]);
    }
  };

  const fetchTasks = async (projectId) => {
    try {
      const res = await fetch(
        `${import.meta.env.VITE_API_URL}/api/tasks/project/${projectId}`,
        {
          headers: { Authorization: `Bearer ${token}` },
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
        setError(backendError);
        setTasks((prev) => ({ ...prev, [projectId]: [] }));
        return;
      }

      const list = toArray(data) || toArray(data, 'tasks') || [];
      setTasks((prev) => ({ ...prev, [projectId]: list }));
    } catch (err) {
      setError(err?.message || '');
      setTasks((prev) => ({ ...prev, [projectId]: [] }));
    }
  };

  const fetchPrompts = async (taskId) => {
    if (!taskId) return;

    try {
      const res = await fetch(
        `${import.meta.env.VITE_API_URL}/api/prompts/task/${taskId}`,
        {
          headers: { Authorization: `Bearer ${token}` },
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
        setError(backendError);
        setPrompts((prev) => ({ ...prev, [taskId]: [] }));
        return;
      }

      const list = toArray(data) || toArray(data, 'value') || [];
      setPrompts((prev) => ({ ...prev, [taskId]: list }));

      const newTexts = {};
      const newNames = {};
      list.forEach((p) => {
        newTexts[p._id] = p.promptText;
        newNames[p._id] = p.name;
      });
      setEditedPromptTexts((prev) => ({ ...prev, ...newTexts }));
      setEditedPromptNames((prev) => ({ ...prev, ...newNames }));
    } catch (err) {
      setError(err?.message || '');
      setPrompts((prev) => ({ ...prev, [taskId]: [] }));
    }
  };

  useEffect(() => {
    fetchProjects();
  }, []);

  const handleProjectClick = async (projectId) => {
    if (selectedProjectId === projectId) {
      setSelectedProjectId(null);
      return;
    }
    setSelectedProjectId(projectId);
    setError('');

    try {
      const res = await fetch(
        `${import.meta.env.VITE_API_URL}/api/projects/${projectId}`,
        {
          headers: { Authorization: `Bearer ${token}` },
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
        setError(backendError);
        return;
      }

      const project = data?.value ?? data;
      setProjects((prev) =>
        arr(prev).map((p) => (p._id === projectId ? { ...p, ...project } : p))
      );
    } catch (err) {
      setError(err?.message || '');
    } finally {
      fetchTasks(projectId);
    }
  };

  const handleTaskClick = (taskId) => {
    setExpandedTasks((prev) => ({ ...prev, [taskId]: !prev[taskId] }));
    if (!prompts[taskId]) fetchPrompts(taskId);
  };

  const handleTaskCreated = async (projectId, newTaskId) => {
    await fetchTasks(projectId);
    setExpandedTasks((prev) => ({ ...prev, [newTaskId]: true }));
    fetchPrompts(newTaskId);
  };

  const handleDeleteProject = async (projectId) => {
    setError('');
    try {
      const res = await fetch(
        `${import.meta.env.VITE_API_URL}/api/projects/${projectId}`,
        {
          method: 'DELETE',
          headers: { Authorization: `Bearer ${token}` },
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
        setError(backendError);
        return;
      }

      setProjects((prev) => arr(prev).filter((p) => p._id !== projectId));
    } catch (err) {
      setError(err?.message || '');
    }
  };

  const handleUpdateProject = async (projectId) => {
    setError('');
    try {
      const payload = {};
      if (editingProjectName?.trim()) payload.name = editingProjectName.trim();
      const trimmedKey = (editingProjectApiKey || '').trim();
      if (trimmedKey !== '') payload.apiKey = trimmedKey;

      const res = await fetch(
        `${import.meta.env.VITE_API_URL}/api/projects/${projectId}`,
        {
          method: 'PATCH',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(payload),
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
        setError(backendError);
        return;
      }

      await fetchProjects();
      setEditingProjectId(null);
    } catch (err) {
      setError(err?.message || '');
    }
  };

  const handleSetProjectStatus = async (projectId, status) => {
    setError('');
    try {
      const res = await fetch(
        `${import.meta.env.VITE_API_URL}/api/projects/${projectId}`,
        {
          method: 'PATCH',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({ status: String(status).toLowerCase() }),
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
        setError(backendError);
        return;
      }

      await fetchProjects();
    } catch (err) {
      setError(err?.message || '');
    }
  };

  const handleDeleteTask = async (taskId, projectId) => {
    setError('');
    try {
      const res = await fetch(
        `${import.meta.env.VITE_API_URL}/api/tasks/${taskId}`,
        {
          method: 'DELETE',
          headers: { Authorization: `Bearer ${token}` },
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
        setError(backendError);
        return;
      }

      setTasks((prev) => {
        const list = arr(prev[projectId]).filter((t) => t._id !== taskId);
        return { ...prev, [projectId]: list };
      });
    } catch (err) {
      setError(err?.message || '');
    }
  };

  const handleUpdateTask = async (taskId, projectId) => {
    setError('');
    try {
      const res = await fetch(
        `${import.meta.env.VITE_API_URL}/api/tasks/${taskId}`,
        {
          method: 'PATCH',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            name: editingTaskName,
            apiMethod: editingTaskApiMethod,
            description: editingTaskDescription,
          }),
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
        setError(backendError);
        return;
      }

      await fetchTasks(projectId);
      setEditingTaskId(null);
    } catch (err) {
      setError(err?.message || '');
    }
  };

  const handleSetTaskStatus = async (taskId, status, projectId) => {
    setError('');
    try {
      const res = await fetch(
        `${import.meta.env.VITE_API_URL}/api/tasks/${taskId}`,
        {
          method: 'PATCH',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({ status: String(status).toLowerCase() }),
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
        setError(backendError);
        return;
      }

      await fetchTasks(projectId);
    } catch (err) {
      setError(err?.message || '');
    }
  };

  const handleDeletePrompt = async (promptId, taskId) => {
    setError('');
    try {
      const res = await fetch(
        `${import.meta.env.VITE_API_URL}/api/prompts/${promptId}`,
        {
          method: 'DELETE',
          headers: { Authorization: `Bearer ${token}` },
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
        setError(backendError);
        return;
      }

      setPrompts((prev) => {
        const list = arr(prev[taskId]).filter((p) => p._id !== promptId);
        return { ...prev, [taskId]: list };
      });
    } catch (err) {
      setError(err?.message || '');
    }
  };

  const handleSavePrompt = async (promptId) => {
    setError('');
    try {
      const res = await fetch(
        `${import.meta.env.VITE_API_URL}/api/prompts/${promptId}`,
        {
          method: 'PATCH',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            promptText: editedPromptTexts[promptId],
            name: editedPromptNames[promptId],
          }),
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
        setError(backendError);
        return;
      }

      const updatedPrompt = Array.isArray(data?.value)
        ? data.value[0]
        : data?.value ?? data;

      const taskId = Object.keys(prompts).find((tid) =>
        arr(prompts[tid]).some((p) => p._id === promptId)
      );

      if (updatedPrompt) {
        setEditedPromptTexts((prev) => ({
          ...prev,
          [promptId]: updatedPrompt.promptText,
        }));
        setEditedPromptNames((prev) => ({
          ...prev,
          [promptId]: updatedPrompt.name,
        }));
      }

      setEditingPromptIds((prev) => ({ ...prev, [promptId]: false }));

      if (taskId) {
        await fetchPrompts(taskId);
      }
    } catch (e) {
      setError(e?.message || '');
    }
  };

  const handleRunPrompt = async (promptId, taskId) => {
    setError('');
    try {
      const res = await fetch(
        `${import.meta.env.VITE_API_URL}/api/prompts/${promptId}/run`,
        {
          method: 'POST',
          headers: { Authorization: `Bearer ${token}` },
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
        setError(backendError);
        return;
      }

      await fetchPrompts(taskId);
    } catch (e) {
      setError(e?.message || '');
    }
  };

  const handleRollback = async (promptId, version, taskId) => {
    setError('');
    try {
      const res = await fetch(
        `${import.meta.env.VITE_API_URL}/api/prompts/${promptId}/rollback`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({ version }),
        }
      );

      let data = null;
      try {
        data = await res.json();
      } catch {
        data = null;
      }

      if (!res.ok || (data && data.ok === false)) {
        const backendError =
          data && typeof data.error === 'string' ? data.error : '';
        setError(backendError);
        return;
      }

      await fetchPrompts(taskId);
    } catch (e) {
      setError(e?.message || '');
    }
  };

  const projectList = arr(projects);

  return (
    <div>
      <h2>Your projects</h2>

      {error && <p style={{ color: 'red' }}>{error}</p>}

      <CreateProjectForm onProjectCreated={fetchProjects} />
      <JoinProjectForm onJoined={fetchProjects} />
      <hr />

      {projectList.length === 0 ? (
        <p>No projects</p>
      ) : (
        projectList.map((project) => (
          <div key={project._id} style={{ marginBottom: '2rem' }}>
            <ProjectRow
              project={project}
              selectedProjectId={selectedProjectId}
              editingProjectId={editingProjectId}
              editingProjectName={editingProjectName}
              editingProjectApiKey={editingProjectApiKey}
              setEditingProjectId={setEditingProjectId}
              setEditingProjectName={setEditingProjectName}
              setEditingProjectApiKey={setEditingProjectApiKey}
              onToggle={() => handleProjectClick(project._id)}
              onDelete={() => handleDeleteProject(project._id)}
              onUpdate={() => handleUpdateProject(project._id)}
              onStatusChange={(s) => handleSetProjectStatus(project._id, s)}
            />

            {selectedProjectId === project._id && (
              <div style={{ paddingLeft: '1rem', marginTop: '10px' }}>
                <CreateTaskForm
                  token={token}
                  projectId={project._id}
                  onTaskCreated={handleTaskCreated}
                />
                <h4>Tasks:</h4>

                {arr(tasks[project._id]).length > 0 ? (
                  arr(tasks[project._id]).map((task) => (
                    <TaskBlock
                      key={task._id}
                      task={task}
                      projectId={project._id}
                      expanded={!!expandedTasks[task._id]}
                      editingTaskId={editingTaskId}
                      editingTaskName={editingTaskName}
                      editingTaskApiMethod={editingTaskApiMethod}
                      editingTaskDescription={editingTaskDescription}
                      setEditingTaskId={setEditingTaskId}
                      setEditingTaskName={setEditingTaskName}
                      setEditingTaskApiMethod={setEditingTaskApiMethod}
                      setEditingTaskDescription={setEditingTaskDescription}
                      onToggle={() => handleTaskClick(task._id)}
                      onDelete={() => handleDeleteTask(task._id, project._id)}
                      onUpdate={() => handleUpdateTask(task._id, project._id)}
                      onStatusChange={(s) =>
                        handleSetTaskStatus(task._id, s, project._id)
                      }
                    >
                      <div style={{ paddingLeft: '1rem', marginTop: '5px' }}>
                        <CreatePromptForm
                          token={token}
                          taskId={task._id}
                          onPromptCreated={() => fetchPrompts(task._id)}
                        />
                        <h5>Prompts:</h5>

                        {arr(prompts[task._id]).length > 0 ? (
                          <PromptReorderList
                            taskId={task._id}
                            token={token}
                            items={arr(prompts[task._id])}
                            refresh={() => fetchPrompts(task._id)}
                            onLocalReorder={(reordered) => {
                              setPrompts((prev) => ({
                                ...prev,
                                [task._id]: reordered,
                              }));
                            }}
                            SortablePromptItem={SortablePromptItem}
                            renderItem={(prompt) => (
                              <>
                                <div>
                                  <strong>Name:</strong>{' '}
                                  {editingPromptIds[prompt._id] ? (
                                    <>
                                      <input
                                        style={{ width: '60%' }}
                                        value={
                                          editedPromptNames[prompt._id] ??
                                          prompt.name
                                        }
                                        onChange={(e) =>
                                          setEditedPromptNames((prev) => ({
                                            ...prev,
                                            [prompt._id]: e.target.value,
                                          }))
                                        }
                                      />
                                      <button
                                        onClick={() => {
                                          handleSavePrompt(prompt._id);
                                          setEditingPromptIds((prev) => ({
                                            ...prev,
                                            [prompt._id]: false,
                                          }));
                                        }}
                                      >
                                        💾
                                      </button>
                                      <button
                                        onClick={() =>
                                          setEditingPromptIds((prev) => ({
                                            ...prev,
                                            [prompt._id]: false,
                                          }))
                                        }
                                      >
                                        ❌
                                      </button>
                                    </>
                                  ) : (
                                    <>
                                      {prompt.name}{' '}
                                      <button
                                        onClick={() =>
                                          setEditingPromptIds((prev) => ({
                                            ...prev,
                                            [prompt._id]: true,
                                          }))
                                        }
                                      >
                                        ✏️
                                      </button>
                                    </>
                                  )}
                                </div>

                                <div
                                  style={{
                                    marginTop: '4px',
                                    fontSize: '0.9em',
                                    color: '#666',
                                  }}
                                >
                                  <strong>Version:</strong> {prompt.version}
                                </div>

                                <div>
                                  <strong>Prompt:</strong>
                                  <br />
                                  <textarea
                                    style={{ width: '100%', minHeight: '60px' }}
                                    value={
                                      editedPromptTexts[prompt._id] ??
                                      prompt.promptText
                                    }
                                    onChange={(e) =>
                                      setEditedPromptTexts((prev) => ({
                                        ...prev,
                                        [prompt._id]: e.target.value,
                                      }))
                                    }
                                  />
                                </div>

                                {prompt.responseText && (
                                  <div>
                                    <strong>Response:</strong>{' '}
                                    {prompt.responseText}
                                  </div>
                                )}

                                {arr(prompt.history).length > 0 && (
                                  <div style={{ marginTop: '5px' }}>
                                    <strong>History:</strong>
                                    <ul>
                                      {arr(prompt.history).map(
                                        (entry, index) => (
                                          <li key={index}>
                                            <div>
                                              <em>
                                                {new Date(
                                                  entry.createdAt
                                                ).toLocaleString()}
                                              </em>
                                            </div>
                                            <div>
                                              <strong>Version:</strong>{' '}
                                              {entry.version ?? '—'}
                                            </div>
                                            <div>
                                              <strong>Prompt:</strong>{' '}
                                              {entry.prompt}
                                            </div>
                                            <div>
                                              <strong>Response:</strong>{' '}
                                              {entry.response}
                                            </div>
                                            <button
                                              onClick={() =>
                                                handleRollback(
                                                  prompt._id,
                                                  entry.version,
                                                  task._id
                                                )
                                              }
                                              style={{ marginTop: '5px' }}
                                            >
                                              🔄 Roll back to V
                                              {entry.version}
                                            </button>
                                          </li>
                                        )
                                      )}
                                    </ul>
                                  </div>
                                )}

                                <button
                                  onClick={() => handleSavePrompt(prompt._id)}
                                  style={{ marginRight: '10px' }}
                                >
                                  💾 Save V{prompt.version}
                                </button>
                                <button
                                  onClick={() =>
                                    handleRunPrompt(prompt._id, task._id)
                                  }
                                >
                                  🚀 Run
                                </button>
                                <button
                                  onClick={() =>
                                    handleDeletePrompt(prompt._id, task._id)
                                  }
                                  style={{ marginTop: '5px', color: 'red' }}
                                >
                                  🗑️ Delete
                                </button>
                              </>
                            )}
                          />
                        ) : (
                          <p>No prompts</p>
                        )}
                      </div>
                    </TaskBlock>
                  ))
                ) : (
                  <p>No tasks</p>
                )}
              </div>
            )}
          </div>
        ))
      )}
    </div>
  );
}
