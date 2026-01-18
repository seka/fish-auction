import { Box } from '@/src/core/ui';
import { ToastContainer } from '../atoms/Toast';

type MainLayoutTemplateProps = {
    navbar: React.ReactNode;
    children: React.ReactNode;
};

export const MainLayoutTemplate = ({ navbar, children }: MainLayoutTemplateProps) => {
    return (
        <Box minH="screen" display="flex" flexDirection="column" bg="gray.50">
            {navbar}
            <Box flex="1">
                {children}
            </Box>
            <ToastContainer />
        </Box>
    );
};
