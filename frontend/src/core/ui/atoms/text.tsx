import { cva, type RecipeVariantProps } from 'styled-system/css';
import { styled } from 'styled-system/jsx';

const textRecipe = cva({
    base: {},
    variants: {
        variant: {
            h1: {
                fontSize: '3rem',
                fontWeight: 'bold',
                lineHeight: '1.2',
            },
            h2: {
                fontSize: '2.25rem',
                fontWeight: 'bold',
                lineHeight: '1.3',
            },
            h3: {
                fontSize: '1.875rem',
                fontWeight: 'bold',
                lineHeight: '1.4',
            },
            h4: {
                fontSize: '1.5rem',
                fontWeight: 'semibold',
                lineHeight: '1.4',
            },
            body: {
                fontSize: '1rem',
                lineHeight: '1.6',
            },
            small: {
                fontSize: '0.875rem',
                lineHeight: '1.5',
            },
        },
        color: {
            default: {
                color: 'gray.900',
            },
            muted: {
                color: 'gray.600',
            },
            primary: {
                color: 'primary.600',
            },
            secondary: {
                color: 'secondary.600',
            },
        },
    },
    defaultVariants: {
        variant: 'body',
        color: 'default',
    },
});

export type TextProps = RecipeVariantProps<typeof textRecipe> &
    React.HTMLAttributes<HTMLElement> & {
        as?: 'h1' | 'h2' | 'h3' | 'h4' | 'p' | 'span';
    };

export const Text = styled('p', textRecipe);
