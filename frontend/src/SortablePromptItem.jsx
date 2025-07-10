import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';

export function SortablePromptItem({ id, children }) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    background: '#f9f9f9',
    border: '1px solid #ccc',
    padding: '10px',
    marginBottom: '10px',
    borderRadius: '6px',
    opacity: isDragging ? 0.5 : 1,
  };

  return (
    <div ref={setNodeRef} style={style}>
      {/* ручка для перетягування */}
      <div
        {...attributes}
        {...listeners}
        style={{
          cursor: 'grab',
          background: '#ddd',
          padding: '5px',
          borderRadius: '4px',
          marginBottom: '8px',
          fontSize: '0.9em',
          userSelect: 'none',
        }}
      >
        ☰ Перетягнути
      </div>

      {/* сам вміст, який можна редагувати */}
      {children}
    </div>
  );
}
