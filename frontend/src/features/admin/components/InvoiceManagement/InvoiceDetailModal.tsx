'use client';

import { Box, Button, Text } from '@atoms';
import { ModalBackdrop, ModalContent } from '@molecules';
import { css } from 'styled-system/css';
import { useTranslations } from 'next-intl';

import { InvoiceItem } from '../../types';

interface InvoiceDetailModalProps {
  invoice: InvoiceItem;
  onClose: () => void;
}

export const InvoiceDetailModal = ({ invoice, onClose }: InvoiceDetailModalProps) => {
  const t = useTranslations();

  return (
    <ModalBackdrop onClick={onClose}>
      <ModalContent onClick={(e) => e.stopPropagation()}>
        <Box p="6" borderBottom="1px solid" borderColor="gray.200">
          <Text variant="h3" fontWeight="bold" className={css({ color: 'gray.900' })}>
            {t('Admin.Invoice.modal_title')}
          </Text>
        </Box>
        <Box p="6">
          <Box mb="6">
            <Text fontSize="sm" className={css({ color: 'gray.700' })} mb="1">
              {t('Admin.Invoice.buyer_id')}
            </Text>
            <Text fontWeight="bold" fontFamily="mono">
              {invoice.buyerId}
            </Text>
          </Box>
          <Box mb="6">
            <Text fontSize="sm" className={css({ color: 'gray.700' })} mb="1">
              {t('Admin.Invoice.buyer_name')}
            </Text>
            <Text fontWeight="bold" fontSize="lg">
              {invoice.buyerName}
            </Text>
          </Box>
          <Box mb="6">
            <Text fontSize="sm" className={css({ color: 'gray.700' })} mb="1">
              {t('Admin.Invoice.total_amount')}
            </Text>
            <Text fontWeight="bold" fontSize="2xl" className={css({ color: 'indigo.700' })}>
              ¥{invoice.totalAmount.toLocaleString()}
            </Text>
          </Box>
        </Box>
        <Box
          p="6"
          borderTop="1px solid"
          borderColor="gray.200"
          display="flex"
          justifyContent="flex-end"
        >
          <Button variant="outline" onClick={onClose}>
            {t('Admin.Invoice.close')}
          </Button>
        </Box>
      </ModalContent>
    </ModalBackdrop>
  );
};
