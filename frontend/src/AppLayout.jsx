import { Link, Outlet } from 'react-router-dom';

export default function AppLayout() {
  return (
    <div>
      <nav style={{ padding: '1rem', backgroundColor: '#f0f0f0' }}>
        <Link to="/ai-hub/keys" style={{ marginRight: '1rem' }}>API Ключі</Link>
        <Link to="/ai-hub/projects" style={{ marginRight: '1rem' }}>Проєкти</Link>
        <Link to="/ai-hub/stats">Статистика</Link>
      </nav>

      <div style={{ padding: '1rem' }}>
        <Outlet /> 
      </div>
    </div>
  );
}
