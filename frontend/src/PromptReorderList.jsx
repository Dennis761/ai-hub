import React, { useMemo } from 'react';
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
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';

export default function PromptReorderList({
  taskId,
  token,
  items,
  onLocalReorder,
  refresh,
  SortablePromptItem,
  renderItem,
  onError,
}) {
  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } })
  );

  const itemIds = useMemo(() => items.map((p) => p._id), [items]);

  const handleDragEnd = async ({ active, over }) => {
    if (!over || active.id === over.id) return;

    const oldIndex = itemIds.indexOf(active.id);
    const newIndex = itemIds.indexOf(over.id);
    if (oldIndex < 0 || newIndex < 0) return;

    const reordered = arrayMove(items, oldIndex, newIndex);
    onLocalReorder?.(reordered);

    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/prompts/reorder`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ items: reordered.map((p) => p._id) }),
      });

      let data = null;
      try {
        data = await res.json();
      } catch {
        data = null;
      }

      if (!res.ok) {
        const backendError =
          data && typeof data.error === 'string' ? data.error : 'Failed to reorder prompts';
        console.error('Reorder error:', backendError);
        onError?.(backendError);
      }
    } catch (err) {
      console.error('Reorder failed:', err);
      onError?.(err?.message || 'Network error');
    }
  };

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
      onDragEnd={handleDragEnd}
    >
      <SortableContext items={itemIds} strategy={verticalListSortingStrategy}>
        <ul style={{ listStyle: 'none', padding: 0, margin: 0 }}>
          {items.map((prompt) => (
            <SortablePromptItem key={prompt._id} id={prompt._id}>
              {renderItem(prompt)}
            </SortablePromptItem>
          ))}
        </ul>
      </SortableContext>
    </DndContext>
  );
}
