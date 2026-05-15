import { css } from 'styled-system/css';

const wrapper = css({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  minHeight: '50vh',
});

const spinner = css({
  width: '8',
  height: '8',
  border: '3px solid',
  borderColor: 'gray.200',
  borderTopColor: 'gray.500',
  borderRadius: 'full',
  animation: 'spin 0.8s linear infinite',
});

export const Spinner = () => (
  <div className={wrapper} role="status" aria-busy="true" aria-label="Loading">
    <span className={spinner} />
  </div>
);
