// AdminLayout.tsx
import { Outlet } from 'react-router-dom';
import Sidebar from './Sidebar';

const AdminLayout = () => {
  return (
    <div className="flex">
      <Sidebar isOpen={true} />
      <div className="flex flex-col flex-1">
        <main>
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default AdminLayout;