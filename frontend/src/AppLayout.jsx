import { Link, Outlet } from 'react-router-dom';

export default function AppLayout() {
  return (
    <div>
      <nav style={{ padding: '1rem', backgroundColor: '#f0f0f0' }}>
        <Link to="/ai-hub/keys" style={{ marginRight: '1rem' }}>API Keys</Link>
        <Link to="/ai-hub/projects" style={{ marginRight: '1rem' }}>Projects</Link>
        <Link to="/ai-hub/stats">Statistics</Link>
      </nav>

      <div style={{ padding: '1rem' }}>
        <Outlet />
      </div>
    </div>
  );
}
