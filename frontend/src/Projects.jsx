import { useEffect, useState } from 'react';
import CreateProjectForm from './CreateProjectForm.jsx';
import CreateTaskForm from './CreateTaskForm.jsx';
import CreatePromptForm from './CreatePromptForm.jsx';
import JoinProjectForm from './JoinProjectForm.jsx';
import {
  DndContext,
  closestCenter,
  PointerSensor,
  useSensor,
  useSensors,
} from '@dnd-kit/core';
import {
  arrayMove,
  SortableContext,
  useSortable,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import {SortablePromptItem} from './SortablePromptItem.jsx'

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

  const sensors = useSensors(useSensor(PointerSensor));

  const token = localStorage.getItem('token');

  const fetchProjects = () => {
    fetch(`${import.meta.env.VITE_API_URL}/api/projects/my`, {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((res) => res.json())
      .then((data) => {
        setProjects(data.projects);
      });
  };

  const fetchTasks = async (projectId) => {
    const res = await fetch(`${import.meta.env.VITE_API_URL}/api/tasks/project/${projectId}`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    const data = await res.json();
    setTasks((prev) => ({ ...prev, [projectId]: data }));
  };

  const fetchPrompts = async (taskId) => {
    const res = await fetch(`${import.meta.env.VITE_API_URL}/api/prompts/task/${taskId}`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    const data = await res.json();
    setPrompts((prev) => ({ ...prev, [taskId]: data }));

    const newTexts = {};
    const newNames = {};
    data.forEach((p) => {
      newTexts[p._id] = p.promptText;
      newNames[p._id] = p.name;
    });
    setEditedPromptTexts((prev) => ({ ...prev, ...newTexts }));
    setEditedPromptNames((prev) => ({ ...prev, ...newNames }));
  };

  useEffect(() => {
    fetchProjects();
  }, []);

  const handleProjectClick = async (projectId) => {
    if (selectedProjectId === projectId) {
      setSelectedProjectId(null);
    } else {
      setSelectedProjectId(projectId);
  
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/projects/${projectId}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
  
      const project = await res.json();
  
      setProjects((prev) =>
        prev.map((p) => (p._id === projectId ? { ...p, ...project } : p))
      );
  
      fetchTasks(projectId);
    }
  };  

  const handleTaskClick = (taskId) => {
    setExpandedTasks(prev => ({
      ...prev,
      [taskId]: !prev[taskId]
    }));
    if (!prompts[taskId]) fetchPrompts(taskId);
  };

  const handleTaskCreated = async (projectId, newTaskId) => {
    await fetchTasks(projectId);
    setExpandedTasks((prev) => ({ ...prev, [newTaskId]: true }));
    fetchPrompts(newTaskId);
  };

  const handleDeleteProject = async (projectId) => {
    await fetch(`${import.meta.env.VITE_API_URL}/api/projects/${projectId}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` },
    });
    setProjects(prev => prev.filter(p => p._id !== projectId));
  };

  const handleDeleteTask = async (taskId, projectId) => {
    await fetch(`${import.meta.env.VITE_API_URL}/api/tasks/${taskId}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` },
    });
    setTasks(prev => ({
      ...prev,
      [projectId]: prev[projectId].filter(t => t._id !== taskId)
    }));
  };

  const handleDeletePrompt = async (promptId, taskId) => {
    await fetch(`${import.meta.env.VITE_API_URL}/api/prompts/${promptId}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` },
    });
    setPrompts(prev => ({
      ...prev,
      [taskId]: prev[taskId].filter(p => p._id !== promptId)
    }));
  };

  const handleSavePrompt = async (promptId) => {
    try {
      const updatedText = editedPromptTexts[promptId];
      const updatedName = editedPromptNames[promptId];
  
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/prompts/${promptId}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ promptText: updatedText, name: updatedName }),
      });
  
      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || 'Save failed');
      }
  
      const updatedPrompt = await res.json();
      alert('Промпт успішно збережено!');
  
      const taskId = Object.keys(prompts).find(taskId =>
        prompts[taskId].some(p => p._id === promptId)
      );
  
      if (taskId) {
        setPrompts(prev => ({
          ...prev,
          [taskId]: prev[taskId].map(p =>
            p._id === promptId ? updatedPrompt : p
          )
        }));
      }
    } catch (error) {
      console.error('Помилка при збереженні промпта:', error);
      alert('Помилка при збереженні промпта');
    }
  };
  

  const handleRunPrompt = async (promptId, taskId) => {
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/prompts/${promptId}/run`, {
        method: 'POST',
        headers: { Authorization: `Bearer ${token}` },
      });

      const data = await res.json();

        if (res.ok) {
          alert('Промпт успішно виконано!');
          fetchPrompts(taskId);
        } else {
          const errorMessage =
            typeof data.error === 'string'
              ? data.error
              : JSON.stringify(data.error, null, 2);
          alert(`Помилка:\n${errorMessage}`);
        }        
    } catch (error) {
      console.error('Помилка при запуску промпта:', error);
      alert('Помилка при запуску промпта');
    }
  };

  const handleUpdateProject = async (projectId) => {
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/projects/${projectId}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          name: editingProjectName,
          apiKey: editingProjectApiKey.trim(),
        }),
      });
  
      const data = await res.json();
  
      if (!res.ok) {
        throw new Error(data?.details || data?.error || 'Update failed');
      }
  
      fetchProjects();
      setEditingProjectId(null);
    } catch (err) {
      alert(`Помилка оновлення проєкту: ${err.message}`);
    }
  };  

  const handleSetProjectStatus = async (projectId, status) => {
    await fetch(`${import.meta.env.VITE_API_URL}/api/projects/${projectId}/status`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ status }),
    });
    fetchProjects();
  };

  const handleUpdateTask = async (taskId, projectId) => {
    await fetch(`${import.meta.env.VITE_API_URL}/api/tasks/${taskId}`, {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        name: editingTaskName,
        apiMethod: editingTaskApiMethod,
        description: editingTaskDescription
      }),
    });
    fetchTasks(projectId);
    setEditingTaskId(null);
  };

  const handleSetTaskStatus = async (taskId, status, projectId) => {
    await fetch(`${import.meta.env.VITE_API_URL}/api/tasks/${taskId}/status`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ status }),
    });
    fetchTasks(projectId);
  };

  const handleRollback = async (promptId, versionIndex, taskId) => {
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/prompts/${promptId}/rollback`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ versionIndex }),
      });
  
      const data = await res.json();
  
      if (!res.ok) {
        throw new Error(data.error || 'Rollback failed');
      }
   
      fetchPrompts(taskId); 
    } catch (err) {
      alert(`Помилка rollback: ${err.message}`);
    }
  };
  
  return (
    <div>
      <h2>Ваші проєкти</h2>
        <CreateProjectForm onProjectCreated={fetchProjects} />
        <JoinProjectForm onJoined={fetchProjects} />

      <hr />

      {projects.map((project) => (
        <div key={project._id} style={{ marginBottom: '2rem' }}>
          {editingProjectId === project._id ? (
            <>
            <input
              value={editingProjectName}
              onChange={(e) => setEditingProjectName(e.target.value)}
              placeholder="Назва"
            />
            <input
              value={editingProjectApiKey}
              onChange={(e) => setEditingProjectApiKey(e.target.value)}
              placeholder="API ключ"
            />
            <button onClick={() => handleUpdateProject(project._id)}>💾</button>
            <button onClick={() => setEditingProjectId(null)}>❌</button>
          </>          
          ) : (
            <>
              <button onClick={() => handleProjectClick(project._id)}>
                {selectedProjectId === project._id ? '▼' : '►'} {project.name}
              </button>
              <button onClick={() => {
                setEditingProjectId(project._id);
                setEditingProjectName(project.name);
                setEditingProjectApiKey(project.apiKey || '');
              }}>
                ✏️
              </button>
              <button onClick={() => handleDeleteProject(project._id)} style={{ color: 'red' }}>🗑️</button>
              <select
                value={project.status}
                onChange={(e) => handleSetProjectStatus(project._id, e.target.value)}
              >
                <option value="active">Active</option>
                <option value="inactive">Inactive</option>
                <option value="archived">Archived</option>
              </select>
            </>
          )}

          {selectedProjectId === project._id && (
            <div style={{ paddingLeft: '1rem', marginTop: '10px' }}>
              <CreateTaskForm token={token} projectId={project._id} onTaskCreated={handleTaskCreated} />
              <h4>Задачі:</h4>
              {tasks[project._id]?.length > 0 ? (
                tasks[project._id].map((task) => (
                  <div key={task._id} style={{ marginBottom: '1rem' }}>
                    {editingTaskId === task._id ? (
                      <div>
                        <input
                          placeholder="Назва"
                          value={editingTaskName}
                          onChange={(e) => setEditingTaskName(e.target.value)}
                        />
                        <input
                          placeholder="API метод"
                          value={editingTaskApiMethod}
                          onChange={(e) => setEditingTaskApiMethod(e.target.value)}
                        />.
                        <input
                          placeholder="Опис"
                          value={editingTaskDescription}
                          onChange={(e) => setEditingTaskDescription(e.target.value)}
                        />
                        <button onClick={() => handleUpdateTask(task._id, project._id)}>💾</button>
                        <button onClick={() => setEditingTaskId(null)}>❌</button>
                      </div>
                    ) : (
                      <>
                        <button onClick={() => handleTaskClick(task._id)}>
                          {expandedTasks[task._id] ? '▼' : '►'} {task.name}
                        </button>
                        <button onClick={() => {
                          setEditingTaskId(task._id);
                          setEditingTaskName(task.name);
                          setEditingTaskApiMethod(task.apiMethod || '');
                          setEditingTaskDescription(task.description || '');
                        }}>✏️</button>
                        <button onClick={() => handleDeleteTask(task._id, project._id)} style={{ color: 'red' }}>🗑️</button>
                        <select
                          value={task.status}
                          onChange={(e) => handleSetTaskStatus(task._id, e.target.value, project._id)}
                        >
                          <option value="inactive">Inactive</option>
                          <option value="active">Active</option>
                          <option value="archived">Archived</option>
                        </select>
                      </>
                    )}

                    {expandedTasks[task._id] && (
                      <div style={{ paddingLeft: '1rem', marginTop: '5px' }}>
                        <CreatePromptForm token={token} taskId={task._id} onPromptCreated={() => fetchPrompts(task._id)} />
                        <h5>Промпти:</h5>
                        {prompts[task._id]?.length > 0 ? (
                          <DndContext
                          sensors={sensors} 
                          collisionDetection={closestCenter}
                          onDragEnd={({ active, over }) => {
                            if (active.id !== over?.id) {
                              const oldIndex = prompts[task._id].findIndex(p => p._id === active.id);
                              const newIndex = prompts[task._id].findIndex(p => p._id === over?.id);
                              const reordered = arrayMove(prompts[task._id], oldIndex, newIndex);
                        
                              setPrompts(prev => ({
                                ...prev,
                                [task._id]: reordered,
                              }));
                        
                              fetch(`${import.meta.env.VITE_API_URL}/api/prompts/reorder`, {
                                method: 'POST',
                                headers: {
                                  'Content-Type': 'application/json',
                                  Authorization: `Bearer ${token}`,
                                },
                                body: JSON.stringify(reordered.map((p, index) => ({
                                  _id: p._id,
                                  executionOrder: index + 1,
                                }))),
                              });
                            }
                          }}
                        >
                          <SortableContext
                            items={prompts[task._id].map(p => p._id)}
                            strategy={verticalListSortingStrategy}
                          >
                            {prompts[task._id].map((prompt) => (
                              <SortablePromptItem key={prompt._id} id={prompt._id}>
                                {/* Весь вміст li перенеси сюди */}
                                <div>
                                  <strong>Назва:</strong>{' '}
                                  {editingPromptIds[prompt._id] ? (
                                    <>
                                      <input
                                        style={{ width: '60%' }}
                                        value={editedPromptNames[prompt._id] ?? prompt.name}
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
                                          setEditingPromptIds((prev) => ({ ...prev, [prompt._id]: false }));
                                        }}
                                      >
                                        💾
                                      </button>
                                      <button
                                        onClick={() =>
                                          setEditingPromptIds((prev) => ({ ...prev, [prompt._id]: false }))
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
                                          setEditingPromptIds((prev) => ({ ...prev, [prompt._id]: true }))
                                        }
                                      >
                                        ✏️
                                      </button>
                                    </>
                                  )}
                                </div>
                        
                                <div style={{ marginTop: '4px', fontSize: '0.9em', color: '#666' }}>
                                  <strong>Версія:</strong> {prompt.version}
                                </div>
                        
                                <div>
                                  <strong>Промпт:</strong><br />
                                  <textarea
                                    style={{ width: '100%', minHeight: '60px' }}
                                    value={editedPromptTexts[prompt._id] ?? prompt.promptText}
                                    onChange={(e) =>
                                      setEditedPromptTexts((prev) => ({
                                        ...prev,
                                        [prompt._id]: e.target.value,
                                      }))
                                    }
                                  />
                                </div>
                        
                                {prompt.responseText && (
                                  <div><strong>Відповідь:</strong> {prompt.responseText}</div>
                                )}
                        
                        {prompt.history?.length > 0 && (
                          <div style={{ marginTop: '5px' }}>
                            <strong>Історія:</strong>
                            <ul>
                            {prompt.history.map((entry, index) => (
                              <li key={index}>
                                <div><em>{new Date(entry.createdAt).toLocaleString()}</em></div>
                                <div><strong>Версія:</strong> {entry.version ?? '—'}</div>
                                <div><strong>Промпт:</strong> {entry.prompt}</div>
                                <div><strong>Відповідь:</strong> {entry.response}</div>
                                <button
                                  onClick={() => handleRollback(prompt._id, entry.version, task._id)} 
                                  style={{ marginTop: '5px' }}
                                >
                                  🔄 Повернутись до V{entry.version}
                                </button>
                              </li>
                            ))}

                            </ul>
                          </div>
                        )}

                        
                                <button onClick={() => handleSavePrompt(prompt._id)} style={{ marginRight: '10px' }}>
                                  💾 Зберегти V{prompt.version}
                                </button>
                                <button onClick={() => handleRunPrompt(prompt._id, task._id)}>
                                  🚀 Запустити
                                </button>
                                <button
                                  onClick={() => handleDeletePrompt(prompt._id, task._id)}
                                  style={{ marginTop: '5px', color: 'red' }}
                                >
                                  🗑️ Видалити
                                </button>
                              </SortablePromptItem>
                            ))}
                          </SortableContext>
                        </DndContext>
                        
                        ) : (
                          <p>Немає промптів</p>
                        )}
                      </div>
                    )}
                  </div>
                ))
              ) : (
                <p>Немає задач</p>
              )}
            </div>
          )}
        </div>
      ))}
    </div>
  );
}

