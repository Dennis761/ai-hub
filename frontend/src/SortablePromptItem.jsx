import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';

export function SortablePromptItem({ id, children }) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } =
    useSortable({ id });

  const wrapperStyle = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.85 : 1,
    listStyle: 'none',
  };

  const rowStyle = {
    display: 'flex',
    gap: 12,
    alignItems: 'flex-start',
    padding: '10px 12px',
  };

  const handleStyle = {
    cursor: 'grab',
    userSelect: 'none',
    width: 20,
    minWidth: 20,
    height: 24,
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'space-between',
    padding: '2px 0',
    opacity: 0.6,
  };

  const barStyle = {
    height: 2,
    width: '100%',
    background: 'rgba(0,0,0,0.35)',
    borderRadius: 2,
  };

  const separatorStyle = {
    height: 1,
    background: 'rgba(0,0,0,0.08)',
    margin: 0,
  };

  return (
    <li ref={setNodeRef} style={wrapperStyle}>
      <div style={rowStyle}>
        <div
          {...attributes}
          {...listeners}
          style={handleStyle}
          aria-label="Drag"
          title="Drag"
        >
          <div style={barStyle} />
          <div style={barStyle} />
          <div style={barStyle} />
        </div>

        <div style={{ flex: 1, minWidth: 0 }}>
          {children}
        </div>
      </div>

      <div style={separatorStyle} />
    </li>
  );
}

export default SortablePromptItem;
