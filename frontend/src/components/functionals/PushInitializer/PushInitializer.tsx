'use client';

import { css } from 'styled-system/css';
import { useSetup } from './useSetup';

export const PushInitializer = () => {
  const { shouldShowPrompt, subscribeToPush } = useSetup();

  if (!shouldShowPrompt) {
    return null;
  }

  return (
    <div
      className={css({
        position: 'fixed',
        bottom: '4',
        right: '4',
        bg: 'white',
        border: '1px solid',
        borderColor: 'orange.400',
        p: '4',
        borderRadius: 'lg',
        boxShadow: 'xl',
        zIndex: 2147483647,
        display: 'flex',
        flexDirection: 'column',
        gap: '3',
        maxWidth: '300px',
      })}
    >
      <div className={css({ display: 'flex', gap: '2', alignItems: 'flex-start' })}>
        <span className={css({ fontSize: 'xl' })}>🔔</span>
        <div className={css({ display: 'flex', flexDirection: 'column', gap: '1' })}>
          <p className={css({ fontSize: 'sm', fontWeight: 'bold', color: 'gray.800' })}>
            通知を有効にしてください
          </p>
          <p className={css({ fontSize: 'xs', color: 'gray.600', lineHeight: '1.4' })}>
            セリの高値更新情報をリアルタイムで受け取るには通知の許可が必要です。
          </p>
        </div>
      </div>
      <button
        onClick={() => subscribeToPush()}
        className={css({
          bg: 'orange.500',
          _hover: { bg: 'orange.600' },
          color: 'white',
          py: '2',
          px: '4',
          borderRadius: 'md',
          fontSize: 'xs',
          fontWeight: 'bold',
          cursor: 'pointer',
          transition: 'background 0.2s',
          textAlign: 'center',
        })}
      >
        通知を有効にする
      </button>
    </div>
  );
};
