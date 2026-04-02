import { AdminLayoutTemplate } from '@templates';
import { AuthorizableAdminSidebar } from '@/src/features/auth/components/AuthorizableAdminSidebar';

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  return <AdminLayoutTemplate sidebar={<AuthorizableAdminSidebar />}>{children}</AdminLayoutTemplate>;
}
