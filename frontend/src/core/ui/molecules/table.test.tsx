import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { Table, Thead, Tbody, Tr, Th, Td } from './table';

describe('Table Molecules', () => {
    it('renders table structure correctly', () => {
        render(
            <Table>
                <Thead>
                    <Tr>
                        <Th>Header 1</Th>
                        <Th>Header 2</Th>
                    </Tr>
                </Thead>
                <Tbody>
                    <Tr>
                        <Td>Cell 1</Td>
                        <Td>Cell 2</Td>
                    </Tr>
                </Tbody>
            </Table>
        );

        expect(screen.getByRole('table')).toBeInTheDocument();
        expect(screen.getByText('Header 1')).toBeInTheDocument();
        expect(screen.getByText('Cell 1')).toBeInTheDocument();

        // Verify semantic structure
        const rows = screen.getAllByRole('row');
        expect(rows).toHaveLength(2); // 1 header row + 1 body row
    });
});
