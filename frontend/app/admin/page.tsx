'use client';

import Link from 'next/link';
import { css } from 'styled-system/css';
import { Box, HStack, Text, Card } from '@/src/core/ui';
import { useTranslations } from 'next-intl';

export default function AdminDashboard() {
    const t = useTranslations();

    const menuItems = [
        {
            title: t('Admin.Dashboard.fishermen_title'),
            description: t('Admin.Dashboard.fishermen_desc'),
            href: '/admin/fishermen',
            icon: '🎣',
            color: 'blue',
        },
        {
            title: t('Admin.Dashboard.buyers_title'),
            description: t('Admin.Dashboard.buyers_desc'),
            href: '/admin/buyers',
            icon: '🛒',
            color: 'green',
        },
        {
            title: t('Admin.Dashboard.items_title'),
            description: t('Admin.Dashboard.items_desc'),
            href: '/admin/items',
            icon: '🐟',
            color: 'indigo',
        },
        {
            title: t('Admin.Dashboard.venues_title'),
            description: t('Admin.Dashboard.venues_desc'),
            href: '/admin/venues',
            icon: '📍',
            color: 'purple',
        },
        {
            title: t('Admin.Dashboard.auctions_title'),
            description: t('Admin.Dashboard.auctions_desc'),
            href: '/admin/auctions',
            icon: '🔨',
            color: 'orange',
        },
        {
            title: t('Admin.Dashboard.invoice_title'),
            description: t('Admin.Dashboard.invoice_desc'),
            href: '/invoice',
            icon: '💰',
            color: 'yellow',
        },
    ];

    const colorStyles: Record<string, { bg: string; hover: string; iconBg: string; iconText: string }> = {
        indigo: { bg: 'indigo.50', hover: 'indigo.100', iconBg: 'indigo.100', iconText: 'indigo.600' },
        green: { bg: 'green.50', hover: 'green.100', iconBg: 'green.100', iconText: 'green.600' },
        orange: { bg: 'orange.50', hover: 'orange.100', iconBg: 'orange.100', iconText: 'orange.600' },
        blue: { bg: 'blue.50', hover: 'blue.100', iconBg: 'blue.100', iconText: 'blue.600' },
        purple: { bg: 'purple.50', hover: 'purple.100', iconBg: 'purple.100', iconText: 'purple.600' },
        yellow: { bg: 'yellow.50', hover: 'yellow.100', iconBg: 'yellow.100', iconText: 'yellow.600' },
    };

    const getStyles = (color: string) => colorStyles[color] || colorStyles.indigo;

    return (
        <Box className={css({ maxW: '7xl', mx: 'auto', p: '6' })}>
            <Box className={css({ mb: '8' })}>
                <Text variant="h1" className={css({ fontSize: '3xl', fontWeight: 'bold', color: 'gray.800' })}>{t('Admin.Dashboard.title')}</Text>
                <Text className={css({ color: 'gray.600', mt: '2' })}>{t('Admin.Dashboard.subtitle')}</Text>
            </Box>

            <div className={css({ display: 'grid', gridTemplateColumns: { base: 'repeat(1, 1fr)', md: 'repeat(2, 1fr)', lg: 'repeat(3, 1fr)' }, gap: '6' })}>
                {menuItems.map((item) => {
                    const styles = getStyles(item.color);
                    return (
                        <Link
                            key={item.href}
                            href={item.href}
                            className={css({ textDecoration: 'none', display: 'block', h: 'full' })}
                        >
                            <Card
                                className={css({
                                    p: '6',
                                    h: 'full',
                                    borderWidth: '2px',
                                    borderColor: 'transparent',
                                    bg: styles.bg,
                                    transition: 'all 0.2s',
                                    display: 'flex',
                                    flexDirection: 'column',
                                    _hover: {
                                        bg: styles.hover,
                                        shadow: 'lg',
                                        transform: 'scale(1.05)',
                                    }
                                })}
                            >
                                <HStack spacing="4" align="start" className={css({ flex: '1' })}>
                                    <Box className={css({ p: '3', borderRadius: 'lg', bg: styles.iconBg })}>
                                        <span className={css({ fontSize: '3xl' })}>{item.icon}</span>
                                    </Box>
                                    <Box className={css({ flex: '1' })}>
                                        <Text variant="h3" className={css({ fontSize: 'lg', fontWeight: 'bold', color: 'gray.900', mb: '1' })}>{item.title}</Text>
                                        <Text variant="small" className={css({ color: 'gray.600' })}>{item.description}</Text>
                                    </Box>
                                </HStack>
                            </Card>
                        </Link>
                    );
                })}
            </div>

            <Box className={css({ mt: '12', p: '6', bg: 'blue.50', border: '1px solid', borderColor: 'blue.200', borderRadius: 'xl' })}>
                <Text variant="h2" className={css({ fontSize: 'lg', fontWeight: 'bold', color: 'blue.900', mb: '2' })}>{t('Admin.Dashboard.usage_title')}</Text>
                <ol className={css({ listStyleType: 'decimal', listStylePosition: 'inside', spaceY: '1', fontSize: 'sm', color: 'blue.800' })}>
                    <li>{t('Admin.Dashboard.usage_step1')}</li>
                    <li>{t('Admin.Dashboard.usage_step2')}</li>
                    <li>{t('Admin.Dashboard.usage_step3')}</li>
                    <li>{t('Admin.Dashboard.usage_step4')}</li>
                    <li>{t('Admin.Dashboard.usage_step5')}</li>
                    <li>{t('Admin.Dashboard.usage_step6')}</li>
                </ol>
            </Box>
        </Box>
    );
}
