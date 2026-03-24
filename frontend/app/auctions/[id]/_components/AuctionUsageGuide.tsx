'use client';

import { Box, Text } from '@atoms';
import { css } from 'styled-system/css';

interface AuctionUsageGuideProps {
  t: (key: string) => string;
}

export const AuctionUsageGuide = ({ t }: AuctionUsageGuideProps) => {
  return (
    <Box
      className={css({
        mb: '8',
        p: '6',
        bg: 'blue.50',
        border: '1px solid',
        borderColor: 'blue.200',
        borderRadius: 'xl',
      })}
    >
      <Text
        as="h2"
        className={css({ fontSize: 'lg', fontWeight: 'bold', color: 'blue.900', mb: '2' })}
      >
        {t('Public.AuctionDetail.Usage.title')}
      </Text>
      <ol
        className={css({
          listStyleType: 'decimal',
          listStylePosition: 'inside',
          spaceY: '1',
          fontSize: 'sm',
          color: 'blue.800',
        })}
      >
        <li>{t('Public.AuctionDetail.Usage.step1')}</li>
        <li>{t('Public.AuctionDetail.Usage.step2')}</li>
        <li>{t('Public.AuctionDetail.Usage.step3')}</li>
        <li>{t('Public.AuctionDetail.Usage.step4')}</li>
        <li>{t('Public.AuctionDetail.Usage.step5')}</li>
      </ol>
    </Box>
  );
};
