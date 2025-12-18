import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { Select } from './select';

describe('Select', () => {
    it('renders correctly', () => {
        render(
            <Select>
                <option value="1">Option 1</option>
                <option value="2">Option 2</option>
            </Select>
        );
        expect(screen.getByRole('combobox')).toBeInTheDocument();
        expect(screen.getByText('Option 1')).toBeInTheDocument();
    });

    it('handles onChange events', () => {
        const handleChange = vi.fn();
        render(
            <Select onChange={handleChange}>
                <option value="1">Option 1</option>
                <option value="2">Option 2</option>
            </Select>
        );
        const select = screen.getByRole('combobox');
        fireEvent.change(select, { target: { value: '2' } });
        expect(handleChange).toHaveBeenCalledTimes(1);
    });
});
