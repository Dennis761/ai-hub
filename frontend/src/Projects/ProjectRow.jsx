export default function ProjectRow({
  project,
  selectedProjectId,
  editingProjectId,
  editingProjectName,
  editingProjectApiKey,
  setEditingProjectId,
  setEditingProjectName,
  setEditingProjectApiKey,
  onToggle,
  onDelete,
  onUpdate,
  onStatusChange,
}) {
  const isEditing = editingProjectId === project._id;

  if (isEditing) {
    return (
      <>
        <input
          value={editingProjectName}
          onChange={(e) => setEditingProjectName(e.target.value)}
          placeholder="Name"
        />
        <input
          value={editingProjectApiKey}
          onChange={(e) => setEditingProjectApiKey(e.target.value)}
          placeholder="API key"
        />
        <button onClick={onUpdate}>💾</button>
        <button onClick={() => setEditingProjectId(null)}>❌</button>
      </>
    );
  }

  return (
    <>
      <button onClick={onToggle}>
        {selectedProjectId === project._id ? '▼' : '►'} {project.name}
      </button>
      <button
        onClick={() => {
          setEditingProjectId(project._id);
          setEditingProjectName(project.name);
          setEditingProjectApiKey(project.apiKey || '');
        }}
      >
        ✏️
      </button>
      <button onClick={onDelete} style={{ color: 'red' }}>
        🗑️
      </button>
      <select value={project.status} onChange={(e) => onStatusChange(e.target.value)}>
        <option value="active">Active</option>
        <option value="inactive">Inactive</option>
        <option value="archived">Archived</option>
      </select>
    </>
  );
}
