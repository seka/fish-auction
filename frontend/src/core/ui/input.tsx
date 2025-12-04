import { cva, type RecipeVariantProps } from 'styled-system/css';
import { styled } from 'styled-system/jsx';

const inputRecipe = cva({
    base: {
        display: 'block',
        width: 'full',
        px: '3',
        py: '2',
        bg: 'white',
        border: '1px solid',
        borderColor: 'gray.300',
        borderRadius: 'md',
        fontSize: 'sm',
        outline: 'none',
        transition: 'border-color 0.2s',
        _focus: {
            borderColor: 'primary.500',
            ring: '1px',
            ringColor: 'primary.500',
        },
        _disabled: {
            bg: 'gray.50',
            cursor: 'not-allowed',
        },
        _placeholder: {
            color: 'gray.500',
        },
    },
    variants: {
        error: {
            true: {
                borderColor: 'red.500',
                _focus: {
                    borderColor: 'red.500',
                    ringColor: 'red.500',
                },
            },
        },
    },
});

export type InputProps = RecipeVariantProps<typeof inputRecipe> &
    React.InputHTMLAttributes<HTMLInputElement>;

export const Input = styled('input', inputRecipe);
