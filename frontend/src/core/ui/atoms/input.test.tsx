import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { Input } from './input';

describe('Input', () => {
    it('renders correctly', () => {
        render(<Input placeholder="Enter text" />);
        const input = screen.getByPlaceholderText('Enter text');
        expect(input).toBeInTheDocument();
    });

    it('handles onChange events', () => {
        const handleChange = vi.fn();
        render(<Input onChange={handleChange} />);
        const input = screen.getByRole('textbox');
        fireEvent.change(input, { target: { value: 'Hello' } });
        expect(handleChange).toHaveBeenCalledTimes(1);
    });

    it('renders with error variant', () => {
        render(<Input error />);
        const input = screen.getByRole('textbox');
        expect(input).toBeInTheDocument();
        // Visual verification via styles is implicit in implementation
    });

    it('renders as disabled', () => {
        render(<Input disabled />);
        const input = screen.getByRole('textbox');
        expect(input).toBeDisabled();
    });
});
