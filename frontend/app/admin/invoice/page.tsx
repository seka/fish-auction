'use client';

import { useInvoicePage } from './_hooks/useInvoicePage';
import { Box, Text, Card, Button, ModalBackdrop, ModalContent, Table, Thead, Tbody, Tr, Th, Td } from '@/src/core/ui';
import { css } from 'styled-system/css';

export default function InvoicePage() {
    const { state, actions } = useInvoicePage();

    return (
        <Box maxW="5xl" mx="auto" p="6">
            <Text as="h1" variant="h2" className={css({ color: 'gray.800' })} mb="8" pb="4" borderBottom="1px solid" borderColor="gray.200">
                請求書発行
            </Text>

            <Card padding="none" overflow="hidden">
                <Table>
                    <Thead>
                        <Tr className={css({ cursor: 'default', _hover: { bg: 'gray.50' } })}>
                            <Th>中買人ID</Th>
                            <Th>屋号・氏名</Th>
                            <Th className={css({ textAlign: 'right' })}>請求総額 (税込)</Th>
                        </Tr>
                    </Thead>
                    <Tbody>
                        {state.isLoading ? (
                            <Tr>
                                <Td colSpan={3} className={css({ py: '12', textAlign: 'center', color: 'gray.500' })}>
                                    読み込み中...
                                </Td>
                            </Tr>
                        ) : state.invoices.length === 0 ? (
                            <Tr className={css({ cursor: 'default', _hover: { bg: 'white' } })}>
                                <Td colSpan={3} className={css({ py: '12', textAlign: 'center', color: 'gray.500' })}>
                                    請求データはありません。
                                </Td>
                            </Tr>
                        ) : (
                            state.invoices.map((invoice) => (
                                <Tr key={invoice.buyerId} onClick={() => actions.setSelectedInvoice(invoice)}>
                                    <Td className={css({ fontSize: 'sm', color: 'gray.500', fontFamily: 'mono' })}>
                                        {invoice.buyerId}
                                    </Td>
                                    <Td className={css({ fontSize: 'sm', fontWeight: 'bold', color: 'gray.900' })}>
                                        {invoice.buyerName}
                                    </Td>
                                    <Td className={css({ textAlign: 'right', fontWeight: 'bold', color: 'indigo.700', fontSize: 'lg' })}>
                                        ¥{invoice.totalAmount.toLocaleString()}
                                    </Td>
                                </Tr>
                            ))
                        )}
                    </Tbody>
                </Table>
            </Card>

            {/* Detail Modal */}
            {state.selectedInvoice && (
                <ModalBackdrop onClick={() => actions.setSelectedInvoice(null)}>
                    <ModalContent onClick={(e) => e.stopPropagation()}>
                        <Box p="6" borderBottom="1px solid" borderColor="gray.200">
                            <Text variant="h3" fontWeight="bold" className={css({ color: 'gray.900' })}>請求書詳細</Text>
                        </Box>
                        <Box p="6">
                            <Box mb="6">
                                <Text fontSize="sm" className={css({ color: 'gray.700' })} mb="1">中買人ID</Text>
                                <Text fontWeight="bold" fontFamily="mono">{state.selectedInvoice.buyerId}</Text>
                            </Box>
                            <Box mb="6">
                                <Text fontSize="sm" className={css({ color: 'gray.700' })} mb="1">屋号・氏名</Text>
                                <Text fontWeight="bold" fontSize="lg">{state.selectedInvoice.buyerName}</Text>
                            </Box>
                            <Box mb="6">
                                <Text fontSize="sm" className={css({ color: 'gray.700' })} mb="1">請求総額 (税込)</Text>
                                <Text fontWeight="bold" fontSize="2xl" className={css({ color: 'indigo.700' })}>
                                    ¥{state.selectedInvoice.totalAmount.toLocaleString()}
                                </Text>
                            </Box>
                        </Box>
                        <Box p="6" borderTop="1px solid" borderColor="gray.200" display="flex" justifyContent="flex-end">
                            <Button variant="outline" onClick={() => actions.setSelectedInvoice(null)}>
                                閉じる
                            </Button>
                        </Box>
                    </ModalContent>
                </ModalBackdrop>
            )}
        </Box>
    );
}
