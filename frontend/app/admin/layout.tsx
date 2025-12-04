import Link from 'next/link';
import { Box, Stack, Text, HStack } from '@/src/core/ui';
import { css } from 'styled-system/css';

export default function AdminLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    return (
        <Box display="flex" minH="screen" bg="gray.100">
            {/* Sidebar */}
            <Box w="64" bg="indigo.900" color="white" flexShrink={0} shadow="xl" display="flex" flexDirection="column">
                <Box p="6" bg="indigo.950">
                    <Text as="h2" fontSize="xl" fontWeight="bold" letterSpacing="wider" color="white">管理画面</Text>
                    <Text fontSize="xs" color="indigo.300" mt="1">Fish Auction Admin</Text>
                </Box>
                <Stack as="nav" mt="6" px="2" spacing="1">
                    <Link href="/" className={css({ display: 'block', py: '3', px: '4', borderRadius: 'md', _hover: { bg: 'indigo.800' }, transition: 'colors', fontSize: 'sm', fontWeight: 'medium', color: 'white' })}>
                        &larr; トップに戻る
                    </Link>
                    <Box borderTop="1px solid" borderColor="indigo.800" my="4" mx="2"></Box>
                    <Link href="/admin" className={css({ display: 'block', py: '3', px: '4', borderRadius: 'md', _hover: { bg: 'indigo.800' }, transition: 'colors', fontSize: 'sm', fontWeight: 'medium', color: 'white' })}>
                        📊 ダッシュボード
                    </Link>
                    <Box borderTop="1px solid" borderColor="indigo.800" my="4" mx="2"></Box>
                    <Link href="/admin/fishermen" className={css({ display: 'block', py: '3', px: '4', borderRadius: 'md', _hover: { bg: 'indigo.800' }, transition: 'colors', fontSize: 'sm', fontWeight: 'medium', color: 'white' })}>
                        👨‍🌾 漁師管理
                    </Link>
                    <Link href="/admin/buyers" className={css({ display: 'block', py: '3', px: '4', borderRadius: 'md', _hover: { bg: 'indigo.800' }, transition: 'colors', fontSize: 'sm', fontWeight: 'medium', color: 'white' })}>
                        👔 中買人管理
                    </Link>
                    <Link href="/admin/venues" className={css({ display: 'block', py: '3', px: '4', borderRadius: 'md', _hover: { bg: 'indigo.800' }, transition: 'colors', fontSize: 'sm', fontWeight: 'medium', color: 'white' })}>
                        🏢 会場管理
                    </Link>
                    <Link href="/admin/auctions" className={css({ display: 'block', py: '3', px: '4', borderRadius: 'md', _hover: { bg: 'indigo.800' }, transition: 'colors', fontSize: 'sm', fontWeight: 'medium', color: 'white' })}>
                        📅 セリ管理
                    </Link>
                    <Link href="/admin/items" className={css({ display: 'block', py: '3', px: '4', borderRadius: 'md', _hover: { bg: 'indigo.800' }, transition: 'colors', fontSize: 'sm', fontWeight: 'medium', color: 'white' })}>
                        🐟 出品管理
                    </Link>
                    <Box borderTop="1px solid" borderColor="indigo.800" my="4" mx="2"></Box>
                    <Link href="/admin/invoice" className={css({ display: 'block', py: '3', px: '4', borderRadius: 'md', _hover: { bg: 'indigo.800' }, transition: 'colors', fontSize: 'sm', fontWeight: 'medium', color: 'white' })}>
                        💰 請求書発行
                    </Link>
                </Stack>
            </Box>

            {/* Main Content */}
            <Box as="main" flex="1" p="8" overflowY="auto">
                {children}
            </Box>
        </Box>
    );
}
