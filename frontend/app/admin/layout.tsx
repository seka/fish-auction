import { Box } from '@atoms';
import { Sidebar } from '@organisms';
import { AdminLayoutTemplate } from '@templates';

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  return <AdminLayoutTemplate sidebar={<Sidebar />}>{children}</AdminLayoutTemplate>;
}
