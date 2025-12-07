import { Box } from '@/src/core/ui';
import { Sidebar } from './_components/Sidebar';

export default function AdminLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    return (
        <Box display="flex" minH="screen" bg="gray.100">
            {/* Sidebar */}
            <Sidebar />

            {/* Main Content */}
            <Box as="main" flex="1" p="8" overflowY="auto">
                {children}
            </Box>
        </Box>
    );
}
