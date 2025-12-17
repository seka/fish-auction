import Link from 'next/link';
import Image from 'next/image';
import { css } from 'styled-system/css';
import { Box, Stack, Text, Card } from '@/src/core/ui';
import { COMMON_TEXT_KEYS } from '@/src/core/assets/text';
import { useTranslations } from 'next-intl';

export default function Home() {
  const t = useTranslations();
  return (
    <Box
      className={css({
        minH: '50vh',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        p: '8',
        fontFamily: 'sans',
        bgGradient: 'to-br',
        gradientFrom: 'blue.50',
        gradientTo: 'indigo.50',
      })}
    >
      <Image src="/logo_text.png" alt="FISHING AUCTION Logo" width={300} height={300} className={css({ mx: 'auto' })} />

      <div className={css({ display: 'grid', gridTemplateColumns: { base: 'repeat(1, 1fr)', md: 'repeat(2, 1fr)' }, gap: '8', w: 'full', maxW: '5xl' })}>
        {/* Admin Portal */}
        <Link href="/admin" className={`group ${css({ textDecoration: 'none', display: 'block', h: 'full' })}`}>
          <Card variant="interactive" className={css({ position: 'relative', overflow: 'hidden', p: '10', h: 'full', display: 'flex', flexDirection: 'column' })}>
            <Box className={css({ position: 'absolute', top: '0', right: '0', w: '32', h: '32', bg: 'primary.50', borderBottomLeftRadius: 'full', mr: '-8', mt: '-8', transition: 'transform 0.3s', _groupHover: { transform: 'scale(1.1)' } })} />
            <Stack spacing="6" align="center" className={css({ position: 'relative', zIndex: '10', textAlign: 'center', flex: '1' })}>
              {/* Icon */}
              <Box className={css({ p: '5', bg: 'primary.100', borderRadius: '2xl', color: 'primary.600', transition: 'colors 0.3s', _groupHover: { bg: 'primary.600', color: 'white' } })}>
                <svg className={css({ w: '14', h: '14' })} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
              </Box>
              <Box>
                <Text variant="h2" className={css({ mb: '2' })}>{t(COMMON_TEXT_KEYS.admin_panel)}</Text>
                <Text className={css({ color: 'gray.500', lineHeight: 'relaxed' })}>
                  <div dangerouslySetInnerHTML={{ __html: t.raw(COMMON_TEXT_KEYS.admin_panel_description) }} />
                </Text>
              </Box>
              <Box className={css({ display: 'inline-flex', alignItems: 'center', color: 'primary.600', fontWeight: 'bold', transition: 'transform 0.3s', _groupHover: { transform: 'translateX(4px)' } })}>
                {t(COMMON_TEXT_KEYS.go_to_admin)} <span className={css({ ml: '2' })}>&rarr;</span>
              </Box>
            </Stack>
          </Card>
        </Link>

        {/* Auction Floor */}
        <Link href="/auctions" className={`group ${css({ textDecoration: 'none', display: 'block', h: 'full' })}`}>
          <Card variant="interactive" className={css({ position: 'relative', overflow: 'hidden', p: '10', h: 'full', display: 'flex', flexDirection: 'column' })}>
            <Box className={css({ position: 'absolute', top: '0', right: '0', w: '32', h: '32', bg: 'secondary.50', borderBottomLeftRadius: 'full', mr: '-8', mt: '-8', transition: 'transform 0.3s', _groupHover: { transform: 'scale(1.1)' } })} />
            <Stack spacing="6" align="center" className={css({ position: 'relative', zIndex: '10', textAlign: 'center', flex: '1' })}>
              {/* Icon */}
              <Box className={css({ p: '5', bg: 'secondary.100', borderRadius: '2xl', color: 'secondary.600', transition: 'colors 0.3s', _groupHover: { bg: 'secondary.600', color: 'white' } })}>
                <svg className={css({ w: '14', h: '14' })} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122" />
                </svg>
              </Box>
              <Box>
                <Text variant="h2" className={css({ mb: '2' })}>{t(COMMON_TEXT_KEYS.auction_venue)}</Text>
                <Text className={css({ color: 'gray.500', lineHeight: 'relaxed' })}>
                  <div dangerouslySetInnerHTML={{ __html: t.raw(COMMON_TEXT_KEYS.auction_venue_description) }} />
                </Text>
              </Box>
              <Box className={css({ display: 'inline-flex', alignItems: 'center', color: 'secondary.600', fontWeight: 'bold', transition: 'transform 0.3s', _groupHover: { transform: 'translateX(4px)' } })}>
                {t(COMMON_TEXT_KEYS.enter_venue)} <span className={css({ ml: '2' })}>&rarr;</span>
              </Box>
            </Stack>
          </Card>
        </Link>
      </div>
    </Box>
  );
}
