import { Sidebar } from './_components/organisms/Sidebar';
import { AdminLayoutTemplate } from './_components/templates/AdminLayoutTemplate';

export default function AdminLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    return <AdminLayoutTemplate sidebar={<Sidebar />}>{children}</AdminLayoutTemplate>;
}
