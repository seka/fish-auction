'use client';

import Link from 'next/link';
import { css } from 'styled-system/css';
import { Box, HStack, Text, Card } from '@/src/core/ui';

export default function AdminDashboard() {
    const menuItems = [
        {
            title: '漁師管理',
            description: '漁師の登録・一覧表示',
            href: '/admin/fishermen',
            icon: '👨‍🌾',
            color: 'indigo',
        },
        {
            title: '中買人管理',
            description: '中買人の登録・一覧表示',
            href: '/admin/buyers',
            icon: '👔',
            color: 'green',
        },
        {
            title: '出品管理',
            description: 'セリへの出品登録',
            href: '/admin/items',
            icon: '🐟',
            color: 'orange',
        },
        {
            title: '会場管理',
            description: 'セリ会場の登録・管理',
            href: '/admin/venues',
            icon: '🏢',
            color: 'blue',
        },
        {
            title: 'セリ管理',
            description: 'セリの作成・スケジュール管理',
            href: '/admin/auctions',
            icon: '📅',
            color: 'purple',
        },
        {
            title: '請求書発行',
            description: '落札後の請求書発行',
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
                <Text variant="h1" className={css({ fontSize: '3xl', fontWeight: 'bold', color: 'gray.800' })}>管理ダッシュボード</Text>
                <Text className={css({ color: 'gray.600', mt: '2' })}>各管理メニューを選択してください</Text>
            </Box>

            <div className={css({ display: 'grid', gridTemplateColumns: { base: '1', md: '2', lg: '3' }, gap: '6' })}>
                {menuItems.map((item) => {
                    const styles = getStyles(item.color);
                    return (
                        <Link
                            key={item.href}
                            href={item.href}
                            className={css({ textDecoration: 'none' })}
                        >
                            <Card
                                className={css({
                                    p: '6',
                                    borderWidth: '2px',
                                    borderColor: 'transparent',
                                    bg: styles.bg,
                                    transition: 'all 0.2s',
                                    _hover: {
                                        bg: styles.hover,
                                        shadow: 'lg',
                                        transform: 'scale(1.05)',
                                    }
                                })}
                            >
                                <HStack spacing="4" align="start">
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
                <Text variant="h2" className={css({ fontSize: 'lg', fontWeight: 'bold', color: 'blue.900', mb: '2' })}>📌 使い方</Text>
                <ol className={css({ listStyleType: 'decimal', listStylePosition: 'inside', spaceY: '1', fontSize: 'sm', color: 'blue.800' })}>
                    <li>まず「会場管理」でセリを行う会場を登録します</li>
                    <li>「セリ管理」で開催日時を設定してセリを作成します</li>
                    <li>「漁師管理」「中買人管理」で参加者を登録します</li>
                    <li>「出品管理」で魚を登録してセリに出品します</li>
                    <li>セリ会場で入札が行われます</li>
                    <li>「請求書発行」で落札後の請求書を発行します</li>
                </ol>
            </Box>
        </Box>
    );
}
