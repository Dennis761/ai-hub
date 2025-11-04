export default function PromptList({
  items,
  editingPromptIds,
  editedPromptTexts,
  editedPromptNames,
  setEditingPromptIds,
  setEditedPromptTexts,
  setEditedPromptNames,
  onSavePrompt,
  onRunPrompt,
  onDeletePrompt,
  onRollback,
  SortablePromptItem,
}) {
  return (
    <>
      {items.map((prompt) => (
        <SortablePromptItem key={prompt._id} id={prompt._id}>
          <div>
            <strong>Name:</strong>{' '}
            {editingPromptIds[prompt._id] ? (
              <>
                <input
                  style={{ width: '60%' }}
                  value={editedPromptNames[prompt._id] ?? prompt.name}
                  onChange={(e) =>
                    setEditedPromptNames((prev) => ({ ...prev, [prompt._id]: e.target.value }))
                  }
                />
                <button
                  onClick={() => {
                    onSavePrompt(prompt._id);
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
            <strong>Version:</strong> {prompt.version}
          </div>

          <div>
            <strong>Prompt:</strong>
            <br />
            <textarea
              style={{ width: '100%', minHeight: '60px' }}
              value={editedPromptTexts[prompt._id] ?? prompt.promptText}
              onChange={(e) =>
                setEditedPromptTexts((prev) => ({ ...prev, [prompt._id]: e.target.value }))
              }
            />
          </div>

          {prompt.responseText !== null && prompt.responseText !== undefined && (
            <div>
              <strong>Response:</strong> {prompt.responseText}
            </div>
          )}

          {prompt.history?.length > 0 && (
            <div style={{ marginTop: '5px' }}>
              <strong>History:</strong>
              <ul>
                {prompt.history.map((entry, index) => (
                  <li key={index}>
                    <div>
                      <em>{new Date(entry.createdAt).toLocaleString()}</em>
                    </div>
                    <div>
                      <strong>Version:</strong> {entry.version ?? '—'}
                    </div>
                    <div>
                      <strong>Prompt:</strong> {entry.prompt}
                    </div>
                    <div>
                      <strong>Response:</strong> {entry.response}
                    </div>
                    <button
                      onClick={() => onRollback(prompt._id, entry.version)}
                      style={{ marginTop: '5px' }}
                    >
                      🔄 Roll back to V{entry.version}
                    </button>
                  </li>
                ))}
              </ul>
            </div>
          )}

          <button
            onClick={() => onSavePrompt(prompt._id)}
            style={{ marginRight: '10px' }}
          >
            💾 Save V{prompt.version}
          </button>
          <button onClick={() => onRunPrompt(prompt._id)}>🚀 Run</button>
          <button
            onClick={() => onDeletePrompt(prompt._id)}
            style={{ marginTop: '5px', color: 'red' }}
          >
            🗑️ Delete
          </button>
        </SortablePromptItem>
      ))}
    </>
  );
}
