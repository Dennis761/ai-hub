export default function TaskBlock({
  task,
  projectId,
  expanded,
  editingTaskId,
  editingTaskName,
  editingTaskApiMethod,
  editingTaskDescription,
  setEditingTaskId,
  setEditingTaskName,
  setEditingTaskApiMethod,
  setEditingTaskDescription,
  onToggle,
  onDelete,
  onUpdate,
  onStatusChange,
  children,
}) {
  const isEditing = editingTaskId === task._id;

  return (
    <div style={{ marginBottom: '1rem' }}>
      {isEditing ? (
        <div>
          <input
            placeholder="Name"
            value={editingTaskName}
            onChange={(e) => setEditingTaskName(e.target.value)}
          />
          <input
            placeholder="API method"
            value={editingTaskApiMethod}
            onChange={(e) => setEditingTaskApiMethod(e.target.value)}
          />
          <input
            placeholder="Description"
            value={editingTaskDescription}
            onChange={(e) => setEditingTaskDescription(e.target.value)}
          />
          <button onClick={onUpdate}>💾</button>
          <button onClick={() => setEditingTaskId(null)}>❌</button>
        </div>
      ) : (
        <>
          <button onClick={onToggle}>
            {expanded ? '▼' : '►'} {task.name}
          </button>
          <button
            onClick={() => {
              setEditingTaskId(task._id);
              setEditingTaskName(task.name);
              setEditingTaskApiMethod(task.apiMethod || '');
              setEditingTaskDescription(task.description || '');
            }}
          >
            ✏️
          </button>
          <button onClick={onDelete} style={{ color: 'red' }}>
            🗑️
          </button>
          <select value={task.status} onChange={(e) => onStatusChange(e.target.value)}>
            <option value="inactive">Inactive</option>
            <option value="active">Active</option>
            <option value="archived">Archived</option>
          </select>
        </>
      )}

      {expanded && <div style={{ marginTop: '5px' }}>{children}</div>}
    </div>
  );
}
