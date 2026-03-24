'use client';

import { use } from 'react';
import { useAuctionDetailPage } from './_hooks/useAuctionDetailPage';
import { AuctionHeader } from './_components/AuctionHeader';
import { BuyerLoginForm } from './_components/BuyerLoginForm';
import { AuctionUsageGuide } from './_components/AuctionUsageGuide';
import { AuctionItemList } from './_components/AuctionItemList';
import { BiddingPanel } from './_components/BiddingPanel';
import { Box, Text, Card } from '@atoms';
import { css } from 'styled-system/css';

export default function AuctionDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const auctionId = Number(id);

  const {
    auction,
    items,
    isLoading,
    isChecking,
    isLoggedIn,
    selectedItem,
    selectedItemId,
    auctionActive,
    message,
    loginError,
    isBidLoading,
    bidForm,
    loginForm,
    onSelectItem,
    onSubmitLogin,
    onSubmitBid,
    t,
  } = useAuctionDetailPage(auctionId);

  if (isNaN(auctionId)) {
    return <Box>Invalid Auction ID</Box>;
  }

  if (isChecking || isLoading) {
    return (
      <Box minH="screen" display="flex" alignItems="center" justifyContent="center" bg="gray.50">
        <Text fontSize="xl" className={css({ color: 'gray.700' })}>
          {t('Common.loading')}
        </Text>
      </Box>
    );
  }

  if (!auction) {
    return <Box>{t('Common.no_data')}</Box>;
  }

  // ログインしていない場合はログインフォームを表示
  if (!isLoggedIn) {
    return (
      <BuyerLoginForm
        loginForm={loginForm}
        onSubmit={onSubmitLogin}
        loginError={loginError}
        t={t}
      />
    );
  }

  return (
    <Box minH="screen" bg="gray.50" p={{ base: '4', md: '8' }}>
      <Box maxW="7xl" mx="auto">
        <AuctionHeader auction={auction} t={t} />

        {message && (
          <Card
            mb="6"
            padding="md"
            className={css({
              borderLeft: '4px solid',
              borderColor: 'green.500',
              bg: 'green.50',
              animation: 'bounce 1s infinite',
            })}
          >
            <Text fontWeight="bold" className={css({ color: 'green.700' })}>
              {message}
            </Text>
          </Card>
        )}

        <AuctionUsageGuide t={t} />

        <Box
          display="grid"
          gridTemplateColumns={{ base: 'repeat(1, 1fr)', lg: 'repeat(3, 1fr)' }}
          gap="8"
        >
          <AuctionItemList
            items={items || []}
            selectedItemId={selectedItemId}
            onSelectItem={onSelectItem}
            t={t}
          />

          <BiddingPanel
            selectedItem={selectedItem}
            auction={auction}
            auctionActive={auctionActive}
            bidForm={bidForm}
            onSubmitBid={onSubmitBid}
            isBidLoading={isBidLoading}
            t={t}
          />
        </Box>
      </Box>
    </Box>
  );
}
