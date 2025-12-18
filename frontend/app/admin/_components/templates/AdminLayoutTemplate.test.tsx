import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { AdminLayoutTemplate } from './AdminLayoutTemplate';

describe('AdminLayoutTemplate', () => {
    it('renders sidebar and children correctly', () => {
        render(
            <AdminLayoutTemplate sidebar={<div data-testid="sidebar">Sidebar</div>}>
                <div data-testid="content">Content</div>
            </AdminLayoutTemplate>
        );

        expect(screen.getByTestId('sidebar')).toBeInTheDocument();
        expect(screen.getByTestId('content')).toBeInTheDocument();
    });
});
