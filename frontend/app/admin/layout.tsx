import { AdminLayoutTemplate } from '@templates';
import { AuthorizableAdminSidebar } from '@/src/features/admin';

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  return (
    <AdminLayoutTemplate sidebar={<AuthorizableAdminSidebar />}>{children}</AdminLayoutTemplate>
  );
}
