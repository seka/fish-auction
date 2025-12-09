import { Box } from '@/src/core/ui';

type AdminLayoutTemplateProps = {
    sidebar: React.ReactNode;
    children: React.ReactNode;
};

export const AdminLayoutTemplate = ({ sidebar, children }: AdminLayoutTemplateProps) => {
    return (
        <Box display="flex" minH="screen" bg="gray.100">
            {/* Sidebar Injection */}
            {sidebar}

            {/* Main Content */}
            <Box as="main" flex="1" p="8" overflowY="auto">
                {children}
            </Box>
        </Box>
    );
};
