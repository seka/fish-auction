'use client';

import { useInvoicePage } from './_hooks/useInvoicePage';
import { useTranslations } from 'next-intl';
import { Box, Text, Card, Button, ModalBackdrop, ModalContent, Table, Thead, Tbody, Tr, Th, Td } from '@/src/core/ui';
import { css } from 'styled-system/css';
import { EmptyState } from '../../_components/atoms/EmptyState';

export default function InvoicePage() {
    const t = useTranslations();
    const { state, actions } = useInvoicePage();

    return (
        <Box maxW="5xl" mx="auto" p="6">
            <Text as="h1" variant="h2" className={css({ color: 'gray.800' })} mb="8" pb="4" borderBottom="1px solid" borderColor="gray.200">
                {t('Admin.Invoice.title')}
            </Text>

            <Card padding="none" overflow="hidden">
                {state.isLoading ? (
                    <Box p="6" textAlign="center" className={css({ color: 'gray.600' })}>{t('Common.loading')}</Box>
                ) : state.invoices.length === 0 ? (
                    <EmptyState
                        message={t('Admin.Invoice.no_data')}
                        icon={<span role="img" aria-label="invoice">🧾</span>}
                    />
                ) : (
                    <Table>
                        <Thead>
                            <Tr className={css({ cursor: 'default', _hover: { bg: 'gray.50' } })}>
                                <Th>{t('Admin.Invoice.buyer_id')}</Th>
                                <Th>{t('Admin.Invoice.buyer_name')}</Th>
                                <Th className={css({ textAlign: 'right' })}>{t('Admin.Invoice.total_amount')}</Th>
                            </Tr>
                        </Thead>
                        <Tbody>
                            {state.invoices.map((invoice) => (
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
                            ))}
                        </Tbody>
                    </Table>
                )}
            </Card>

            {/* Detail Modal */}
            {state.selectedInvoice && (
                <ModalBackdrop onClick={() => actions.setSelectedInvoice(null)}>
                    <ModalContent onClick={(e) => e.stopPropagation()}>
                        <Box p="6" borderBottom="1px solid" borderColor="gray.200">
                            <Text variant="h3" fontWeight="bold" className={css({ color: 'gray.900' })}>{t('Admin.Invoice.modal_title')}</Text>
                        </Box>
                        <Box p="6">
                            <Box mb="6">
                                <Text fontSize="sm" className={css({ color: 'gray.700' })} mb="1">{t('Admin.Invoice.buyer_id')}</Text>
                                <Text fontWeight="bold" fontFamily="mono">{state.selectedInvoice.buyerId}</Text>
                            </Box>
                            <Box mb="6">
                                <Text fontSize="sm" className={css({ color: 'gray.700' })} mb="1">{t('Admin.Invoice.buyer_name')}</Text>
                                <Text fontWeight="bold" fontSize="lg">{state.selectedInvoice.buyerName}</Text>
                            </Box>
                            <Box mb="6">
                                <Text fontSize="sm" className={css({ color: 'gray.700' })} mb="1">{t('Admin.Invoice.total_amount')}</Text>
                                <Text fontWeight="bold" fontSize="2xl" className={css({ color: 'indigo.700' })}>
                                    ¥{state.selectedInvoice.totalAmount.toLocaleString()}
                                </Text>
                            </Box>
                        </Box>
                        <Box p="6" borderTop="1px solid" borderColor="gray.200" display="flex" justifyContent="flex-end">
                            <Button variant="outline" onClick={() => actions.setSelectedInvoice(null)}>
                                {t('Admin.Invoice.close')}
                            </Button>
                        </Box>
                    </ModalContent>
                </ModalBackdrop>
            )}
        </Box>
    );
}
