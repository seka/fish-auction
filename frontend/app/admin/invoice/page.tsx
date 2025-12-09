'use client';

import { useState } from 'react';
import { useInvoices } from './_hooks/useInvoice';
import { InvoiceItem } from '@/src/models';
import { Box, Text, Card, Button, ModalBackdrop, ModalContent, Table, Thead, Tbody, Tr, Th, Td } from '@/src/core/ui';
import { css } from 'styled-system/css';
import { styled } from 'styled-system/jsx';







export default function InvoicePage() {
    const { invoices } = useInvoices();
    const [selectedInvoice, setSelectedInvoice] = useState<InvoiceItem | null>(null);

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
                        {invoices.length === 0 ? (
                            <Tr className={css({ cursor: 'default', _hover: { bg: 'white' } })}>
                                <Td colSpan={3} className={css({ py: '12', textAlign: 'center', color: 'gray.500' })}>
                                    請求データはありません。
                                </Td>
                            </Tr>
                        ) : (
                            invoices.map((invoice) => (
                                <Tr key={invoice.buyerId} onClick={() => setSelectedInvoice(invoice)}>
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
            {selectedInvoice && (
                <ModalBackdrop onClick={() => setSelectedInvoice(null)}>
                    <ModalContent onClick={(e) => e.stopPropagation()}>
                        <Box p="6" borderBottom="1px solid" borderColor="gray.200">
                            <Text variant="h3" fontWeight="bold" className={css({ color: 'gray.900' })}>請求書詳細</Text>
                        </Box>
                        <Box p="6">
                            <Box mb="6">
                                <Text fontSize="sm" className={css({ color: 'gray.700' })} mb="1">中買人ID</Text>
                                <Text fontWeight="bold" fontFamily="mono">{selectedInvoice.buyerId}</Text>
                            </Box>
                            <Box mb="6">
                                <Text fontSize="sm" className={css({ color: 'gray.700' })} mb="1">屋号・氏名</Text>
                                <Text fontWeight="bold" fontSize="lg">{selectedInvoice.buyerName}</Text>
                            </Box>
                            <Box mb="6">
                                <Text fontSize="sm" className={css({ color: 'gray.700' })} mb="1">請求総額 (税込)</Text>
                                <Text fontWeight="bold" fontSize="2xl" className={css({ color: 'indigo.700' })}>
                                    ¥{selectedInvoice.totalAmount.toLocaleString()}
                                </Text>
                            </Box>
                        </Box>
                        <Box p="6" borderTop="1px solid" borderColor="gray.200" display="flex" justifyContent="flex-end">
                            <Button variant="outline" onClick={() => setSelectedInvoice(null)}>
                                閉じる
                            </Button>
                        </Box>
                    </ModalContent>
                </ModalBackdrop>
            )}
        </Box>
    );
}
