import React from 'react';
import ApiKeyItem from './ApiKeyItem.jsx';

const ApiKeyList = ({ keys, loading, onChange }) => {
  if (loading) return <p>Loading keys...</p>;
  if (!keys || !keys.length) return <p>No API keys found.</p>;

  return (
    <div>
      <h3>All API Keys</h3>
      {keys.map((key) => (
        <ApiKeyItem key={key.id} keyData={key} onChange={onChange} />
      ))}
    </div>
  );
};

export default ApiKeyList;
